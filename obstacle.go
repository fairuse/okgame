package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Obstacle struct {
	geom    Polygon
	enabled bool
	color   color.RGBA
}

func (o *Obstacle) translate(p Point) Obstacle {
	return Obstacle{geom: o.geom.translated(p), enabled: o.enabled, color: o.color}
}

func (o *Obstacle) draw() (vertices []ebiten.Vertex, indices []uint16) {
	// todo, check function signature
	if !o.enabled {
		return []ebiten.Vertex{}, []uint16{}
	}
	v, i := render(o.geom, o.color)
	return v, i
}
