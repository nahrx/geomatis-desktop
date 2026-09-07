package types

import (
	"mime/multipart"
	"regexp"
)

type RasterKeySettings struct {
	Category string         `json:"type"`
	NumChar  int            `json:"num_char"`
	Regex    *regexp.Regexp `json:"-"`
}
type RasterFeatureSettings struct {
	XPosition string  `json:"x_position"`
	YPosition string  `json:"y_position"`
	Margin    float64 `json:"margin"`
}
type GeoreferenceSettings struct {
	MasterMapSource   string  `json:"master_map_source"`
	MasterMap         string  `json:"master_map"`
	AttrKey           string  `json:"attr_key"`
	RasterRotation    float64 `json:"raster_rotation"` // 0, 90, -90, 180
	RasterKeySettings *RasterKeySettings `json:"raster_key_settings"`
	// TargetDir             string
	// SeparateDirAttrs      []string
	RasterFeatureSettings *RasterFeatureSettings `json:"raster_feature_settings"`
}

func (g *GeoreferenceSettings) Prepare() {
	// Initialize nested pointers
	g.RasterKeySettings = &RasterKeySettings{}
	g.RasterFeatureSettings = &RasterFeatureSettings{}

	g.AttrKey = "idsubsls"
	g.RasterRotation = 0
	g.RasterKeySettings.Category = "prefix"
	g.RasterKeySettings.NumChar = 16
	g.RasterFeatureSettings.Margin = 0.05 // 5 persen
	g.RasterFeatureSettings.XPosition = "none"
	g.RasterFeatureSettings.YPosition = "none"
}

type GeoreferenceRequest struct {
	Raster   []*multipart.FileHeader
	Settings *GeoreferenceSettings
}

type MasterMap struct {
	Name      string `json:"name"`
	Dimension int    `json:"dimension"`
	Srid      int    `json:"srid"`
	Category  string `json:"type"`
}
type MasterMapAttr struct {
	Name     string `json:"name"`
	Category string `json:"type"`
}
type Dimension struct {
	Length, Width float64
}
type Diagonal struct {
	TopLeft, TopRight, BottomLeft, BottomRight Coord
}
type Margin struct {
	MarginX, MarginY float64
}
type Coord []float64

type FeaturePoints struct {
	Points []Coord `json:"points"`
}

type Extent struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

type WorldFileParameter struct {
	A, D, B, E, C, F float64
}
type Result struct {
	Id    string
	Error error
}
