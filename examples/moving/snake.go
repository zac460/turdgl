package main

import (
	"image/color"
	"math"
	"time"

	tgl "github.com/zac460/turdgl"
	"golang.org/x/exp/constraints"
)

const (
	maxNodeDistPx   = 80
	numSegments     = 30
	headSize        = 30
	bodyScaleFactor = 0.97
)

type snake struct {
	head     *tgl.Circle
	body     []*tgl.Circle
	velocity *tgl.Vec // velocity in px/s
}

// Snake constructs a new snake based on the given head position.
func NewSnake(headPos tgl.Vec) *snake {
	headStyle := tgl.Style{
		Colour:    color.RGBA{255, 255, 255, 255},
		Thickness: 0,
	}

	bodyStyle := tgl.Style{
		Colour:    color.RGBA{255, 255, 255, 255},
		Thickness: 4,
	}

	// Construct body in with segments stretched out...
	var b []*tgl.Circle
	for i := 0; i < numSegments-1; i++ {
		segmentDiameter := headSize * math.Pow(bodyScaleFactor, float64(i))
		segment := tgl.NewCircle(
			segmentDiameter, segmentDiameter,
			tgl.Vec{X: headPos.X, Y: headPos.Y + headSize*float64(i)},
			bodyStyle,
		)
		b = append(b, segment)
	}
	s := snake{
		head: tgl.NewCircle(headSize, headSize, headPos, headStyle),
		body: b,
	}
	// ...then align the body segments
	s.updateBodyPos()

	return &s
}

// Draw draws the snake on the provided frame buffer.
func (s *snake) Draw(buf *tgl.FrameBuffer) {
	const markerSize = 4
	markerStyle := tgl.Style{Colour: color.RGBA{255, 0, 0, 0}, Thickness: 0}

	// Draw head
	s.head.Draw(buf)

	// Draw body
	for _, c := range s.body {
		// Draw segment
		c.Draw(buf)

		// Draw markers
		lPos := c.Marker(math.Pi / 2 * 3)
		lMarker := tgl.NewCircle(markerSize, markerSize, lPos, markerStyle)
		lMarker.Draw(buf)

		rPos := c.Marker(math.Pi / 2)
		rMarker := tgl.NewCircle(markerSize, markerSize, rPos, markerStyle)
		rMarker.Draw(buf)
	}
}

// Update recalculates the snake's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check snake isn't out of bounds.
func (s *snake) Update(dt time.Duration, buf *tgl.FrameBuffer) {
	// Update the head
	newX := s.head.Pos.X + s.velocity.X*dt.Seconds()
	newY := s.head.Pos.Y + s.velocity.Y*dt.Seconds()
	const segmentRad float64 = headSize / 2
	newX = Constrain(newX, segmentRad, float64(buf.Width())-segmentRad-1)
	newY = Constrain(newY, segmentRad, float64(buf.Height())-segmentRad-1)
	s.head.Pos = tgl.Vec{X: newX, Y: newY}
	s.head.Direction = tgl.Normalise(*s.velocity)

	// Update the body
	s.updateBodyPos()
}

func (s *snake) updateBodyPos() {
	for i, node := range s.body {
		var nodeAhead *tgl.Circle
		if i == 0 {
			nodeAhead = s.head
		} else {
			nodeAhead = s.body[i-1]
		}
		// If node is too far away from the node ahead of it...
		if tgl.Dist(node.Pos, nodeAhead.Pos) > nodeAhead.Width() {
			// Move the node to be adjacent to the node ahead
			diff := tgl.Sub(nodeAhead.Pos, node.Pos)
			node.Move(tgl.Sub(diff, diff.SetMag(nodeAhead.Width())))
			node.Direction = tgl.Normalise(tgl.Sub(nodeAhead.Pos, node.Pos))
		}
	}
}

// Constrain keeps a number between lower and upper bounds.
func Constrain[T constraints.Ordered](x, lower, upper T) T {
	switch {
	case x < lower:
		return lower
	case x > upper:
		return upper
	default:
		return x
	}
}
