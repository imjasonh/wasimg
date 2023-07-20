//go:build wasm
// +build wasm

package main

import (
	"log"
	"os"
)

func hello() {
	log.Println("hello adventurous wasm user")
	os.Exit(0)
}
