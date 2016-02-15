package main

import (
    "syscall"
    "os"
    "log"
    "time"
)

func statfile(outfileName string) (logfileSize int64, logfileCreationTime time.Time) {
    var stat syscall.Stat_t
    err := syscall.Stat(outfileName, &stat)
    if os.IsNotExist(err) {
        logfileSize = 0
        logfileCreationTime = time.Now().UTC()
    } else if err != nil {
        log.Fatal(err)
    } else {
        // The file exists. Update our globals.
        logfileSize = stat.Size
        logfileCreationTime = time.Unix(stat.Ctim.Sec,
                                        stat.Ctim.Nsec).UTC()
    }
    return logfileSize, logfileCreationTime
}
