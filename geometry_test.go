package dxfgeometry

import (
	"math"
	"strings"
	"testing"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestPolylineAndArea(t *testing.T) {
	square := []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}
	if got := PolylineLength(square, false); !almostEqual(got, 30, 1e-9) {
		t.Fatalf("open length = %v, want 30", got)
	}
	if got := PolylineLength(square, true); !almostEqual(got, 40, 1e-9) {
		t.Fatalf("closed length = %v, want 40", got)
	}
	if got := PolygonArea(square); !almostEqual(got, 100, 1e-9) {
		t.Fatalf("area = %v, want 100", got)
	}
}

func TestBoundingBox(t *testing.T) {
	pts := []Point{{-2, 5}, {8, -3}, {4, 12}}
	b, ok := BoundingBox(pts)
	if !ok {
		t.Fatal("expected valid bounding box")
	}
	if b.Width() != 10 || b.Height() != 15 {
		t.Fatalf("bbox dimensions = (%v, %v), want (10, 15)", b.Width(), b.Height())
	}
}

func TestScaleAndDPI(t *testing.T) {
	if got := PixelsToMM(300, 300); !almostEqual(got, 25.4, 1e-9) {
		t.Fatalf("PixelsToMM(300, 300) = %v, want 25.4", got)
	}
	if got := PaperToModelScale(20, 50); !almostEqual(got, 1000, 1e-9) {
		t.Fatalf("PaperToModelScale(20, 50) = %v, want 1000", got)
	}
}

func TestSimplifyRDP(t *testing.T) {
	pts := []Point{{0, 0}, {5, 0.05}, {10, 0}, {10, 10}}
	got := SimplifyRDP(pts, 0.2)
	if len(got) != 3 {
		t.Fatalf("simplified len = %d, want 3 (%v)", len(got), got)
	}
}

func TestFormatDXFPolyline(t *testing.T) {
	dxf := FormatDXFPolyline([]Point{{0, 0}, {10, 0}, {10, 5}}, "WALLS", true)
	if !strings.Contains(dxf, "LWPOLYLINE") || !strings.Contains(dxf, "WALLS") || !strings.Contains(dxf, "EOF") {
		t.Fatalf("unexpected DXF output:\n%s", dxf)
	}
}
