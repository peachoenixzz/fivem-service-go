package randalphabet

import (
	"fmt"
	"math/rand"
	"time"
)

func VehiclePlate() string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randomLetters := make([]rune, 2)
	for i := range randomLetters {
		randomLetters[i] = letters[rand.Intn(len(letters))]
	}
	randomNumber := rand.Intn(900) + 100
	result := fmt.Sprintf("P%sK%d", string(randomLetters), randomNumber)

	return result
}
