package main

type BoyerMoorePreprocessing struct {
    bc map[byte][]int
    gp []int
}

var preprocessingCache map[string]*BoyerMoorePreprocessing = map[string]*BoyerMoorePreprocessing{}

func generatePreprocessing(pattern []byte) *BoyerMoorePreprocessing {
    bm := &BoyerMoorePreprocessing{
        bc: map[byte][]int{},
        gp: make([]int, len(pattern) + 1),
    }

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

    // Good prefix rule
    // If a suffix of a pattern matches until a certain symbol, try to find
    // rightmost occurence of the same suffix, but with a different symbol.

    // |f[i]| contains starting position of the widest border of |pattern[i:]|,
    // meaning its biggest suffix that is also its prefix.
    // The borders can be computed in o(n) time.

    f := make([]int, len(pattern) + 1)
    j := len(pattern) + 1
    f[len(pattern)] = len(pattern) + 1

    for i := len(pattern); i > 0; {
        // Try different indices until we get a border
        for ; j <= len(pattern) && pattern[i - 1] != pattern[j - 1]; {
            if bm.gp[j] == 0 {
                // If a border can't be expanded, we should remember it.
                // It may be a good shift.
                bm.gp[j] = j - i;
            }
            j=f[j];
        }
        i--
        j--
        f[i] = j;
    }

    // Second case, when the a part of the matching suffix is actually a prefix
    // of the entire pattern. The pattern can be shifted as far as its widest
    // matching border allows.

    j = f[0];
    for i := 0; i <= len(pattern); i++ {
        if (bm.gp[i]==0) {
            bm.gp[i]=j;
        }
        if (i == j) {
            j = f[j];
        }
    }

    return bm
}

func boyerMoore(text []byte, pattern_str string) bool {
    pattern := []byte(pattern_str)
    bm, ok := preprocessingCache[pattern_str]
    if !ok {
        bm = generatePreprocessing(pattern)
        preprocessingCache[pattern_str] = bm
    }

    for k := len(pattern) - 1; k < len(text); {
        start := k - len(pattern) + 1
        shift := 0
        for i := k; i >= start; i-- {
            p := i - start  // cursor inside pattern
            if text[i] != pattern[p] {
                // Apply the bad character rule
                var bc_shift int
                if bc, ok := bm.bc[text[i]]; ok {
                    bc_shift = bc[p]
                } else {
                    bc_shift = p + 1
                }
                shift = bc_shift

                // Apply the good prefix rule
                var gp_shift = bm.gp[p + 1]
                if gp_shift > shift {
                    shift = gp_shift
                }

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
