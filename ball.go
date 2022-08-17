package main

import (
	//	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

func line(sx float64, sy float64, tx float64, ty float64, col color.RGBA) (vertices []ebiten.Vertex, indices []uint16) {
	vertices = make([]ebiten.Vertex, 4)
	indices = []uint16{0, 1, 3, 1, 2, 3}

	cr := float32(col.R) / 0xff
	cg := float32(col.G) / 0xff
	cb := float32(col.B) / 0xff
	ca := float32(col.A) / 0xff

	l := math.Sqrt((ty-sy)*(ty-sy) + (tx-sx)*(tx-sx))
	dx := float64((tx - sx) / l * 4) // line width = 4
	dy := float64((ty - sy) / l * 4)

	perx := dy
	pery := -dx
	// fmt.Println(l,dx,dy, perx, pery)

	vertices[0] = ebiten.Vertex{DstX: float32(sx + perx), DstY: float32(sy + pery), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
	vertices[1] = ebiten.Vertex{DstX: float32(sx - perx), DstY: float32(sy - pery), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
	vertices[2] = ebiten.Vertex{DstX: float32(tx - perx), DstY: float32(ty - pery), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
	vertices[3] = ebiten.Vertex{DstX: float32(tx + perx), DstY: float32(ty + pery), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}

	// fmt.Println(vertices)

	return vertices, indices

}

func dotline(sx float64, sy float64, tx float64, ty float64, col color.RGBA) (vertices []ebiten.Vertex, indices []uint16) {
	const dotcount = 10
	vertices = make([]ebiten.Vertex, 4*dotcount)
	indices = make([]uint16, 0)
	i := uint16(0)
	for i = 0; i < dotcount; i += 1 {
		indices = append(indices, []uint16{0 + i*4, 1 + i*4, 3 + i*4, 1 + i*4, 2 + i*4, 3 + i*4}...)
	}

	cr := float32(col.R) / 0xff
	cg := float32(col.G) / 0xff
	cb := float32(col.B) / 0xff
	ca := float32(col.A) / 0xff

	l := math.Sqrt((ty-sy)*(ty-sy) + (tx-sx)*(tx-sx))
	dx := float64((tx - sx) / l * 4) // line width = 4
	dy := float64((ty - sy) / l * 4)

	perx := dy
	pery := -dx

	stepx := (tx - sx) / dotcount
	stepy := (ty - sy) / dotcount
	// fmt.Println(l,dx,dy, perx, pery)
	for i = 0; i < dotcount; i += 1 {
		fi := float64(i)

		vertices[0+i*4] = ebiten.Vertex{DstX: float32(sx + perx + stepx*fi), DstY: float32(sy + pery + stepy*fi), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
		vertices[1+i*4] = ebiten.Vertex{DstX: float32(sx - perx + stepx*fi), DstY: float32(sy - pery + stepy*fi), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
		vertices[2+i*4] = ebiten.Vertex{DstX: float32(sx - perx + stepx*(fi+0.5)), DstY: float32(sy - pery + stepy*(fi+0.5)), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
		vertices[3+i*4] = ebiten.Vertex{DstX: float32(sx + perx + stepx*(fi+0.5)), DstY: float32(sy + pery + stepy*(fi+0.5)), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
	}
	// fmt.Println(vertices)

	return vertices, indices

}

func ball(px float64, py float64, col color.RGBA) (vertices []ebiten.Vertex, indices []uint16) {
	const MAXSEG = 64
	const r = 10.0

	vertices = make([]ebiten.Vertex, 3*MAXSEG)
	indices = make([]uint16, 3*MAXSEG)

	cr := float32(col.R) / 0xff
	cg := float32(col.G) / 0xff
	cb := float32(col.B) / 0xff
	ca := float32(col.A) / 0xff

	var segment uint16
	for segment = 0; segment < MAXSEG; segment += 1 {
		dx := math.Sin(float64(segment)*2.0*math.Pi/MAXSEG) * r
		dy := math.Cos(float64(segment)*2.0*math.Pi/MAXSEG) * r
		dxn := math.Sin(float64(segment+1)*2.0*math.Pi/MAXSEG) * r
		dyn := math.Cos(float64(segment+1)*2.0*math.Pi/MAXSEG) * r

		vertices[segment*3+0] = ebiten.Vertex{DstX: float32(px), DstY: float32(py), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
		vertices[segment*3+1] = ebiten.Vertex{DstX: float32(px + dx), DstY: float32(py + dy), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}
		vertices[segment*3+2] = ebiten.Vertex{DstX: float32(px + dxn), DstY: float32(py + dyn), SrcX: 1, SrcY: 1, ColorR: cr, ColorG: cg, ColorB: cb, ColorA: ca}

		indices[segment*3+0] = segment*3 + 0
		indices[segment*3+1] = segment*3 + 1
		indices[segment*3+2] = segment*3 + 2
	}

	return vertices, indices
}
