package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	title      string
	wraparound bool
	obstacles  []Obstacle
}

func (l *Level) draw() (vertices []ebiten.Vertex, indices []uint16) {
	v := make([]ebiten.Vertex, 0)
	i := make([]uint16, 0)
	for _, obj := range l.obstacles {
		ov, oi := obj.draw()
		vlen := uint16(len(v))
		ilen := len(i)
		v = append(v, ov...)
		i = append(i, oi...)
		for pos := ilen; pos < len(i); pos++ {
			i[pos] += vlen
		}
	}
	return v, i
}

func (l *Level) add(o Obstacle) {
	l.obstacles = append(l.obstacles, o)
}

func (l *Level) findHit(ball Point, direction Point) (*Line, *Obstacle) {
	var hitInfo *Line // src = hit point, dst = normal direction
	var hitObj *Obstacle
	for _, obj := range l.obstacles {
		interp := intersectPolygonNorm(ball, direction, obj.geom)
		if interp != nil {
			if hitInfo == nil {
				hitInfo = interp
				hitObj = &obj
			} else {
				if distance(ball,hitInfo.src) > distance(ball,interp.src) {
					hitInfo = interp
					hitObj = &obj
				}
			}
		}
	}
	return hitInfo, hitObj
}
