package storage

import (
	"geomatis-desktop/bpsmap"
	"geomatis-desktop/types"
)

type Storage interface {
	TableExist(string) (bool, error)
	MasterMapExist(string) (bool, error)
	MasterMapAttributeExist(string, string) (bool, error)
	GetMasterMaps() ([]types.MasterMap, error)
	GetMasterMapByName(string) (types.MasterMap, error)
	GetMasterMapAttributes(string) ([]types.MasterMapAttr, error)
	GetExtent(string, string, bpsmap.BpsMap) (*types.Extent, error)
	GetAttributesValue(string, string, string, []string) ([]string, error)
	CreateMasterMaps(string, *[]byte) error
	DeleteMasterMap(string) error
	Close() error
}
