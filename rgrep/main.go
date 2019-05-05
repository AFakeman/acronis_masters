package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "sync"
)

func printLoop(files chan string) {
    for file := range(files) {
        fmt.Println(file)
    }
}

const ChunkSize = 1024

func grep(file string, pattern string) (bool, error) {
    f, err := os.Open(file)
    if err != nil {
        return false, err
    }
    defer f.Close()
    buf := make([]byte, ChunkSize + 2 * len(pattern))
    var pos int64
    pos = 0
    for {
        nn, err := f.Read(buf)
        if err != nil {
            if err == io.EOF {
                break
            }
            return false, err
        }

        contains := boyerMoore(buf[:nn], pattern)
        if contains {
            return true, nil
        }

        pos += int64(nn)
        if nn == len(buf) {
            pos, err = f.Seek(int64(-len(pattern)), 1)
            if err != nil {
                return false, err
            }
        }
    }
    return false, nil
}

func grepLoop(pattern string, input chan string, output chan string, wg *sync.WaitGroup) {
    for file := range(input) {
        contains, err := grep(file, pattern)
        if err != nil {
            log.Printf("Error processing %s: %s\n", file, err.Error())
        }
        if contains {
            output <- file
        }
    }
    wg.Done()
    wg.Wait()
    close(output)
}

func main() {
    if len(os.Args) != 3 {
        log.Fatal("Usage: rgrep <pattern> <directory>")
    }
    pattern := os.Args[1]
    files := os.Args[2]

    preprocessingCache[pattern] = generatePreprocessing([]byte(pattern))

    input := make(chan string, 16)
    output := make(chan string, 16)
    wg := sync.WaitGroup{}
    wg.Add(1)


    go grepLoop(pattern, input, output, &wg)

    err := filepath.Walk(files,
        func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if info.Mode().IsRegular() {
                input <- path
            }
            return nil
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    close(input)
    printLoop(output)

}
