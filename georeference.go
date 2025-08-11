package main

import (
	"fmt"
	"geomatis-desktop/bpsmap"
	"geomatis-desktop/util"
	"os"
	"path"
	"path/filepath"
	"strings"

	"geomatis-desktop/types"
)

func (a *App) GeoreferenceWorker(id int, rasterFilePath <-chan string, results chan<- types.Result, g types.GeoreferenceSettings, mapType bpsmap.BpsMap) {
	fmt.Println("worker : ", id)
	for rasterPath := range rasterFilePath {

		result := types.Result{
			Id:    path.Base(rasterPath),
			Error: nil,
		}
		//Get image dimension
		file1, err := os.Open(rasterPath)
		if err != nil {
			result.Error = fmt.Errorf("Error while converting multipart FileHeader to File. error : %s.", err.Error())
			results <- result
			continue
		}
		defer file1.Close()

		orientation := false
		imgDim := types.Dimension{}
		if orientation {
			file2, err := os.Open(rasterPath)
			if err != nil {
				result.Error = fmt.Errorf("Error while converting multipart FileHeader to File. error : %s.", err.Error())
				results <- result
				continue
			}
			defer file2.Close()

			imgDim, err = util.GetOrientedImageDimensions(file1, file2)
			if err != nil {
				result.Error = fmt.Errorf("Error GetOrientedImageDimensions : %w", err)
				results <- result
				continue
			}
		} else {
			imgDim, err = util.GetImageDimensions(file1)
			if err != nil {
				result.Error = fmt.Errorf("Error GetImageDimensions : %w", err)
				results <- result
				continue
			}
		}

		//Get raster key
		rasterKey, err := GetRasterKey(rasterPath, g.RasterKeySettings)
		if err != nil {
			result.Error = fmt.Errorf("Error GetRasterKey: %s.", err.Error())
			results <- result
			continue
		}

		//Get polygon extent, raster feature point from image
		polygonExtent := &types.Extent{}
		if g.MasterMapSource == "database" {
			polygonExtent, err = a.store.GetExtent(g.MasterMap, rasterKey, mapType)
		} else {
			polygonExtent, err = a.GetExtent(rasterKey)
		}
		if err != nil {
			result.Error = fmt.Errorf("Error GetExtent : %s.", err.Error())
			results <- result
			continue
		}

		//need to edit the featurePoints variable back to the normal orientationtag(value = 1)
		featurePoints := []types.Coord{}
		if orientation {
			featurePoints, err = util.GetRasterFeaturePoints(rasterPath)
			if err != nil {
				result.Error = fmt.Errorf("Error GetRasterFeaturePoints : %s.", err.Error())
				results <- result
				continue
			}
		} else {
			featurePoints, err = util.GetOrientationRemovedRasterFeaturePoints(rasterPath)
			if err != nil {
				result.Error = fmt.Errorf("Error GetRasterFeaturePoints : %s.", err.Error())
				results <- result
				continue
			}
		}

		//Calculate Georeference Parameter and save world file
		parameter := util.CalculateGeoreferenceParameters(imgDim, featurePoints, *polygonExtent, g.RasterFeatureSettings.Margin)
		worldFileExt := GetWorldFileExtlist()[strings.ToLower(path.Ext(rasterPath))]
		fmt.Println("worldFileExt : ", worldFileExt)
		worldFileName := fmt.Sprintf("%s%s", util.FileNameWithoutExtension(rasterPath), worldFileExt)
		err = util.WriteWorldFileParametersToFile(worldFileName, *parameter)
		if err != nil {
			result.Error = fmt.Errorf("Error while creating worldfile. error : %s.", err.Error())
			results <- result
			continue
		}
		results <- result
	}
}

func GetWorldFileExtlist() map[string]string {
	return map[string]string{
		".jpg":  ".jgw",
		".jpeg": ".jgw",
		".png":  ".pgw",
	}
}

func GetRasterKey(filename string, rasterKeySettings *types.RasterKeySettings) (string, error) {
	filename = util.FileNameWithoutExtension(filepath.Base(filename))
	if len(filename) < rasterKeySettings.NumChar {
		return "", fmt.Errorf("Length of string is not enough (less than rasterKeySettings.NumChar)")
	}
	switch rasterKeySettings.Type {
	case "all":
		return filename, nil
	case "prefix":
		return filename[:rasterKeySettings.NumChar], nil
	case "suffix":
		return filename[len(filename)-rasterKeySettings.NumChar:], nil
	case "regex":
		return rasterKeySettings.Regex.FindString(filename), nil
	}
	return "", fmt.Errorf("Type of raster key is not valid. Only all, prefix, suffix, or regex allowed.")
}
