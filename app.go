package main

import (
	"context"
	"encoding/json"
	"fmt"
	"geomatis-desktop/geo"
	"geomatis-desktop/storage"
	"geomatis-desktop/types"
	"geomatis-desktop/util"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	store  storage.Storage
	extent *geo.Extents
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

//DATABASE ################################################################################

func (a *App) LoadDbConfig() (*storage.Config, error) {
	fmt.Println("load DB config")
	file, err := os.Open("config.json")
	if err != nil {
		if err == os.ErrNotExist {
			return nil, fmt.Errorf("there is no database config file yet")
		}
		return nil, fmt.Errorf("unable to open config.json: %w", err)
	}
	defer file.Close()

	var config storage.Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return &config, nil
}

func (a *App) ConnectToDB(config storage.Config) error {
	s, err := storage.NewPostgreStorage(config)
	if err != nil {
		return fmt.Errorf("failed to connect : %w", err)
	}
	a.store = s
	return nil
}

func (a *App) DisconnectDB(config storage.Config) error {
	if a.store != nil {
		err := a.store.Close()
		if err != nil {
			return fmt.Errorf("failed to close DB connection: %w", err)
		}
		a.store = nil // avoid reusing a closed DB
	}
	return nil
}

// SaveConfig saves the user input to config.json
func (a *App) SaveDbConfig(config storage.Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("failed to create config.json: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// MASTER MAPS ##############################################################
func (a *App) GetMasterMaps() ([]types.MasterMap, error) {
	if a.store == nil {
		return nil, fmt.Errorf("database not connected")
	}
	data, err := a.store.GetMasterMaps()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *App) CreateMasterMaps(filePaths string) error {
	if a.store == nil {
		return fmt.Errorf("database not connected")
	}
	file, err := os.Open(filePaths)
	if err != nil {
		return fmt.Errorf("error opening file :  %w ", err)
	}
	defer file.Close()
	// file validation
	fileName := file.Name()

	allowedExt := ".geojson"
	if path.Ext(fileName) != allowedExt {
		return fmt.Errorf("The uploaded file must have the following extensions : %s ", allowedExt)
	}

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file stat :  %w ", err)
	}

	var maxFileSize int64 = 200_000_000
	if fileStat.Size() > maxFileSize {
		return fmt.Errorf("The uploaded file cannot be larger than %v", maxFileSize)
	}

	// file content processing
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error file content processing. error : %s", err.Error())

	}

	//modifiedString := strings.ReplaceAll(fileName, " ", "_")
	name := util.FileNameWithoutExtension(fileStat.Name())

	err = a.store.CreateMasterMaps(name, &fileBytes)
	if err != nil {
		return fmt.Errorf("error when storing master maps. error : %s", err.Error())

	}
	return nil
}
func (a *App) DeleteMasterMap(masterMap string) error {
	if a.store == nil {
		return fmt.Errorf("database not connected")
	}
	if !util.AllNotNil(masterMap) {
		return fmt.Errorf("request parameter not found. master_maps needed in the request.")
	}
	err := a.store.DeleteMasterMap(masterMap)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) SelectGeojsonFile() (string, error) {
	if a.store == nil {
		return "", fmt.Errorf("database not connected")
	}
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Geojson File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Geojson Files",
				Pattern:     "*.geojson",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

// GEOREFERENCE #########################################################################
func (a *App) GetExtent(idsls string) (*types.Extent, error) {
	if extent, ok := (*a.extent)[idsls]; ok {
		return &extent, nil
	} else {
		return nil, fmt.Errorf("Key %s not found\n", idsls)
	}
}
func (a *App) SelectGeojsonFileForGeoreference(mapType string) (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Geojson File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Geojson Files",
				Pattern:     "*.geojson",
			},
		},
	})
	if err != nil {
		return "", err
	}
	if mapType == "ws" {
		e, err := geo.ReadExtents(file, geo.WsMap{})
		if err != nil {
			return "", err
		}
		a.extent = e
	} else if mapType == "wb" {
		e, err := geo.ReadExtents(file, geo.WbMap{})
		if err != nil {
			return "", err
		}
		a.extent = e
	} else {
		return "", fmt.Errorf("Error map type")
	}

	return file, nil
}
func (a *App) SelectFiles() ([]string, error) {
	file, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Image File (jpg/jpeg or png only)",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Image Files",
				Pattern:     "*.jpg;*.jpeg;*.png",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (a *App) ProcessFiles(filePaths []string, masterMap string, masterMapType string, masterMapSource string) ([]string, error) {
	log := []string{}
	if masterMap == "" {
		return nil, fmt.Errorf("missing master map")
	}
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("missing raster file")
	}

	gSettings := types.GeoreferenceSettings{}
	gSettings.Prepare()
	gSettings.MasterMap = masterMap

	if masterMapSource == "database" {
		if a.store == nil {
			return nil, fmt.Errorf("database not connected")
		}

		masterMapExist, err := a.store.MasterMapExist(masterMap)
		if err != nil {
			return nil, fmt.Errorf("Error when calling MasterMapExist. Error :  %s", err.Error())
		}
		if !masterMapExist {
			return nil, fmt.Errorf("%s is not found in the database. Error :  %s", masterMap, err.Error())
		}
		attrKeyExist, err := a.store.MasterMapAttributeExist(masterMap, gSettings.AttrKey)
		if err != nil {
			return nil, fmt.Errorf("Error when calling MasterMapExist. Error :  %s", err.Error())
		}
		if !attrKeyExist {
			return nil, fmt.Errorf("%s is not found in the database. Error :  %s", masterMap, err.Error())
		}

	} else if masterMapSource == "file" {
		if a.extent == nil {
			return nil, fmt.Errorf("geojson file is not ready")
		}
	} else {
		return nil, fmt.Errorf("masterMapSource not found")
	}

	//step : check if master map exist in database
	//step : prepare georeference setting

	gSettings.MasterMapSource = masterMapSource

	var mapType geo.BpsMap
	if masterMapType == "ws" {
		mapType = geo.WsMap{}
	} else if masterMapType == "wb" {
		mapType = geo.WsMap{}
	} else {
		return nil, fmt.Errorf("Error select mapType")
	}

	numJobs := len(filePaths)
	numWorkers := 20
	if numJobs < numWorkers {
		numWorkers = numJobs
	}
	files := make(chan string, numJobs)
	results := make(chan types.Result, numJobs)
	for w := 0; w < numWorkers; w++ {
		go a.worker(w, files, results, gSettings, mapType)
	}

	for j := 0; j < numJobs; j++ {
		files <- filePaths[j]
	}
	close(files)

	success, fail := 0, 0
	for a := 0; a < numJobs; a++ {
		r := <-results
		if r.Error != nil {
			fail++
			log = append(log, fmt.Sprintf("failed %s : %s.", filepath.Base(r.Id), r.Error.Error()))
			continue
		}
		log = append(log, fmt.Sprintf("success : %s", filepath.Base(r.Id)))
	}
	success = numJobs - fail
	//return success, fail, e

	fmt.Printf("%+v", filePaths)
	fmt.Println(success, fail)
	return log, nil
}
