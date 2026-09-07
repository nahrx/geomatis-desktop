package bpsmap

import (
	"fmt"
	"geomatis-desktop/types"
)

type BpsMap interface {
	GetExtentKey(Properties) (string, error)
	GetKeyName() string
	GetGeoreferenceSetting(types.GeoreferenceSettings) types.GeoreferenceSettings
}

type WssMap struct{}

type WsMap struct{}

type WbMap struct{}

func (wss WssMap) GetExtentKey(properties Properties) (string, error) {
	if properties.Idsubsls == "" {
		return "", fmt.Errorf("Idsubsls not found. check the selected map type (WSS/WS/WB)")
	}
	return properties.Idsubsls, nil
}

func (ws WsMap) GetExtentKey(properties Properties) (string, error) {
	if properties.Idsls == "" {
		return "", fmt.Errorf("Idsls not found. check the selected map type (WSS/WS/WB)")
	}
	return properties.Idsls, nil
}

func (wb WbMap) GetExtentKey(properties Properties) (string, error) {
	if properties.Idbs == "" {
		return "", fmt.Errorf("Idbs not found. check the selected map type (WSS/WS/WB)")
	}
	return properties.Idbs, nil
}

func (wss WssMap) GetGeoreferenceSetting(gSetting types.GeoreferenceSettings)types.GeoreferenceSettings{
	if gSetting.RasterKeySettings == nil {
		gSetting.RasterKeySettings = &types.RasterKeySettings{}
	}
	gSetting.AttrKey = "idsubsls"
	gSetting.RasterKeySettings.Category = "prefix"
	gSetting.RasterKeySettings.NumChar = 16
	return gSetting
}

func (ws WsMap) GetGeoreferenceSetting(gSetting types.GeoreferenceSettings)types.GeoreferenceSettings{
	if gSetting.RasterKeySettings == nil {
		gSetting.RasterKeySettings = &types.RasterKeySettings{}
	}
	gSetting.AttrKey = "idsls"
	gSetting.RasterKeySettings.Category = "prefix"
	gSetting.RasterKeySettings.NumChar = 14
	return gSetting
}

func (wb WbMap) GetGeoreferenceSetting(gSetting types.GeoreferenceSettings)types.GeoreferenceSettings{
	if gSetting.RasterKeySettings == nil {
		gSetting.RasterKeySettings = &types.RasterKeySettings{}
	}
	gSetting.AttrKey = "idbs"
	gSetting.RasterKeySettings.Category = "prefix"
	gSetting.RasterKeySettings.NumChar = 14
	return gSetting
}

func (wss WssMap) GetKeyName() string {
	return "idsubsls"
}

func (ws WsMap) GetKeyName() string {
	return "idsls"
}

func (wb WbMap) GetKeyName() string {
	return "idbs"
}
