package main

import (
	"math"
	//	"fmt"
)

const epsilon = 1e-8

type Point struct {
	x float64
	y float64
}

func (p Point) length() float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y)
}

func (p Point) normalized() Point {
	l := p.length()
	return Point{x: p.x / l, y: p.y / l}
}

func (p Point) add(t Point) Point {
	return Point{x: p.x + t.x, y: p.y + t.y}
}

func (p Point) sub(t Point) Point {
	return Point{x: p.x - t.x, y: p.y - t.y}
}

func (p Point) rot90() Point {
	return Point{x: p.y, y: -p.x}
}

func (a Point) dot(b Point) float64 {
	return a.x*b.x + a.y*b.y
}

func (p Point) mul(l float64) Point {
	return Point{x: p.x * l, y: p.y * l}
}

func (p Point) reflect(normal Point) Point {
	return p.sub(normal.mul(2.0 * p.dot(normal)))
}

type Line struct {
	src Point
	dst Point
}

type Polygon struct {
	pts []Point
}

func (p Polygon) translated(t Point) Polygon {
	pts := make([]Point, len(p.pts))
	for i := range p.pts {
		pts[i] = p.pts[i].add(t)
	}
	return Polygon{pts: pts}
}

func makeNGon(center Point, radius float64, segments int) Polygon {
	poly := Polygon{pts: make([]Point, segments)}
	for i := 0; i < segments; i++ {
		x := center.x + radius*math.Sin(float64(i)*2.0*math.Pi/float64(segments))
		y := center.y + radius*math.Cos(float64(i)*2.0*math.Pi/float64(segments))
		poly.pts[i] = Point{x: x, y: y}
	}
	return poly
}

func distance(p1 Point, p2 Point) float64 {
	return math.Sqrt((p2.x-p1.x)*(p2.x-p1.x) + (p2.y-p1.y)*(p2.y-p1.y))
}

// note: only returns the closest interesction point to p1
func intersectPolygon(p1 Point, p2 Point, p Polygon) *Point {
	var closestInter *Point
	for i := 0; i < len(p.pts); i++ {
		pa := p.pts[i]
		pb := p.pts[(i+1)%len(p.pts)]

		inter := intersection(p1, p2, pa, pb)
		if inter != nil {
			if closestInter == nil {
				closestInter = inter
			} else {
				if distance(p1, *inter) < distance(p1, *closestInter) {
					closestInter = inter
				}
			}
		}
	}
	return closestInter
}

// note: only returns the closest interesction point to p1
func intersectPolygonNorm(p1 Point, p2 Point, p Polygon) *Line {
	var closestInter *Line
	// fmt.Println("--")
	for i := 0; i < len(p.pts); i++ {
		pa := p.pts[i]
		pb := p.pts[(i+1)%len(p.pts)]

		inter := intersection(p1, p2, pa, pb)
		normal := pb.sub(pa).normalized().rot90()
		if inter != nil {
			if closestInter == nil {
				closestInter = &Line{src: *inter, dst: normal}
				// fmt.Println(i, distance(p1,*inter))
			} else {
				// fmt.Println(i, distance(p1,*inter), distance(p1,closestInter.src) )
				if distance(p1, *inter) < distance(p1, closestInter.src) {
					closestInter = &Line{src: *inter, dst: normal}
				}
			}
		}
	}
	return closestInter
}

// find the intersection point of the lines between p1--p2 and p3--p4
func intersection(p1 Point, p2 Point, p3 Point, p4 Point) *Point {
	// fmt.Println("intersection",p1,p2,p3,p4)
	d := (p1.x-p2.x)*(p3.y-p4.y) - (p1.y-p2.y)*(p3.x-p4.x)

	// fmt.Println(d)

	// If d is zero, there is no intersection
	if d == 0 {
		return nil
	}

	// Get the x and y
	pre := (p1.x*p2.y - p1.y*p2.x)
	post := (p3.x*p4.y - p3.y*p4.x)
	x := (pre*(p3.x-p4.x) - (p1.x-p2.x)*post) / d
	y := (pre*(p3.y-p4.y) - (p1.y-p2.y)*post) / d

	// fmt.Println(pre,post,x,y)


	// Check if the x and y coordinates are within both lines
	if x < math.Min(p1.x, p2.x)-epsilon || x > math.Max(p1.x, p2.x)+epsilon || x < math.Min(p3.x, p4.x)-epsilon || x > math.Max(p3.x, p4.x)+epsilon {
		return nil
	}
	if y < math.Min(p1.y, p2.y)-epsilon || y > math.Max(p1.y, p2.y)+epsilon || y < math.Min(p3.y, p4.y)-epsilon || y > math.Max(p3.y, p4.y)+epsilon {
		return nil
	}

	return &Point{x: x, y: y}
}
