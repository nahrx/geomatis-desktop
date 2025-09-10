package util

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"geomatis-desktop/types"

	"github.com/rwcarlsen/goexif/exif"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}
func LW(a types.Dimension) types.Dimension { //Correct Length Width Dimension
	d := a
	if a.Width > a.Length {
		d.Length = a.Width
		d.Width = a.Length
	}
	return d
}
func GetFeatureDimensions(d types.Diagonal, margin float64) (types.Dimension, types.Dimension) {

	w := math.Pow(math.Pow(d.TopLeft[0]-d.TopRight[0], 2)+math.Pow(d.TopLeft[1]-d.TopRight[1], 2), 0.5)     // pythagoras formula
	h := math.Pow(math.Pow(d.TopLeft[0]-d.BottomLeft[0], 2)+math.Pow(d.TopLeft[1]-d.BottomLeft[1], 2), 0.5) // pythagoras formula
	featureImg := types.Dimension{
		Length: (1 / (1 + margin)) * w,
		Width:  (1 / (1 + margin)) * h,
	}
	marginAroundFeature := types.Dimension{
		Length: (w - featureImg.Length) / 2,
		Width:  (h - featureImg.Width) / 2,
	}

	return featureImg, marginAroundFeature
}
func DimRatio(dim types.Dimension) float64 {
	return dim.Length / dim.Width
}
func calculateCentroid(points []types.Coord) types.Coord {
	var sumX, sumY float64

	// Calculate the sum of all x and y coordinates
	for _, p := range points {
		sumX += p[0]
		sumY += p[1]
	}

	// Calculate the centroid by averaging the x and y coordinates
	centroidX := sumX / float64(len(points))
	centroidY := sumY / float64(len(points))

	return types.Coord{centroidX, centroidY}
}
func CalculateGeoreferenceParameters(img types.Dimension, rasterPoints []types.Coord, extent types.Extent, margin float64, imageRotation float64) *types.WorldFileParameter {
	// WorldFileParameter : A, D, B, E, C, F
	var scale, x, y float64

	deltaX := extent.MaxX - extent.MinX
	deltaY := extent.MaxY - extent.MinY
	var polygon types.Dimension
	polygon.Length, polygon.Width = deltaX, deltaY

	diagonal, _ := FindDiagonalPoints(rasterPoints)
	featureImg, _ := GetFeatureDimensions(diagonal, margin)

	angle := CalculateRotationAngle(diagonal.TopLeft, diagonal.TopRight)

	//if the polygon is landscape, but the image is portrait
	if deltaX > deltaY && img.Length < img.Width {
		angle = angle + 90
		fmt.Println("angle = angle + 90")
	}
	//if the polygon is portrait, but the image is landscape
	if deltaX < deltaY && img.Length > img.Width {
		angle = angle - 90
		fmt.Println("angle = angle - 90")
	}

	/*if the image is portrait and rotated 90 degrees (counter clockwise), rotate it 180 degrees*/
	/*if the image is portrait and rotated -90 degrees (clockwise), ignore it*/
	if img.Length < img.Width && imageRotation == 90 {
		angle = angle + 180
	}
	/*if the image is landscape and rotated 90 degrees (counter clockwise), ignore it*/
	/*if the image is landscape and rotated -90 degrees (clockwise), rotate it 180 degrees*/
	if img.Length > img.Width && imageRotation == -90 {
		angle = angle + 180
	}

	// if the image is rotated 180 degrees
	if imageRotation == 180 {
		angle = angle + 180
	}

	LWPolygon := LW(polygon)
	LWFeatureImg := LW(featureImg)

	radian := math.Pi * angle / 180.0

	fmt.Println("image rotation : ", imageRotation)

	if DimRatio(LWPolygon) >= DimRatio(LWFeatureImg) { // horizontally tilted
		//Scale X
		scale = LWPolygon.Length / LWFeatureImg.Length

	} else {
		//Scale Y
		scale = LWPolygon.Width / LWFeatureImg.Width
	}
	polygonPoint := []types.Coord{
		{extent.MinX, extent.MinY},
		{extent.MinX, extent.MaxY},
		{extent.MaxX, extent.MinY},
		{extent.MaxX, extent.MaxY},
	}

	featureImgCentroid := calculateCentroid(rasterPoints)
	polygonCentroid := calculateCentroid(polygonPoint)

	x = featureImgCentroid[0]
	y = featureImgCentroid[1]

	// Calculate world file parameters

	var p types.WorldFileParameter
	p.A = scale * math.Cos(radian)
	p.D = scale * math.Sin(radian)
	p.B = scale * math.Sin(radian)
	p.E = -scale * math.Cos(radian)
	p.C = polygonCentroid[0] - p.A*(x) - p.B*(y)
	p.F = polygonCentroid[1] - p.D*(x) - p.E*(y)
	return &p
}
func GetImageDimensions(file io.Reader) (types.Dimension, error) {
	//time.Sleep(2 * time.Second)
	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println(err.Error())
		return types.Dimension{}, fmt.Errorf("Error decoding image config: %s", err)
	}
	return types.Dimension{
		Length: float64(imgConfig.Width),
		Width:  float64(imgConfig.Height),
	}, nil
}
func GeRotationDegree(filepath string) (float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, nil
	}
	defer file.Close()

	var rotation float64
	orientationTag, err := GetOrientationTag(file)

	if err != nil {
		return 0, fmt.Errorf("Error GetOrientationTag : %w", err)
	}
	/*In standard mathematical convention,
	positive angles are measured counterclockwise (anticlockwise) from the positive x-axis,
	while negative angles are measured clockwise.
	*/
	if orientationTag == 8 {
		rotation = -90
	} else if orientationTag == 3 {
		rotation = 180
	} else if orientationTag == 6 {
		rotation = 90
	}
	return rotation, nil
}
func GetOrientationTag(file io.Reader) (int, error) {
	//time.Sleep(2 * time.Second)

	// decode file exif
	x, err := exif.Decode(file)
	if err != nil {
		if err == io.EOF {
			return 1, nil
		}
		return 0, fmt.Errorf("Error exif.Decode : %s", err)
	}

	o, err := x.Get(exif.Orientation) // normally, don't ignore errors!
	if err != nil {
		var tagNotPresentError exif.TagNotPresentError = "Orientation"
		if err.Error() == tagNotPresentError.Error() {
			return 0, nil
		}
		return 0, fmt.Errorf("Error x.Get : %s", err)
	}
	orientation, err := o.Int(0)
	if err != nil {
		return 0, fmt.Errorf("Error o.Int : %s", err)
	}
	return orientation, nil
}

func GetOrientedImageDimensions(file1 io.Reader, file2 io.Reader) (types.Dimension, error) {
	imgDim, err := GetImageDimensions(file1)
	if err != nil {
		return types.Dimension{}, fmt.Errorf("Error GetImageDimensions : %w", err)
	}
	oValue, err := GetOrientationTag(file2)
	fmt.Println("orientation value : ", oValue)
	fmt.Println("orientation before : ", imgDim)
	if err != nil {
		return types.Dimension{}, fmt.Errorf("Error GetOrientationTag : %w", err)
	}
	if oValue == 8 || oValue == 6 { // Swap width and height if rotation 90 or -90 degree
		imgDim.Length, imgDim.Width = imgDim.Width, imgDim.Length
	}
	return imgDim, nil
}

func WriteWorldFileParametersToFile(filePath string, p types.WorldFileParameter) error {
	content := fmt.Sprintf("%.20f\n%.20f\n%.20f\n%.20f\n%.20f\n%.20f\n", p.A, p.D, p.B, p.E, p.C, p.F)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetRasterFeaturePoints(filePath string, orientation float64) ([]types.Coord, error) {
	safePath := strings.ReplaceAll(filePath, `\`, `/`) // Replace all backslashes with forward slashes
	pythonCode := fmt.Sprintf(`import pypy; print(pypy.rasterFeaturePoints('%s', False,%v))`, safePath, orientation)
	cmd := exec.Command("python", "-c", pythonCode)
	//cmd := exec.Command("build/bin/python-embed/python.exe", "-c", pythonCode)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	rasterFeaturePoints, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Failed to call python function. error : %s.", err.Error())
	}

	var points types.FeaturePoints
	err = json.Unmarshal(rasterFeaturePoints, &points)
	if err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal rasterFeaturePoints. error : %s.", err.Error())
	}
	return points.Points, nil
}

func SaveFile(filePath string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("Failed to open file. error : %s.", err.Error())
	}
	defer file.Close()

	// Create a new file in the server's "uploads" directory

	newFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create file. error : %s.", err.Error())
	}
	defer newFile.Close()

	// Copy the uploaded file's contents to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		return fmt.Errorf("Failed to copy file contents. error : %s.", err.Error())
	}
	return nil
}
func FileNameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(path.Ext(fileName))]
}
func AllNotNil(data ...any) bool {
	for _, d := range data {
		if d == nil || d == "" {
			return false
		}
	}
	return true
}

func CalculateRotationAngle(point1, point2 types.Coord) float64 {
	// Calculate the differences in x and y coordinates
	x1, y1 := point1[0], point1[1]
	x2, y2 := point2[0], point2[1]

	dx := x2 - x1
	dy := y2 - y1

	// Calculate the angle using the arctangent function (Atan2)
	// Note that Atan2 returns the angle in radians, so we convert it to degrees.
	angleRad := math.Atan2(dy, dx)
	angleDeg := angleRad * 180.0 / math.Pi

	// Ensure the angle is between 0 and 360 degrees
	if angleDeg < 0 {
		angleDeg += 360.0
	}

	return angleDeg
}

func FindDiagonalPoints(rectangle []types.Coord) (types.Diagonal, error) {
	if len(rectangle) != 4 {
		return types.Diagonal{}, fmt.Errorf("rectangle must have exactly 4 points")
	}

	// Find the two points with the smallest sum of X and Y coordinates (top-left and bottom-right)
	topLeft, bottomRight := rectangle[0], rectangle[0]
	// Find the two points with the smallest difference of X and Y coordinates (top-right and bottom-left)
	topRight, bottomLeft := rectangle[0], rectangle[0]
	for _, point := range rectangle {
		if point[0]+point[1] < bottomLeft[0]+bottomLeft[1] {
			bottomLeft = point //bottomLeft
		}
		if point[0]+point[1] > topRight[0]+topRight[1] {
			topRight = point //topRight
		}

		if point[0]-point[1] < topLeft[0]-topLeft[1] {
			topLeft = point //topLeft
		}
		if point[0]-point[1] > bottomRight[0]-bottomRight[1] {
			bottomRight = point //bottomRight
		}
	}

	// The diagonal points
	diagonal := types.Diagonal{
		TopLeft:     bottomLeft,
		TopRight:    bottomRight,
		BottomLeft:  topLeft,
		BottomRight: topRight,
	}
	return diagonal, nil
}
