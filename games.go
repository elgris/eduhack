package main

import "math/rand"

var games = []int{0, 1, 2, 3, 4, 5, 6}

func getGames() []int {
	g := make([]int, len(games))
	copy(g, games)
	shuffleInt(g)
	return g
}

func shuffleInt(a []int) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func shuffleString(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
