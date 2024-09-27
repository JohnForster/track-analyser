package circuit

import (
	"encoding/json"

	"github.com/johnforster/fitcurves"
	"github.com/johnforster/racetrack-go/core"
)

type Circuit struct {
	inner []fitcurves.Bezier
	outer []fitcurves.Bezier
}

func NewCircuit(inner core.Track, outer core.Track) *Circuit {
	innerBeziers := TrackToBeziers(inner)
	outerBeziers := TrackToBeziers(outer)

	return &Circuit{inner: innerBeziers, outer: outerBeziers}
}

func TrackToBeziers(t core.Track) []fitcurves.Bezier {
	points := []fitcurves.Point{}
	for _, p := range t.ToList() {
		points = append(points, fitcurves.NewPoint(float64(p.X), float64(p.Y)))
	}

	return fitcurves.FitCurves(points, 50.0)
}

type CircuitJSON struct {
	Inner [][][]float64 `json:"inner"`
	Outer [][][]float64 `json:"outer"`
}

func bezierToJSON(b fitcurves.Bezier) [][]float64 {
	p0 := []float64{b.P0.X, b.P0.Y}
	p1 := []float64{b.P1.X, b.P1.Y}
	p2 := []float64{b.P2.X, b.P2.Y}
	p3 := []float64{b.P3.X, b.P3.Y}
	return [][]float64{p0, p1, p2, p3}
}

func (c *Circuit) ToJSON() ([]byte, error) {
	inner := [][][]float64{}
	for _, b := range c.inner {
		inner = append(inner, bezierToJSON(b))
	}

	outer := [][][]float64{}
	for _, b := range c.outer {
		outer = append(outer, bezierToJSON(b))
	}

	circuitJSON := CircuitJSON{Inner: inner, Outer: outer}
	return json.Marshal(circuitJSON)
}
