package main

import (
	"os"
	"testing"

	"github.com/johnforster/racetrack-go/circuit"
	"github.com/johnforster/racetrack-go/track_analyser"
)

func TestCreatingBeziers(t *testing.T) {
	tracks := track_analyser.AnalyseByFilePath("./test_data/big_track.bmp")
	beziers := circuit.NewCircuit(tracks[0], tracks[1])
	json, err := beziers.ToJSON()
	bytes := []byte(json)

	os.WriteFile("test_server/test_beziers.json", bytes, 0644)

	if err != nil {
		panic("Didn't work")
	}
}
