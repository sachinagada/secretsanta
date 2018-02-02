package internal

import (
	"fmt"
	"hash/fnv"
)

var names = []string{"Sam", "Pam", "Michelle", "Sarah"}

//hash takes in a string and generates a hash number. This will be used to randomly select the people to be paired
func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func AssignPeople() {
	assigned := make([]string, len(names))
	for _, v := range names {
		hashed := int(hash(v)) % len(names)
		if assigned[hashed] == "" {
			assigned[hashed] = v
		} else {
			assigned = placeAssigned(assigned, v)
		}
	}
	fmt.Println("assigned string slice", assigned)
}

//TODO: play around with pointers so you don't have to return the assigned back
func placeAssigned(assigned []string, name string) []string {
	for i, v := range assigned {
		fmt.Println(i, v)
		if v == "" {
			assigned[i] = name
			break
		}
	}
	return assigned
}
