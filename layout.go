package main

import "math"

type Orientation struct {
	f0, f1, f2, f3 float32
	b0, b1, b2, b3 float32
	start_angle    float32
}

var PointyTop Orientation = Orientation{
	float32(math.Sqrt(3.0)),
	float32(math.Sqrt(3.0) / 2.0),
	float32(0.0),
	float32(3.0 / 2.0),
	float32(math.Sqrt(3.0) / 3.0),
	float32(-1.0 / 3.0),
	float32(0.0),
	float32(2.0 / 3.0),
	float32(0.5),
}

var FlatTop Orientation = Orientation{
	float32(3.0 / 2.0),
	float32(0.0),
	float32(math.Sqrt(3.0) / 2.0),
	float32(math.Sqrt(3.0)),
	float32(2.0 / 3.0),
	float32(0.0),
	float32(-1.0 / 3.0),
	float32(math.Sqrt(3.0) / 3.0),
	float32(0.0),
}

type Layout struct {
	orientation Orientation
	size        Point
	origin      Point
}

func NewLayout(orientation *Orientation, size, origin Point) *Layout {
	return &Layout{
		orientation: *orientation,
		size:        size,
		origin:      origin,
	}
}

func (l *Layout) HexToPixel(hex *Hex) Point {
	ori := &l.orientation
	x := (ori.f0*float32(hex.q) + ori.f1*float32(hex.r)) * l.size.x
	y := (ori.f2*float32(hex.q) + ori.f3*float32(hex.r)) * l.size.y
	return Point{x + l.origin.x, y + l.origin.y}
}

func (l *Layout) PixelToHex(p Point) FractionalHex {
	ori := &l.orientation
	pt := Point{
		x: (p.x - l.origin.x) / l.size.x,
		y: (p.y - l.origin.x) / l.size.y,
	}
	q := ori.b0*pt.x + ori.b1*pt.y
	r := ori.b2*pt.x + ori.b3*pt.y
	return FractionalHex{q, r, (-q - r)}
}

func (l *Layout) HexCornerOffset(corner int) Point {

	angle := float64(2.0 * math.Pi * (l.orientation.start_angle + float32(corner)) / 6)
	return Point{l.size.x * float32(math.Cos(angle)), l.size.y * float32(math.Sin(angle))}
}

func (l *Layout) GetCorners(hex *Hex) []Point {
	var corners []Point
	center := l.HexToPixel(hex)
	for i := 0; i < 6; i++ {
		offset := l.HexCornerOffset(i)
		corners = append(corners, Point{center.x + offset.x, center.y + offset.y})
	}
	return corners
}
