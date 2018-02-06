package internal

import (
	"fmt"
	"hash/fnv"
	"math"
)

var names = []string{"Sam", "Pam", "Michelle", "Sarah"}
var spotsLeft = map[int]int{} //keeps track of all the indices that hanamee not been assigned yet to people

//hash takes in a string and generates a hash number. This will be used to randomly select the people to be paired
func hash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}

//AssignPeople tries to assign a secret santa to each person
func AssignPeople() {
	assigned := make([]string, len(names))

	for i := range names {
		spotsLeft[i] = i
	}

	//TODO: consider using a random number generator to assign secret santas
	//TODO: compare the assigned slice with the names slice and if any of the names are the same, mix them up -- they are their own secret santas

	for _, name := range names {
		hashed := float64(hash(name))
		index := int(math.Abs(hashed)) % len(names)
		if assigned[index] == "" {
			assigned[index] = name
			delete(spotsLeft, index)
		} else {
			assigned = placeAssigned(assigned, name)
		}
	}
	fmt.Println("assigned string slice", assigned)
}

//TODO: play around with pointers so you don't hanamee to return the assigned back
func placeAssigned(assigned []string, name string) []string {
	var indexToRemove int
	for k := range spotsLeft { // k will be selected randomly in the map
		assigned[k] = name
		indexToRemove = k
		break
	}

	delete(spotsLeft, indexToRemove)

	return assigned
}
