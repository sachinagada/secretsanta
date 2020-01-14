package main

import (
	"fmt"

	"github.com/sachinagada/Secret_Santa/internal"
)

// TODO: add the front-end page where the users can insert their name and email address
func main() {
	var names = []string{"Sarah", "Pam", "Lam", "Bam", "Tam", "Michelle"}
	secretSanta, err := internal.AssignSecretSanta(names)
	fmt.Println("assigned Secret Santas:", secretSanta)
	fmt.Println("err", err)

}
