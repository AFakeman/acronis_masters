package main

import (
    "path/filepath"
    "fmt"
    "io"
    "log"
    "os"
    "sync"
)

type BoyerMoorePreprocessing struct {
    bc map[byte][]int
}

var preprocessingCache map[string]*BoyerMoorePreprocessing = map[string]*BoyerMoorePreprocessing{}

func generatePreprocessing(pattern []byte) *BoyerMoorePreprocessing {
    bm := &BoyerMoorePreprocessing{bc: map[byte][]int{}}

    // Bad character rule
    prev_idx := make(map[byte]int)
    for idx, b := range(pattern) {
        if _, ok := bm.bc[b]; !ok {
            bm.bc[b] = make([]int, len(pattern))
        }
        if prev, ok := prev_idx[b]; ok {
            for i := prev; i < idx; i++ {
                // If the symbol is found, shift to it
                // pattscale
                // test
                //   test
                // bm.bc['t'][2] = 2 - 0 = 2

                bm.bc[b][i] = idx - prev
            }
        } else {
            for i := 0; i < idx; i++ {
                // If the symbol is not found before, move the start of
                // the string past the mismatched symbol
                // pabtscale
                // best
                //    best
                // bm.bc['b'][2] = 2 + 1 = 3

                bm.bc[b][i] = idx + 1
            }
        }
        prev_idx[b] = idx
    }
    for b, bc := range(bm.bc) {
        for idx := prev_idx[b]; idx < len(pattern); idx++ {
            bc[idx] = idx - prev_idx[b]
        }
    }

    // TODO: other two rules

    return bm
}

func boyerMoore(text []byte, pattern_str string) bool {
    pattern := []byte(pattern_str)
    bm, ok := preprocessingCache[pattern_str]
    if !ok {
        preprocessingCache[pattern_str] = generatePreprocessing(pattern)
    }

    for k := len(pattern) - 1; k < len(text); {
        start := k - len(pattern) + 1
        shift := 0
        for i := k; i >= start; i-- {
            p := i - start  // cursor inside pattern
            if text[i] != pattern[p] {
                // Appy the bad character rule
                var bc_shift int
                if bc, ok := bm.bc[text[i]]; ok {
                    bc_shift = bc[p]
                } else {
                    bc_shift = p + 1
                }
                shift = bc_shift
                break
            }
        }
        if shift == 0 {
            return true
        } else {
            k += shift
        }
    }
    return false
}

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
                continue
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
