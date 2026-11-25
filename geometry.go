// Package dxfgeometry provides zero-dependency 2D CAD and DXF geometry utilities:
// polyline length, shoelace polygon area, bounding boxes, DPI-to-millimeter
// calibration, architectural drawing scale conversion, Ramer-Douglas-Peucker
// polyline simplification, and minimal ASCII DXF LWPOLYLINE serialization.
package dxfgeometry

import (
	"fmt"
	"math"
	"strings"
)

// Point represents a 2D Cartesian coordinate in drawing or model space.
type Point struct {
	X float64
	Y float64
}

// BBox represents an axis-aligned 2D bounding box.
type BBox struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// Width returns the horizontal extent of the bounding box.
func (b BBox) Width() float64 {
	return b.MaxX - b.MinX
}

// Height returns the vertical extent of the bounding box.
func (b BBox) Height() float64 {
	return b.MaxY - b.MinY
}

// Distance returns the Euclidean distance between two 2D points.
func Distance(a, b Point) float64 {
	return math.Hypot(b.X-a.X, b.Y-a.Y)
}

// PolylineLength returns the total length of an open or closed polyline.
func PolylineLength(pts []Point, closed bool) float64 {
	if len(pts) < 2 {
		return 0
	}
	total := 0.0
	for i := 1; i < len(pts); i++ {
		total += Distance(pts[i-1], pts[i])
	}
	if closed && len(pts) > 2 {
		total += Distance(pts[len(pts)-1], pts[0])
	}
	return total
}

// PolygonArea computes the unsigned area of a simple polygon using the shoelace formula.
func PolygonArea(pts []Point) float64 {
	if len(pts) < 3 {
		return 0
	}
	sum := 0.0
	n := len(pts)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += pts[i].X*pts[j].Y - pts[j].X*pts[i].Y
	}
	return math.Abs(sum) * 0.5
}

// BoundingBox computes the axis-aligned bounding box of a point slice.
// The boolean return value is false when pts is empty.
func BoundingBox(pts []Point) (BBox, bool) {
	if len(pts) == 0 {
		return BBox{}, false
	}
	b := BBox{MinX: pts[0].X, MinY: pts[0].Y, MaxX: pts[0].X, MaxY: pts[0].Y}
	for _, p := range pts[1:] {
		if p.X < b.MinX {
			b.MinX = p.X
		}
		if p.Y < b.MinY {
			b.MinY = p.Y
		}
		if p.X > b.MaxX {
			b.MaxX = p.X
		}
		if p.Y > b.MaxY {
			b.MaxY = p.Y
		}
	}
	return b, true
}

// PixelsToMM converts a raster pixel span at the given DPI to millimeters.
// Returns 0 when dpi <= 0.
func PixelsToMM(pixels, dpi float64) float64 {
	if dpi <= 0 {
		return 0
	}
	return (pixels / dpi) * 25.4
}

// PaperToModelScale converts a paper-space length to real-world model length
// using a 1:scaleDenominator architectural or engineering scale (e.g., 50 for 1:50).
func PaperToModelScale(paperLength, scaleDenominator float64) float64 {
	if scaleDenominator <= 0 {
		return 0
	}
	return paperLength * scaleDenominator
}

func perpendicularDistance(p, a, b Point) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	if dx == 0 && dy == 0 {
		return Distance(p, a)
	}
	num := math.Abs(dy*p.X - dx*p.Y + b.X*a.Y - b.Y*a.X)
	return num / math.Hypot(dx, dy)
}

// SimplifyRDP simplifies a 2D polyline using the Ramer-Douglas-Peucker algorithm.
func SimplifyRDP(pts []Point, epsilon float64) []Point {
	if len(pts) <= 2 || epsilon <= 0 {
		out := make([]Point, len(pts))
		copy(out, pts)
		return out
	}
	maxDist := -1.0
	index := 0
	last := len(pts) - 1
	for i := 1; i < last; i++ {
		d := perpendicularDistance(pts[i], pts[0], pts[last])
		if d > maxDist {
			maxDist = d
			index = i
		}
	}
	if maxDist > epsilon {
		left := SimplifyRDP(pts[:index+1], epsilon)
		right := SimplifyRDP(pts[index:], epsilon)
		res := make([]Point, 0, len(left)+len(right)-1)
		res = append(res, left[:len(left)-1]...)
		res = append(res, right...)
		return res
	}
	return []Point{pts[0], pts[last]}
}

// FormatDXFPolyline formats a 2D polyline as a minimal ASCII DXF LWPOLYLINE entity block.
func FormatDXFPolyline(pts []Point, layer string, closed bool) string {
	if layer == "" {
		layer = "0"
	}
	flag := 0
	if closed {
		flag = 1
	}
	var sb strings.Builder
	sb.WriteString("0\nSECTION\n2\nENTITIES\n")
	sb.WriteString("0\nLWPOLYLINE\n")
	sb.WriteString(fmt.Sprintf("8\n%s\n", layer))
	sb.WriteString(fmt.Sprintf("90\n%d\n", len(pts)))
	sb.WriteString(fmt.Sprintf("70\n%d\n", flag))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("10\n%.4f\n20\n%.4f\n", p.X, p.Y))
	}
	sb.WriteString("0\nENDSEC\n0\nEOF\n")
	return sb.String()
}
