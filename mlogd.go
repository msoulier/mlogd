package main

import (
    "sort"
    "path"
    "io/ioutil"
    "fmt"
    "runtime"
    "strings"
    "bufio"
    "os"
    "os/exec"
    "time"
    "log"
    "io"
    "flag"
    "github.com/op/go-logging"
    "regexp"
    "os/signal"
    "syscall"
    mlib "github.com/msoulier/mlib"
)

const (
    usage = "mlogd [options] <logfile path>\n"
    VERSION = "1.3.5"
    shutdown_wait_time = 5
)

var (
    timestamps = true
    localtime = true
    maxsize int64 = 0
    maxage int64 = 0
    logfileSize int64 = 0
    logfileCreationTime = time.Now()
    logger *logging.Logger
    debug = false
    isaFile = true
    flush = false
    version = false
    nfiles = 7
    post = ""
    altext = ""
    rotation_required = false
    rotation_frequency time.Duration = 60
    select_timeout time.Duration = 60
    shutdown_asap = false
)

func init() {
    const (
        defaultMaxSize = 5*1024*1024
        defaultMaxAge = 3600*24
    )
    flag.BoolVar(&timestamps, "timestamps", false, "Prefix all output lines with timestamps")
    flag.Int64Var(&maxsize, "maxsize", defaultMaxSize, "Maximum size of logfile in bytes before rotation")
    flag.Int64Var(&maxage, "maxage", defaultMaxAge, "Maximum age of logfile in seconds before rotation")
    flag.BoolVar(&localtime, "localtime", false, "Render timestamps in localtime instead of UTC")
    flag.BoolVar(&debug, "debug", false, "Debug logging in mlogd")
    flag.BoolVar(&flush, "flush", false, "Flush output buffer on each line")
    flag.BoolVar(&version, "version", false, "Print version and quit")
    flag.IntVar(&nfiles, "nfiles", 7, "The number of log files to keep")
    flag.StringVar(&post, "post", "", "A post-rotation shell command to run on the rotated file")
    flag.StringVar(&altext, "altext", "", "Additional comma-separated file extensions to consider log files")
    flag.Parse()

    // The colour logger is problematic for capturing logs in text files.
    //format := logging.MustStringFormatter(
    //    `%{color}%{time:15:04:05.000} ▶ %{level} %{color:reset} %{message}`,
    //    )
    format := logging.MustStringFormatter(
        `%{time:15:04:05.000} %{level} %{message}`,
        )
    stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
    stderrFormatter := logging.NewBackendFormatter(stderrBackend, format)
    stderrBackendLevelled := logging.AddModuleLevel(stderrFormatter)
    logging.SetBackend(stderrBackendLevelled)
    if debug {
        stderrBackendLevelled.SetLevel(logging.DEBUG, "mlogd")
    } else {
        stderrBackendLevelled.SetLevel(logging.INFO, "mlogd")
    }
    logger = logging.MustGetLogger("mlogd")

    logger.Debugf("starting mlogd version %s on %s\n", VERSION, runtime.GOOS)
    logger.Debugf("timestamps is %q", timestamps)
    logger.Debugf("maxsize is %q", maxsize)
    logger.Debugf("maxage is %q", maxage)
    logger.Debugf("localtime is %q", localtime)
    logger.Debugf("debug is %q", debug)
    logger.Debugf("flush is %q", flush)
    logger.Debugf("nfiles is %q", nfiles)
    logger.Debugf("post is %q", post)
    logger.Debugf("altext is %q", altext)
    // set rotation_frequency to the lesser of itself or the maxage argument
    if time.Duration(maxage) < rotation_frequency {
        rotation_frequency = time.Duration(maxage)
    }
    if time.Duration(maxage) < select_timeout {
        select_timeout = time.Duration(maxage)
    }
}

// For sorting FileInfo objects by Name
type ByName []os.FileInfo
func (a ByName) Len() int               { return len(a) }
func (a ByName) Swap(i, j int)          { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool     { return a[i].Name() < a[j].Name() }

func gettimesuffix(now time.Time) string {
    logger.Debugf("gettimesuffix: now is %s", now)
    // http://fuckinggodateformat.com/
    // %Y%m%e%H%M%S
    // rfc 3339 - seriously??
    rv := now.Format("20060102150405")
    logger.Debugf("returning format %s", rv)
    // The timesuffix returned should never have spaces in it
    if strings.Contains(rv, " ") {
        panic(rv)
    }
    return rv
}

/*
 * Run the post command on the file in a goroutine.
 */
func do_post(filePath string) {
    cmd := exec.Command(post, filePath)
    logger.Debugf("in do_post: cmd = %s", cmd)
    err := cmd.Run()
    if err != nil {
        logger.Error(err)
    }
}

/*
 * The purpose of this function is to clean up old logfiles based on the
 * nfiles setting, and run the post action on the latest rotated file if
 * required.
 */
func manage_rotated_files(linkName string, postFile string) {
    logger.Debugf("manage_rotated_files: nfiles is %d", nfiles)
    logger.Debugf("linkName is %s, postFile is %s", linkName, postFile)
    dirname := path.Dir(linkName)
    files, err := ioutil.ReadDir(dirname)
    if err != nil {
        logger.Fatal(err)
    }
    // Compose the array of file extensions to consider logfiles.
    logfile_exts := make([]string, 0, 10)
    // The first extension is always .log
    logfile_exts = append(logfile_exts, ".log")
    for _, ext := range strings.Split(altext, ",") {
        if ext != "" {
            logger.Debugf("adding '%s' to the extension list", ext)
            logfile_exts = append(logfile_exts, ext)
        }
    }
    logger.Debugf("logfile_exts: %s, size %d", logfile_exts, len(logfile_exts))
    // An array of files to return
    old_logfiles := make([]os.FileInfo, 0, 100)
    for _, file := range files {
        logger.Debugf("found file in log dir: %s", file.Name())
        // Skip any dotfiles.
        if strings.HasPrefix(file.Name(), ".") {
            logger.Debugf("skipping dotfile %s", file.Name())
            continue
        } else if regular := file.Mode().IsRegular(); ! regular {
            logger.Debugf("%s is not a regular file, skipping", file.Name())
            continue
        }
        // We only want supported file extensions.
        for _, ext := range logfile_exts {
            logger.Debugf("    looping on extension '%s'", ext)
            if strings.HasSuffix(file.Name(), ext) {
                logger.Debug("    this is a logfile")
                old_logfiles = append(old_logfiles, file)
            }
        }
    }
    sort.Reverse(ByName(old_logfiles))
    logger.Debugf("old_logfiles is now %s, with %d elements", old_logfiles, len(old_logfiles))
    if len(old_logfiles) > nfiles {
        todelete := len(old_logfiles) - nfiles
        logger.Debugf("need to delete old logfiles: %d", todelete)
        for _, file := range old_logfiles[:todelete] {
            todelete_path := path.Join(dirname, file.Name())
            if todelete_path != postFile {
                logger.Debugf("deleting: %s", todelete_path)
                if err := os.Remove(todelete_path); err != nil {
                    logger.Fatal(err)
                }
            } else {
                logger.Warning("Was about to delete rotated file:", postFile)
            }
        }
    }
    // Ok, they've been managed. Now run any post action on the file that was
    // just rotated.
    if post != "" {
        logger.Debugf("post action configured: '%s'", post)
        last_log := postFile
        logger.Debugf("rotated file is %s", last_log)
        go do_post(last_log)
    }
}

func rollover(linkName string, outfileName string, outfile *os.File) (string, *os.File, error) {
    var err error
    newOutfileName := strings.TrimSuffix(linkName, ".log") + "-" + gettimesuffix(time.Now()) + ".log"
    logger.Debugf("rollover: new filename is %q", newOutfileName)
    // Close and reopen outfile
    outfile.Close()
    outfile, err = os.OpenFile(newOutfileName,
                               os.O_WRONLY | os.O_CREATE | os.O_APPEND,
                               0600)
    // Move the symlink
    logger.Debugf("removing symlink: %s", linkName)
    if err = os.Remove(linkName); err != nil {
        logger.Fatal(err)
    }
    logger.Debugf("recreating symlink: %s -> %s", linkName, newOutfileName)
    if err = os.Symlink(newOutfileName, linkName); err != nil {
        logger.Fatal(err)
    }
    manage_rotated_files(linkName, outfileName)
    return newOutfileName, outfile, err
}

func parse_creation(outfileName string) (ctime time.Time) {
    logger.Debugf("parse_creation: outfileName is %s", outfileName)
    var datetime = regexp.MustCompile(`(\d{14})\.log`)
    if datetime.MatchString(outfileName) {
        // Matched name.
        datetime_string := datetime.FindStringSubmatch(outfileName)[1]
        logger.Debugf("parsed out datetime: %q", datetime_string)
        zone, _ := time.Now().Zone()
        logger.Debugf("timezone is now %s", zone)
        t, err := time.Parse("20060102150405 MST", datetime_string + " " + zone)
        if err == nil {
            logger.Debugf("time %q", t)
            return t
        } else {
            logger.Debugf("time parse error, using now: %s", err)
            return time.Now()
        }
    } else {
        logger.Debug("failed to match time string, using now")
        return time.Now()
    }
}

func check_rotation() {
    // Loop forever, checking rotation_frequency seconds if the file needs
    // rotating.
    for {
        // Check rotation conditions - first by size
        logger.Debugf("logfileSize is now %d, rollover at %d",
            logfileSize, maxsize)
        now := time.Now()

        if logfileSize > maxsize && isaFile {
            logger.Debug("flagging rotation required by size")
            rotation_required = true
        }

        // and then by age
        duration := now.Sub(logfileCreationTime)
        logger.Debugf("It has been %f seconds since file creation", duration.Seconds())
        logger.Debugf("maxage is %d seconds", maxage)

        if int64(duration.Seconds()) >= maxage && isaFile {
            logger.Debug("flagging rotation required by age")
            rotation_required = true
        }

        logger.Debugf("sleeping for %d seconds", rotation_frequency)
        time.Sleep(rotation_frequency * time.Second)
    }
}

func handle_hup(input <-chan os.Signal) {
    for {
        logger.Debugf("handle_hup: waiting on signal")
        sig := <-input
        logger.Debugf("HUP signal received: %s", sig)
        // Right now we do nothing with a sighup. svlogd re-reads its config, but
        // ours is all command-line.
    }
}

func handle_alarm(input <-chan os.Signal) {
    // Loop indefinitely to handle alarm signals.
    for {
        logger.Debugf("handle_alarm: waiting on signal")
        sig := <-input
        logger.Info("ALRM signal received: forcing rotation", sig)
        rotation_required = true
    }
}

func handle_shutdown(input <-chan os.Signal) {
    logger.Debugf("handle_shutdown: waiting on signal")
    sig := <-input
    logger.Info("INT or TERM signal received:", sig)
    shutdown_asap = true
    // If we don't shut down in shutdown_wait_time, then we're not getting any
    // input, so trigger an exit.
    time.Sleep(shutdown_wait_time * time.Second)
    logger.Debugf("handle_shutdown: timeout on shutdown")
    os.Exit(0)
}

func main() {
    var outfile *os.File
    var err error
    var linkName string
    var outfileName string

    if version {
        fmt.Printf("mlogd version %s on %s\n", VERSION, runtime.GOOS)
        os.Exit(0)
    }

    // On SIGHUP we should close and reopen all logs.
    hup_sigs := make(chan os.Signal, 1)
    // On SIGALRM we should force logfile rotation.
    alarm_sigs := make(chan os.Signal, 1)
    // On SIGTERM or SIGINT we should shutdown.
    shutdown_sigs := make(chan os.Signal, 1)

    // And launch our goroutines to handle them.
    // FIXME: Can likely consolidate all these.
    go handle_hup(hup_sigs)
    go handle_alarm(alarm_sigs)
    go handle_shutdown(shutdown_sigs)

    // Register the channels for those signals.
    signal.Notify(hup_sigs, syscall.SIGHUP)
    signal.Notify(alarm_sigs, syscall.SIGALRM)
    signal.Notify(shutdown_sigs, syscall.SIGINT, syscall.SIGTERM)

    // Output is the file supplied on the command line.
    if len(flag.Args()) > 0 {
        timesuffix := gettimesuffix(time.Now())
        // FIXME: make .log extension configurable
        linkName = os.Args[len(os.Args)-1]
        if ! strings.HasPrefix(linkName, "/") {
            log.Fatal("not an absolute path: ", linkName)
        }
        outfileName = timesuffix + ".log"
        outfileName = strings.TrimSuffix(linkName, ".log") + "-" + timesuffix + ".log"
        logger.Debugf("linkName is %q, outfileName is %q", linkName, outfileName)
        if linkName == "-" {
            outfile = os.Stdout
            isaFile = false
            logger.Debug("outfile set to stdout")
        } else {
            // If the logfile exists already, stat it and update the
            // logfileSize global.
            if linkContents, err := os.Readlink(linkName); err != nil {
                logger.Debugf("%s", err)
                logger.Debugf("linkName %q does not exist yet - creating", linkName)
                if err := os.Symlink(outfileName, linkName); err != nil {
                    log.Fatal(err)
                }
            } else {
                // The symlink exists. It is now our output file name.
                logger.Debugf("linkName %q exists, reading and using it", linkName)
                logger.Debugf("link points to %q", linkContents)
                outfileName = linkContents
            }
            logfileStat, err := os.Stat(outfileName)
            if err != nil && os.IsNotExist(err) {
                logger.Debugf("outfile %q does not yet exist - creating", outfileName)
                // And if it was just created, then the logfileCreationTime is
                // now, of course.
                logfileCreationTime = time.Now()
            } else {
                // The logfile pointed to by the main log symlink does exist.
                // The stat tells us its size but its creation date and time
                // must be parsed out of the name, as we usually don't find
                // filesystems storing ctime.
                logfileSize = logfileStat.Size()
                logfileCreationTime = parse_creation(outfileName)
                logger.Debugf("outfile %q exists already, size is %d bytes, creation time is %s - using", outfileName, logfileSize, logfileCreationTime)
            }
            outfile, err = os.OpenFile(outfileName,
                                       os.O_WRONLY | os.O_CREATE | os.O_APPEND,
                                       0600)
            if err != nil {
                log.Fatal(err)
            }
        }
    } else {
        os.Stderr.WriteString(usage)
        flag.PrintDefaults()
        os.Exit(1)
    }
    output := bufio.NewWriter(io.Writer(outfile))

    // Input is always stdin.
    input := bufio.NewReader(os.Stdin)

    // Start goroutine to check for rotation.
    go check_rotation();

selectloop:
    for {
        if rotation_required {
            outfileName, outfile, err = rollover(linkName, outfileName, outfile)
            logger.Debug("flushing output")
            output.Flush()
            output = bufio.NewWriter(io.Writer(outfile))
            if err != nil {
                log.Fatal(err)
            }
            logfileSize = 0
            logfileCreationTime = time.Now()
            rotation_required = false
        }

        logger.Debugf("going into select on stdin, timeout is %ds", select_timeout)
        readable := mlib.SelectStdin(select_timeout)
        logger.Debug("back from select")

        if readable {
            // Loop over stdin until EOF.
            var count int64 = 0
            for {
                // Break out of the inner reading loop if rotation is needed.
                if rotation_required {
                    break;
                }
                logger.Debugf("count is %d", count)
                line, readerr := input.ReadString('\n')
                if readerr != nil {
                    logger.Debugf("read error: %#v", readerr)
                    if readerr == io.EOF {
                        logger.Debug("EOF")
                        break selectloop
                    } else if shutdown_asap {
                        logger.Debug("shutdown_asap")
                        break selectloop
                    } else {
                        logger.Debugf("breaking read loop after %d lines", count+1)
                        break
                    }
                }
                count++
                if timestamps {
                    var now time.Time
                    if localtime {
                        now = time.Now()
                    } else {
                        now = time.Now().UTC()
                    }
                    output.WriteString(now.Format(time.StampMicro) + " ")
                }
                outBytes, err := output.WriteString(line)
                if err != nil {
                    log.Fatalf("Write error: %s\n", err)
                }
                logger.Debugf("wrote %d bytes to output file", outBytes)
                logfileSize += int64(outBytes)
                if flush {
                    logger.Debug("flushing output")
                    output.Flush()
                }
            }
        } else {
            logger.Debug("stdin not readable")
        }

    }
    output.Flush()
    if isaFile {
        outfile.Close()
    }
    os.Exit(0)
}
