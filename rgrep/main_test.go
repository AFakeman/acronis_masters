package main

import (
    "testing"
)

func TestSearchTrivial(t *testing.T) {
    text := []byte("simple")
    pattern := []byte("simple")

    if !boyerMoore(text, pattern, nil) {
        t.Fail()
    }
}

func TestSearchTrivialFalse(t *testing.T) {
    text := []byte("simple")
    pattern := []byte("s1mple")

    if boyerMoore(text, pattern, nil) {
        t.Fail()
    }
}

func TestSearchNonTrivial(t *testing.T) {
    text := []byte("A quick brown fox jumps over the lazy dog")
    pattern := []byte("lazy")

    if !boyerMoore(text, pattern, nil) {
        t.Fail()
    }
}

func TestSearchNonTrivial2(t *testing.T) {
    text := []byte(" lazy")
    pattern := []byte("lazy")

    if !boyerMoore(text, pattern, nil) {
        t.Fail()
    }
}

func TestSearchNonTrivialFail(t *testing.T) {
    text := []byte("A quick brown fox jumps over the lazy dog")
    pattern := []byte("la3y")

    if boyerMoore(text, pattern, nil) {
        t.Fail()
    }
}
