package track_analyser

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"iter"
	"os"

	"github.com/johnforster/racetrack-go/core"
	"github.com/johnforster/racetrack-go/set"

	"github.com/sergeymakinen/go-bmp"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
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

// ? - Is this worth keeping?
// func GetTracksFromImage2(image image.Image) ([]core.Track, error) {
// 	bounds := image.Bounds()

// 	tracks := []core.Track{}

// 	var start_point core.Coordinate

// 	all_coordinates := []core.Coordinate{}
// 	for y := 0; y < bounds.Max.Y; y++ {
// 		for x := 0; x < bounds.Max.X; x++ {
// 			r, _, _, _ := image.At(x, y).RGBA()
// 			is_black := r == 0
// 			all_coordinates = append(all_coordinates, core.Coordinate{X: x, Y: y})

// 			if is_black && start_point.X == 0 && start_point.Y == 0 {
// 				start_point = core.Coordinate{X: x, Y: y}
// 			}
// 		}
// 	}

// 	complete := false
// 	current_point := start_point
// 	checked := map[core.Coordinate]bool{start_point: true}
// 	track_coords := []core.Coordinate{}
// 	for complete {
// 		var unused_neighbours []core.Coordinate
// 		for n := range surroundingPixels(current_point, 1) {
// 			_, ok := checked[n]
// 			if ok {
// 				unused_neighbours = append(unused_neighbours, n)
// 			}
// 		}

// 		if len(unused_neighbours) == 0 {
// 			complete = true
// 		}
// 		if len(unused_neighbours) == 1 {
// 			track_coords = append(track_coords, unused_neighbours[0])
// 		}
// 		if len(unused_neighbours) > 1 {
// 			panic("Not yet implemented")
// 		}
// 	}

// 	var track *set.OrderedSet[core.Coordinate]
// 	track.AddMulti(track_coords...)
// 	tracks = append(tracks, track)

// 	return tracks, nil

// }

func GetTracksFromImage(image image.Image) ([]core.Track, error) {
	// TODO - Apply image thinning on submitted images.
	// image, err := thinImage(image)

	// if err != nil {
	// 	return []core.Track{}, err
	// }

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
	if should_include(c) && !track.Has(c) {
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

// Not currently used.
func thinImage(img image.Image) (image.Image, error) {
	result := image.NewGray(img.Bounds())
	draw.Draw(result, result.Bounds(), img, img.Bounds().Min, draw.Src)

	mat, err := gocv.ImageGrayToMatGray(result)

	if err != nil {
		return nil, err
	}

	inv := gocv.NewMat()
	defer inv.Close()

	gocv.BitwiseNot(mat, &inv)

	dst := gocv.NewMat()
	defer dst.Close()

	contrib.Thinning(inv, &dst, contrib.ThinningGuoHall)

	res := gocv.NewMat()
	defer res.Close()
	gocv.BitwiseNot(dst, &res)

	img, err = res.ToImage()
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Not currently used.
func writePNG(img image.Image, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outFile, img)

	return nil
}
