package track_analyser

import (
	"fmt"
	"image"
	"iter"
	"os"

	"github.com/johnforster/racetrack-go/core"
	"github.com/johnforster/racetrack-go/set"

	"github.com/sergeymakinen/go-bmp"
)

func AnalyseByFilePath(path string) []core.Track {
	file, err := os.Open(path)

	if err != nil {
		fmt.Print("Error opening path")
		panic(err)
	}

	image, err := bmp.Decode(file)

	if err != nil {
		fmt.Print("Error decoding image")
		panic(err)
	}

	tracks, err := GetTracksFromImage(image)

	if err != nil {
		fmt.Print("Error decoding image")
		panic(err)
	}

	return tracks
}

func GetTracksFromImage(image image.Image) ([]core.Track, error) {
	bounds := image.Bounds()

	tracks := []core.Track{}
	accounted_for := set.NewSet[core.Coordinate]()

	create_predicate := func(t core.Track) func(c core.Coordinate) bool {
		return func(c core.Coordinate) bool {
			already_counted := accounted_for.Has(c) || t.Has(c)
			if already_counted {
				return false
			}
			within_bounds := c.X > 0 && c.X < bounds.Max.X && c.Y >= 0 && c.Y < bounds.Max.Y

			if !within_bounds {
				return false
			}
			r, _, _, _ := image.At(c.X, c.Y).RGBA()
			is_black := r == 0

			if !is_black {
				return false
			}

			return true
		}
	}

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, _, _, _ := image.At(x, y).RGBA()
			is_black := r == 0

			if is_black && !accounted_for.Has(core.Coordinate{X: x, Y: y}) {
				track := createNewTrack(create_predicate, core.Coordinate{X: x, Y: y})
				tracks = append(tracks, track)
				accounted_for = accounted_for.UnionWithTracked(track)
			}
		}
	}

	return tracks, nil
}

func createNewTrack(create_test func(t core.Track) func(c core.Coordinate) bool, start core.Coordinate) core.Track {
	track := set.NewOrderedSet[core.Coordinate]()
	test := create_test(track)

	recursivelyFollowTrack(track, test, start)
	return track
}

func recursivelyFollowTrack(track core.Track, should_include func(core.Coordinate) bool, c core.Coordinate) {
	if should_include(c) {
		track.Add(c)

		for neighbour := range surroundingPixels(c, 2) {
			recursivelyFollowTrack(track, should_include, neighbour)
		}
	}
}

func surroundingPixels(original core.Coordinate, distance int) iter.Seq[core.Coordinate] {
	return func(yield func(core.Coordinate) bool) {
		for dy := -distance; dy <= distance; dy++ {
			for dx := -distance; dx <= distance; dx++ {
				new_coords := core.Coordinate{X: original.X + dx, Y: original.Y + dy}
				if new_coords == original {
					continue
				}

				if !yield(new_coords) {
					return
				}
			}
		}
	}
}

func findCentre(t core.Track) core.Coordinate {
	totalX := 0
	totalY := 0

	for _, c := range t.ToList() {
		totalX += c.X
		totalY += c.Y
	}

	return core.Coordinate{X: totalX / t.Size(), Y: totalY / t.Size()}
}
