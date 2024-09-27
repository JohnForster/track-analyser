package core

import "github.com/johnforster/racetrack-go/set"

type Track = *set.OrderedSet[Coordinate]

type Coordinate struct {
	X, Y int
}
