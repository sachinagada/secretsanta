package internal

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

var names = []string{"Sam", "Pam", "Michelle", "Sarah"}

//AssignPeople tries to assign a secret santa to each person
func AssignPeople() {
	assigned := make([]string, len(names))

	//TODO: compare the assigned slice with the names slice and if any of the names are the same, mix them up -- they are their own secret santas

	for _, name := range names {
		index := getRandomNumber(len(names))
		for assigned[index] != "" { //keep randomly generating numbers until a free spot in assigned is found
			fmt.Println("index", index)
			index = getRandomNumber(len(names))
		}
		assigned[index] = name

		fmt.Println(assigned)

	}
}

func getRandomNumber(n int) int {

	//TODO: look into seeding new sources for random generators:
	/* s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)*/

	i := rand.Intn(n * 77)
	fmt.Println("random number", i)
	return i % n
}

func finalAssign(assign []string) []string {
	changeNames := []string{}

	for i, v := range names {
		if assign[i] == v {
			changeNames = append(changeNames, v)
		}
	}

	//TODO: add the logic for changing the positions of the people who have themselves as their secret santa
	return assign
}

//HashingTest is testing to see if I can hash and assign Secret Santas that way instead of random number
func HashingTest(names []string) {
	secretSantaList := make([]string, len(names))
	h := fnv.New32a()

	for i := 0; i < len(names); i++ {
		h.Write([]byte(names[i]))
		secretSantaIndex := int(h.Sum32()) % len(names)
		fmt.Println(names[i], secretSantaIndex)

	}
	fmt.Println(secretSantaList)
}
