package main

import (
    "bufio"
    "os"
    "time"
    "log"
    "io"
    "flag"
    "syscall"
)

const (
    usage = "mlogd [options] <logfile path>\n"
)

var timestamps = true
var localtime = true
var maxsize = 0
var maxage = 0
var logfileSize int64 = 0
var logfileCreationTime = time.Now().UTC()

func init() {
    const (
        defaultMaxSize = 5*1024*1024
        defaultMaxAge = 3600*24
    )
    flag.BoolVar(&timestamps, "timestamps", false, "Prefix all output lines with timestamps")
    flag.IntVar(&maxsize, "maxsize", defaultMaxSize, "Maximum size of logfile in bytes before rotation")
    flag.IntVar(&maxage, "maxage", defaultMaxAge, "Maximum age of logfile in seconds before rotation")
    flag.BoolVar(&localtime, "localtime", false, "Render timestamps in localtime instead of UTC")
}

func statfile(outfileName string) {
    var stat syscall.Stat_t
    err := syscall.Stat(outfileName, &stat)
    if os.IsNotExist(err) {
        return
    } else if err != nil {
        log.Fatal(err)
    } else {
        // The file exists. Update our globals.
        logfileSize = stat.Size
        logfileCreationTime = stat.Ctimespec
    }
}

func main() {
    flag.Parse()
    // Input is always stdin.
    input := bufio.NewScanner(os.Stdin)
    var outfile *os.File
    var err error
    // Output is the file supplied on the command line.
    if (len(os.Args[1:]) > 0) {
        outfileName := os.Args[len(os.Args)-1]
        if (outfileName == "-") {
            outfile = os.Stdout
        } else {
            // If the logfile exists already, stat it and update the
            // logfileSize and logfileAge globals.
            statfile(outfileName)
            outfile, err = os.OpenFile(outfileName,
                                       os.O_WRONLY | os.O_CREATE,
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
    for input.Scan() {
        if timestamps {
            var now time.Time
            if localtime {
                now = time.Now()
            } else {
                now = time.Now().UTC()
            }
            output.WriteString(now.Format(time.StampMicro) + ": ")
        }
        output.WriteString(input.Text() + "\n")
        output.Flush()
    }
}
