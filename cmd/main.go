package main

import (
	"fmt"

	"github.com/sachinagada/Secret_Santa/internal"
)

func main() {
	var names = []string{"Sarah", "Pam", "Lam", "Bam", "Tam", "Michelle", "Michelle "}

	fmt.Println(internal.HashingTest(names))
}
