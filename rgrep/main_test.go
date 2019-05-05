package main

import (
    "testing"
)

func checkBoyerMoore(t *testing.T, text string, pattern string, expect bool) {
    if boyerMoore([]byte(text), pattern) != expect {
        if expect {
            t.Errorf("Could not find [%s] in [%s]", pattern, text)
        } else {
            t.Errorf("Found [%s] in [%s]", pattern, text)
        }
    }
}

func TestSearchTrivial(t *testing.T) {
    checkBoyerMoore(t, "simple", "simple", true)
}

func TestSearchTrivialFalse(t *testing.T) {
    checkBoyerMoore(t, "simple", "s1mple", false)
}

func TestSearchNonTrivial(t *testing.T) {
    checkBoyerMoore(t,
        "A quick brown fox jumps over the lazy dog",
        "lazy",
        true,
    )
}

func TestSearchNonTrivialSimple(t *testing.T) {
    checkBoyerMoore(t,
        " lazy",
        "lazy",
        true,
    )
}

func TestSearchNonTrivialFail(t *testing.T) {
    checkBoyerMoore(t,
        "A quick brown fox jumps over the lazy dog",
        "la3y",
        false,
    )
}
