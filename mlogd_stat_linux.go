package main

import (
    "syscall"
    "os"
    "log"
    "time"
)

func statfile(outfileName string) (logfileSize int64, logfileCreationTime time.Time, err error) {
    var stat syscall.Stat_t
    err = syscall.Stat(outfileName, &stat)
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
    return logfileSize, logfileCreationTime, err
}

func select_stdin() {
    var r_fdset syscall.FdSet
    for i := 0; i < 16; i++ {
        r_fdset.Bits[i] = 0
    }
    r_fdset.Bits[0] = 1
    var timeval syscall.Timeval
    timeval.Sec = 60
    timeval.Usec = 0
    n, selerr := syscall.Select(1, &r_fdset, nil, nil, &timeval)
    if selerr != nil {
        logger.Warning(selerr)
    }
}
