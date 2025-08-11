package bpsmap

import (
	"encoding/json"
	"fmt"
	"geomatis-desktop/types"
	"os"
)

type Extents map[string]types.Extent

type FeatureCollection struct {
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	CRS      CRS       `json:"crs"`
	Features []Feature `json:"features"`
}

type CRS struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	Idsls string `json:"idsls"`
	Idbs  string `json:"idbs"`
}

type Geometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

type Geometries map[string][][]float64

func ParseExtents(filePath string, bpsMap BpsMap) (*Extents, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc FeatureCollection
	if err := json.Unmarshal(file, &fc); err != nil {
		return nil, fmt.Errorf("failed to parse FeatureCollection GeoJSON: %w", err)
	}

	extents := make(Extents)

	for _, feature := range fc.Features {
		var multiPolygon [][][][]float64
		if err := json.Unmarshal(feature.Geometry.Coordinates, &multiPolygon); err != nil {
			return nil, fmt.Errorf("Error due to geometry error: %w\n", err)
		}

		if len(multiPolygon) > 0 && len(multiPolygon[0]) > 0 {
			outerRing := multiPolygon[0][0]
			extentKey, err := bpsMap.GetExtentKey(feature.Properties)
			if err != nil {
				return nil, err
			}
			extents[extentKey] = computeExtent(outerRing)
		}
	}
	return &extents, nil
}

func computeExtent(coords [][]float64) types.Extent {
	if len(coords) == 0 {
		return types.Extent{}
	}
	minX, minY := coords[0][0], coords[0][1]
	maxX, maxY := coords[0][0], coords[0][1]

	for _, pt := range coords {
		x, y := pt[0], pt[1]
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return types.Extent{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
}
