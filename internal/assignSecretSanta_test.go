package internal

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssignSecretSanta(t *testing.T) {
	for i := 0; i <= 10; i++ {
		namesArray := generateRandomNameArray()
		fmt.Println("Testing with namesArray of length: ", len(namesArray))
		secretSantaArray, err := AssignSecretSanta(namesArray)
		if err != nil {
			assert.True(t, len(namesArray) < 2)
		} else {
			assert.Equal(t, len(namesArray), len(secretSantaArray)) //makle sure the arrays are the same length
			secretSantaMap := map[string]string{}
			for i, name := range secretSantaArray {
				assert.NotEqual(t, namesArray[i], name) //make sure the names at the same position aren't the same
				//adding the names from the secretSanta array to the map to ensure all the names from the original namesArray are in the secretSantaArray
				secretSantaMap[name] = name
			}

			for _, name := range namesArray {
				_, ok := secretSantaMap[name]        //if the name is in the secretSantaMap, ok will be true. Should always be true
				assert.True(t, ok, name, namesArray) //prints out the name and original namesArray if the name is not in the map
			}
		}
	}
}

func generateRandomNameArray() []string {
	const s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@._/12345567890"

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	lengthOfArray := r.Intn(len(s)) //generating a random number for the length of the array
	names := make([]string, lengthOfArray)
	for i := 0; i < lengthOfArray; i++ {
		lengthOfName := r.Intn(len(s)) + 1 //the name can be of any length greater than 0. Adding 1 to ensure it's not an empty string
		name := make([]byte, lengthOfName) //create a byte array. The string of this will be the name in the ith position of the names array
		for nameIndex := 0; nameIndex < lengthOfName; nameIndex++ {
			name[nameIndex] = s[rand.Intn(len(s))]
		}
		names[i] = string(name)
	}
	return names
}
