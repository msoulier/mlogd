package main

import (
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
)

var (
    timestamps = true
    localtime = true
    maxsize int64 = 0
    maxage = 0
    logfileSize int64 = 0
    logfileCreationTime = time.Now().UTC()
    logger = logging.MustGetLogger("mlogd")
    debug = false
    isaFile = true
)

func init() {
    const (
        defaultMaxSize = 5*1024*1024
        defaultMaxAge = 3600*24
    )
    flag.BoolVar(&timestamps, "timestamps", false, "Prefix all output lines with timestamps")
    flag.Int64Var(&maxsize, "maxsize", defaultMaxSize, "Maximum size of logfile in bytes before rotation")
    flag.IntVar(&maxage, "maxage", defaultMaxAge, "Maximum age of logfile in seconds before rotation")
    flag.BoolVar(&localtime, "localtime", false, "Render timestamps in localtime instead of UTC")
    flag.BoolVar(&debug, "debug", false, "Debug logging in mlogd")
    flag.Parse()

    if debug {
        logging.SetLevel(logging.DEBUG, "mlogd")
    } else {
        logging.SetLevel(logging.INFO, "mlogd")
    }
    format := logging.MustStringFormatter(
        `%{color}%{time:15:04:05.000} ▶ %{level} %{color:reset} %{message}`,
        )
    logging.SetFormatter(format)
    /* stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
    stderrFormatter := logging.NewBackendFormatter(stderrBackend, format)
    stderrBackendLevelled := logging.AddModuleLevel(stderrBackend)
    if debug {
        stderrBackendLevelled.SetLevel(logging.DEBUG, "mlogd")
    } else {
        stderrBackendLevelled.SetLevel(logging.INFO, "mlogd")
    }
    logging.SetBackend(stderrBackendLevelled, stderrFormatter)
    log.Printf("debug is %s\n", debug)
    log.Printf("level is %d\n", logging.GetLevel("")) */
}

func main() {
    // Input is always stdin.
    input := bufio.NewScanner(os.Stdin)
    var outfile *os.File
    var err error
    // Output is the file supplied on the command line.
    if (len(os.Args[1:]) > 0) {
        outfileName := os.Args[len(os.Args)-1]
        if (outfileName == "-") {
            outfile = os.Stdout
            isaFile = false
            logger.Debug("outfile set to stdout")
        } else {
            // If the logfile exists already, stat it and update the
            // logfileSize and logfileAge globals.
            logfileSize, logfileCreationTime = statfile(outfileName)
            logger.Debugf("outfile exists already, size is %d bytes, creation time is %s", logfileSize, logfileCreationTime)
            outfile, err = os.OpenFile(outfileName,
                                       os.O_WRONLY | os.O_CREATE | os.O_APPEND,
                                       0600)
            if (err != nil) {
                log.Fatal(err)
            }
        }
    } else {
        os.Stderr.WriteString(usage)
        flag.PrintDefaults()
        os.Exit(1)
    }
    output := bufio.NewWriter(io.Writer(outfile))

    // Loop over stdin until EOF.
    var count int64 = 0
    for input.Scan() {
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
        outBytes, err := output.WriteString(input.Text() + "\n")
        if err != nil {
            log.Fatalf("Write error: %s\n", err)
        }
        logfileSize += int64(outBytes)
        output.Flush()
        if count % lineFrequencyCheck == 0 {
            logger.Debugf("logfileSize is now %d, rollover at %d",
                logfileSize, maxsize)
            if logfileSize > maxsize && isaFile {
                logger.Debug("Rolling over logfile")
            }
        }
    }
    logger.Infof("EOF @ %d bytes", logfileSize)
}
