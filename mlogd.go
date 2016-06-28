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
    "time"
    "log"
    "io"
    "flag"
    "github.com/op/go-logging"
)

const (
    usage = "mlogd [options] <logfile path>\n"
    lineFrequencyCheck = 100
    VERSION = "0.91"
)

var (
    timestamps = true
    localtime = true
    maxsize int64 = 0
    maxage int64 = 0
    logfileSize int64 = 0
    logfileCreationTime = time.Now().UTC()
    logger *logging.Logger
    debug = false
    isaFile = true
    flush = false
    version = false
    nfiles = 7
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
    flag.Parse()

    // The colour logger is problematic for capturing logs in text files.
    //format := logging.MustStringFormatter(
    //    `%{color}%{time:15:04:05.000} ▶ %{level} %{color:reset} %{message}`,
    //    )
    format := logging.MustStringFormatter(
        `%{time:15:04:05.000} ▶ %{level} %{message}`,
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
    rv := now.Format("200601_2150405")
    return rv
}

func manage_rotated_files(linkName string, nfiles int) {
    logger.Debugf("manage_rotated_files: nfiles is %d", nfiles)
    dirname := path.Dir(linkName)
    basename := path.Base(linkName)
    logger.Debugf("dirname is %s, basename is %s", dirname, basename)
    files, err := ioutil.ReadDir(dirname)
    if err != nil {
        logger.Fatal(err)
    }
    // An array of files to return
    old_logfiles := make([]os.FileInfo, 0, 100)
    for _, file := range files {
        logger.Debugf("found file in log dir: %s", file.Name())
        // We only want .log files.
        if strings.HasSuffix(file.Name(), ".log") {
            logger.Debug("    this is a logfile")
            if file.Name() == basename {
                logger.Debug("    don't count this one")
            } else {
                old_logfiles = append(old_logfiles, file)
            }
        }
    }
    logger.Debugf("old_logfiles is now %s, with %d elements", old_logfiles, len(old_logfiles))
    if len(old_logfiles) > nfiles {
        todelete := len(old_logfiles) - nfiles
        logger.Debugf("need to delete old logfiles: %d", todelete)
        sort.Reverse(ByName(old_logfiles))
        for _, file := range old_logfiles[:todelete] {
            todelete_path := dirname + "/" + file.Name()
            logger.Debugf("deleting: %s", todelete_path)
            if err := os.Remove(todelete_path); err != nil {
                logger.Fatal(err)
            }
        }
    }
}

func rollover(linkName string, outfileName string, outfile *os.File, nfiles int) (string, *os.File, error) {
    var err error
    newOutfileName := strings.TrimSuffix(linkName, ".log") + "-" + gettimesuffix(time.Now()) + ".log"
    logger.Debugf("rollover: new filename is %q", newOutfileName)
    // Close and reopen outfile
    outfile.Close()
    outfile, err = os.OpenFile(newOutfileName,
                               os.O_WRONLY | os.O_CREATE | os.O_APPEND,
                               0600)
    // Move the symlink
    if err = os.Remove(linkName); err != nil {
        logger.Fatal(err)
    }
    if err = os.Symlink(newOutfileName, linkName); err != nil {
        logger.Fatal(err)
    }
    manage_rotated_files(linkName, nfiles)
    return newOutfileName, outfile, err
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

    // Output is the file supplied on the command line.
    if len(os.Args[1:]) > 0 {
        timesuffix := gettimesuffix(time.Now())
        // FIXME: make .log extension configurable
        linkName = os.Args[len(os.Args)-1]
        outfileName = timesuffix + ".log"
        outfileName = strings.TrimSuffix(linkName, ".log") + "-" + timesuffix + ".log"
        logger.Debugf("linkName is %q, outfileName is %q", linkName, outfileName)
        manage_rotated_files(linkName, nfiles)
        if linkName == "-" {
            outfile = os.Stdout
            isaFile = false
            logger.Debug("outfile set to stdout")
        } else {
            // If the logfile exists already, stat it and update the
            // logfileSize and logfileAge globals.
            if linkContents, err := os.Readlink(linkName); err != nil {
                logger.Debugf("linkName %q does not exist yet", linkName)
                if err := os.Symlink(outfileName, linkName); err != nil {
                    log.Fatal(err)
                }
            } else {
                // The symlink exists. It is now our output file name.
                logger.Debugf("linkName %q exists, reading and using it", linkName)
                logger.Debugf("link points to %q", linkContents)
                outfileName = linkContents
            }
            logfileSize, logfileCreationTime, err = statfile(outfileName)
            if err != nil && os.IsNotExist(err) {
                logger.Debugf("outfile %q does not yet exist - creating", outfileName)
            } else {
                logger.Debugf("outfile %q exists already, size is %d bytes, creation time is %s", outfileName, logfileSize, logfileCreationTime)
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

    // A convenience for running checks at regular intervals.
    // FIXME: use a goroutine
    var ticker int64 = 0

selectloop:
    for {
        logger.Debug("going into select on stdin")
        select_stdin()

        // Loop over stdin until EOF.
        var count int64 = 0
        for {
            logger.Debugf("count is %d", count)
            line, readerr := input.ReadString('\n')
            if readerr != nil {
                logger.Debugf("read error: %#v", readerr)
                if readerr == io.EOF {
                    logger.Debug("EOF")
                    break selectloop
                } else {
                    logger.Debugf("breaking read loop after %d lines", count+1)
                    break
                }
            }
            count++
            ticker++
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
            logfileSize += int64(outBytes)
            if flush {
                output.Flush()
            }
            // FIXME: check at startup too in case we don't hit this frequency count
            if ticker % lineFrequencyCheck == 0 {
                logger.Debugf("logfileSize is now %d, rollover at %d",
                    logfileSize, maxsize)
                now := time.Now().UTC()
                if logfileSize > maxsize && isaFile {
                    logger.Debug("Rolling over logfile")
                    outfileName, outfile, err = rollover(linkName, outfileName, outfile, nfiles)
                    output.Flush()
                    output = bufio.NewWriter(io.Writer(outfile))
                    if err != nil {
                        log.Fatal(err)
                    }
                    logfileSize = 0
                    logfileCreationTime = now.UTC()
                }
                // And check current time for rollover.
                duration := now.Sub(logfileCreationTime)
                logger.Debugf("It has been %f seconds since file creation", duration.Seconds())
                logger.Debugf("maxage is %d seconds", maxage)
                if int64(duration.Seconds()) >= maxage && isaFile {
                    logger.Debug("Rolling over logfile")
                    outfileName, outfile, err = rollover(linkName, outfileName, outfile, nfiles)
                    output.Flush()
                    output = bufio.NewWriter(io.Writer(outfile))
                    if err != nil {
                        log.Fatal(err)
                    }
                    logfileSize = 0
                    logfileCreationTime = now.UTC()
                }
            }
        }
    }
    output.Flush()
    if isaFile {
        outfile.Close()
    }
}
