package bpsmap

import "fmt"

type BpsMap interface {
	GetExtentKey(Properties) (string, error)
	GetKeyName() string
}

type WsMap struct{}

type WbMap struct{}

func (ws WsMap) GetExtentKey(properties Properties) (string, error) {
	if properties.Idsls == "" {
		return "", fmt.Errorf("Idsls not found. check the selected map type (WS/WB)")
	}
	return properties.Idsls, nil
}

func (wb WbMap) GetExtentKey(properties Properties) (string, error) {
	if properties.Idbs == "" {
		return "", fmt.Errorf("Idbs not found. check the selected map type (WS/WB)")
	}
	return properties.Idbs, nil
}

func (ws WsMap) GetKeyName() string {
	return "idsls"
}

func (wb WbMap) GetKeyName() string {
	return "idbs"
}
