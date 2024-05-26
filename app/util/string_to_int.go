package util

import (
	"hash/fnv"
	"math/rand"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RandomInt returns a random integer for the given string.
func UserIdStringToInt(input string) int {
	// Initialize the hash function
	h := fnv.New64a()
	h.Write([]byte(input))
	seed := h.Sum64()

	// Initialize the random number generator with the seed
	generator := rand.New(rand.NewSource(int64(seed)))

	// Generate a random integer
	return generator.Int()
}

var adjectives = [...]string{
	"adorable", "hermoso", "inteligente", "determinado", "ansioso", "fiero", "gracioso", "feliz",
	"inquisitivo", "alegre", "amable", "animado", "misterioso", "noble", "optimista", "orgulloso",
	"peculiar", "resiliente", "sincero", "considerado", "animado", "vibrante", "ingenioso", "apasionado",
}

var animals = [...]string{
	"oso", "tigre", "león", "elefante", "jirafa", "cebra", "guepardo", "hipopótamo",
	"rinoceronte", "canguro", "koala", "panda", "mono", "gorila", "loro", "águila",
	"búho", "lobo", "zorro", "conejo", "ciervo", "ardilla", "tortuga", "serpiente",
}

var adjectives_animals = []string{}

func init() {
	// Perform inner join of adjectives and animals
	for _, adjective := range adjectives {
		for _, animal := range animals {
			capitalized := cases.Title(language.Und).String(animal + " " + adjective)
			adjectives_animals = append(adjectives_animals, capitalized)
		}
	}
}

// AssignUserName assigns a name to a user based on two lists of words.
func UserName(user_id string) string {
	adjectives_animals_idx := UserIdStringToInt(user_id) % len(adjectives_animals)
	return adjectives_animals[adjectives_animals_idx]
}
