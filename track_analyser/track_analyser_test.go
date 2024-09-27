package track_analyser

import (
	"fmt"
	"testing"

	"github.com/johnforster/racetrack-go/core"
)

func TestAnalyser(t *testing.T) {
	tracks := AnalyseByFilePath("./test_data/big_track.bmp")
	expected_size := 2
	result := len(tracks)

	if result != expected_size {
		t.Errorf("Result was incorrect, got: %d, want: %d.", result, expected_size)
	}
}

func TestLookAround(t *testing.T) {
	base := core.Coordinate{5, 5}
	for c := range surroundingPixels(base, 1) {
		fmt.Println(c)
	}
}

func TestCountRings(t *testing.T) {
	tracks := AnalyseByFilePath("./test_data/big_track_six_rings.bmp")
	expected_size := 6
	result := len(tracks)

	fmt.Println("Number of loops: ", result)

	if result != expected_size {
		t.Errorf("Result was incorrect, got: %d, want: %d.", result, expected_size)
	}
}

func TestWithSinglePixelGaps(t *testing.T) {
	tracks := AnalyseByFilePath("./test_data/small_with_gaps.bmp")

	result := len(tracks)
	expected := 2
	if result != expected {
		t.Errorf("Result was incorrect, got: %d, wanted: %d.", result, expected)
	}
}

func TestFindCentre(t *testing.T) {
	tracks := AnalyseByFilePath("./test_data/find_centre.bmp")

	result := findCentre(tracks[0])
	expected := core.Coordinate{7, 5}
	if result != expected {
		t.Errorf("Result was incorrect, got: %d, wanted: %d.", result, expected)
	}
}
