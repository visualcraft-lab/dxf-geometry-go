# dxf-geometry-go

Zero-dependency Go library for 2D CAD and DXF geometry calculations, raster-to-CAD DPI calibration, architectural scale conversion, Ramer-Douglas-Peucker polyline simplification, and minimal ASCII DXF `LWPOLYLINE` serialization.

## Installation

```bash
go get github.com/visualcraft-lab/dxf-geometry-go@v0.1.0
```

## Quick Start

```go
package main

import (
	"fmt"
	dxf "github.com/visualcraft-lab/dxf-geometry-go"
)

func main() {
	wall := []dxf.Point{{X: 0, Y: 0}, {X: 5000, Y: 0}, {X: 5000, Y: 4000}, {X: 0, Y: 4000}}
	fmt.Println("Perimeter (mm):", dxf.PolylineLength(wall, true))
	fmt.Println("Area (mm^2):", dxf.PolygonArea(wall))
	fmt.Println("300 px @ 300 DPI (mm):", dxf.PixelsToMM(300, 300))
	fmt.Println("20 mm @ 1:50 (mm):", dxf.PaperToModelScale(20, 50))
}
```

## Features

- `Distance`, `PolylineLength`, `PolygonArea`, `BoundingBox`
- `PixelsToMM` for calibrating raster scan dimensions from DPI
- `PaperToModelScale` for standard 1:N architectural and engineering drawing scales
- `SimplifyRDP` (Ramer-Douglas-Peucker) for reducing noisy vectorized polyline vertices
- `FormatDXFPolyline` for emitting minimal ASCII DXF `LWPOLYLINE` entity blocks

## Online Reference

- Interactive blueprint and image-to-CAD converter: [ImgToCAD](https://imgtocad.com/)
- Sketch vectorization reference: [Sketch to CAD](https://imgtocad.com/sketch-to-cad)

## License

MIT License — see [LICENSE](LICENSE).
