package internal

import (
	"hash/fnv"
)

//HashingTest is testing to see if I can hash and assign Secret Santas that way instead of random number
func HashingTest(names []string) (secretSantaList []string) {
	secretSantaList = make([]string, len(names))
	h := fnv.New32a()

	for i := 0; i < len(names); i++ {
		currentName := names[i]
		h.Write([]byte(names[i]))
		secretSantaIndex := int(h.Sum32()) % len(names) //hash the name and get the index to assign to the secret santa

		//don't want the current person to be their own secret santa
		for i == secretSantaIndex || secretSantaList[secretSantaIndex] != "" {

			//if it's the last name in the array and the only available space is the last space,
			//put the currentName in the 0th index and the one at 0th index in the current position.
			//Otherwise, this would be stuck in an infinite for loop
			if i == len(names)-1 && i == secretSantaIndex && secretSantaList[secretSantaIndex] == "" {
				secretSantaList[0], secretSantaList[i] = currentName, secretSantaList[0]
				return
			}
			//use linear probing to just place it in the next available position
			secretSantaIndex = (secretSantaIndex + 1) % len(names)
		}
		secretSantaList[secretSantaIndex] = currentName
	}
	return
}
