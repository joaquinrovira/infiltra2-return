package util

import (
	"hash/fnv"
	"math/rand"
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

var names = [...]string{
	"Alex Carter", "Jordan West", "Taylor Ward", "Casey Whateley", "Morgan Armitage",
	"Riley Thurston", "Robin Pickman", "Avery Gilman", "Quinn Alhazred",
	"Sam Waite", "Jamie Nyarla", "Cameron Curwen", "Reese Sothoth", "River Dagon",
	"Shay Azathoth", "Emery Shub", "Parker Allen", "Sage Lynch", "Rowan Marsh",
	"Devin Whateley", "Kai Olney", "Hunter Blake", "Taylor Lorry",
}

func UserName(user_id string) string {
	adjectives_animals_idx := UserIdStringToInt(user_id) % len(names)
	return names[adjectives_animals_idx]
}
