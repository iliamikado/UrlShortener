package main

import (
	"os"
)

func main() {
	os.Exit(0) // want "os.Exit is forbidden"
	f()
}

func f() {
	var a = 0
	a++
}
