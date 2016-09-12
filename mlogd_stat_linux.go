package main

import (
    "syscall"
    "os"
    "log"
)

func statfile(outfileName string) (logfileSize int64, err error) {
    var stat syscall.Stat_t
    err = syscall.Stat(outfileName, &stat)
    if os.IsNotExist(err) {
        logfileSize = 0
    } else if err != nil {
        log.Fatal(err)
    } else {
        // The file exists. Update our globals.
        logfileSize = stat.Size
    }
    return logfileSize, err
}

func select_stdin(timeout_secs int64) (bool) {
    var r_fdset syscall.FdSet
    var timeout syscall.Timeval
    timeout.Sec = timeout_secs
    timeout.Usec = 0
    for i := 0; i < 16; i++ {
        r_fdset.Bits[i] = 0
    }
    r_fdset.Bits[0] = 1
    _, selerr := syscall.Select(1, &r_fdset, nil, nil, &timeout)
    if selerr != nil {
        logger.Warning(selerr)
    }
    // Is it really ready to read or did we time out?
    if r_fdset.Bits[0] == 1 {
        return true
    } else {
        return false
    }
}
