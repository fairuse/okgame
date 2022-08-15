package main

import (
	"math"
	"image/color"
        "github.com/hajimehoshi/ebiten/v2"
)

func ball(px float64, py float64, col color.RGBA) (vertices []ebiten.Vertex, indices []uint16) {
     const MAXSEG = 64
     const r = 4.0

     vertices = make([]ebiten.Vertex,3*MAXSEG)
     indices = make([]uint16,3*MAXSEG)

     cr := float32(col.R)/0xff
     cg := float32(col.G)/0xff
     cb := float32(col.B)/0xff
     ca := float32(col.A)/0xff

     var segment uint16
     for segment=0;segment<MAXSEG;segment+=1 {
          dx := math.Sin(float64(segment) * 2.0 * math.Pi / MAXSEG) * r
          dy := math.Cos(float64(segment) * 2.0 * math.Pi / MAXSEG) * r
          dxn := math.Sin(float64(segment+1) * 2.0 * math.Pi/ MAXSEG) * r
          dyn := math.Cos(float64(segment+1) * 2.0 * math.Pi / MAXSEG) * r

          vertices[segment*3+0] = ebiten.Vertex{DstX:float32(px),DstY:float32(py), SrcX:1, SrcY:1, ColorR:cr, ColorG:cg, ColorB:cb, ColorA: ca}
          vertices[segment*3+1] = ebiten.Vertex{DstX:float32(px+dx),DstY:float32(py+dy), SrcX:1, SrcY:1, ColorR:cr, ColorG:cg, ColorB:cb, ColorA: ca}
          vertices[segment*3+2] = ebiten.Vertex{DstX:float32(px+dxn),DstY:float32(py+dyn), SrcX:1, SrcY:1, ColorR:cr, ColorG:cg, ColorB:cb, ColorA: ca}

          indices[segment*3+0] = segment*3+0
          indices[segment*3+1] = segment*3+1
          indices[segment*3+2] = segment*3+2
     }

     return vertices,indices
}
