package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float32
}

type Hex struct {
	q int
	r int
	s int
}

var hex_directions [6]Hex = [6]Hex{
	{1, 0, -1},
	{1, -1, 0},
	{0, -1, 1},
	{-1, 0, 1},
	{-1, 1, 0},
	{0, 1, -1},
}

func (h1 *Hex) Add(h2 Hex) Hex {
	return Hex{h1.q + h2.q, h1.r + h2.r, h1.s + h2.s}
}

func (h *Hex) Neighbors() [6]Hex {
	var neighbors [6]Hex
	for i, d := range hex_directions {
		neighbors[i] = h.Add(d)
	}
	return neighbors
}

func (h *Hex) ToString() string {
	return fmt.Sprintf("Hex{q: %d, r: %d, s: %d}", h.q, h.r, h.s)
}

type FractionalHex struct {
	q float32
	r float32
	s float32
}

func (fh *FractionalHex) Round() *Hex {
	q := int(math.RoundToEven(float64(fh.q)))
	r := int(math.RoundToEven(float64(fh.r)))
	s := int(math.RoundToEven(float64(fh.s)))

	q_diff := math.Abs(float64(float32(q) - fh.q))
	r_diff := math.Abs(float64(float32(r) - fh.r))
	s_diff := math.Abs(float64(float32(s) - fh.s))
	if q_diff > r_diff && q_diff > s_diff {
		q = -r - s
	} else if r_diff > s_diff {
		r = -q - s
	} else {
		s = -q - r
	}
	return &Hex{q, r, s}
}
