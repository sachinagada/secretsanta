package main

import (
	"fmt"

	"github.com/sachinagada/Secret_Santa/internal"
)

func main() {
	var names = []string{"Sarah", "Pam", "Lam", "Bam", "Tam", "Michelle"}
	secretSanta, err := internal.AssignSecretSanta(names)
	fmt.Println("assigned Secret Santas:", secretSanta)
	fmt.Println("err", err)

}
