package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	title      string
	wraparound bool
	obstacles  []Obstacle
}

func (l *Level) enableAll() {
	for nr, _ := range l.obstacles {
		l.obstacles[nr].enabled = true
	}
}

func (l *Level) draw() (vertices []ebiten.Vertex, indices []uint16) {
	v := make([]ebiten.Vertex, 0)
	i := make([]uint16, 0)
	for _, obj := range l.obstacles {
		if obj.enabled == false {
			continue
		}
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
	for nr, obj := range l.obstacles {
		if obj.enabled == false {
			continue
		}
		interp := intersectPolygonNorm(ball, direction, obj.geom)
		if interp != nil {
			if hitInfo == nil {
				hitInfo = interp
				hitObj = &l.obstacles[nr]
			} else {
				if distance(ball, hitInfo.src) > distance(ball, interp.src) {
					hitInfo = interp
					hitObj = &l.obstacles[nr]
				}
			}
		}
	}
	return hitInfo, hitObj
}
