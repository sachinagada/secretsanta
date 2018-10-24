package main

import (
	"github.com/sachinagada/Secret_Santa/internal"
)

func main() {
	// internal.AssignPeople()
	var names = []string{"Sam", "Pam", "Michelle", "Sarah"}

	internal.HashingTest(names)
}
