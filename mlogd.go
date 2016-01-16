package main

import (
    "bufio"
    "os"
    "time"
    "log"
    "io"
)

var timestamps = true

func main() {
    // Input is always stdin.
    input := bufio.NewScanner(os.Stdin)
    var outfile *os.File
    var err error
    // Output is the file supplied on the command line.
    if (len(os.Args[1:]) > 0) {
        outfileName := os.Args[1]
        if (outfileName == "-") {
            outfile = os.Stdout
        } else {
            outfile, err = os.OpenFile(outfileName,
                                       os.O_WRONLY | os.O_CREATE,
                                       0600)
            if (err != nil) {
                log.Fatal(err)
            }
        }
    }
    output := bufio.NewWriter(io.Writer(outfile))

    // Loop over stdin until EOF.
    for input.Scan() {
        if (timestamps) {
            now := time.Now().UTC()
            output.WriteString(now.Format(time.StampMicro) + ": ")
        }
        output.WriteString(input.Text() + "\n")
        output.Flush()
    }
}
