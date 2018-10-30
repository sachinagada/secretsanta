package internal

import (
	"errors"
	"hash/fnv"
)

//AssignSecretSanta uses hashing algorithm to hash each name and assign it to the person in the hash value % len(names) position.
//If a name is already at that position, it will use linear probing method to put in the next available position. This function
//uses the same concepts as a hash map to hash the key (name in this case) to assign the key to a location.
func AssignSecretSanta(names []string) (secretSantaList []string, err error) {
	if len(names) < 2 {
		return nil, errors.New("Need more than 1 participant")
	}

	secretSantaList = make([]string, len(names))
	h := fnv.New32a()

	for i, currentName := range names {
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
			//use linear probing to place the name in the next available position
			secretSantaIndex = (secretSantaIndex + 1) % len(names)
		}
		secretSantaList[secretSantaIndex] = currentName
	}
	return
}
