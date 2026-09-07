package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"geomatis-desktop/bpsmap"
	"geomatis-desktop/types"

	"github.com/lib/pq"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type PostgreStorage struct {
	Db *sql.DB
}

type Config struct {
	DBHost     string `json:"DB_HOST"`
	DBPort     int    `json:"DB_PORT"`
	DBDatabase string `json:"DB_DATABASE"`
	DBUsername string `json:"DB_USERNAME"`
	DBPassword string `json:"DB_PASSWORD"`
}

func makeSqlScanFunc[T comparable](columns []T) []interface{} {
	columnPointers := make([]interface{}, cap(columns))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}
	return columnPointers
}

func NewPostgreStorage(c Config) (*PostgreStorage, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }
	dbHost := c.DBHost
	dbPort := c.DBPort
	dbDatabase := c.DBDatabase
	dbUsername := c.DBUsername
	dbPassword := c.DBPassword

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbDatabase)
	fmt.Println(connStr)
	//connStr := "postgres://postgres:password@localhost/geomatis?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	err = EnsurePostGISExtension(db)
	if err != nil {
		return nil, fmt.Errorf("PostGIS check/install failed: %w", err)
	}

	return &PostgreStorage{
		Db: db,
	}, nil
}
func EnsurePostGISExtension(db *sql.DB) error {
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'postgis')`

	err := db.QueryRow(checkQuery).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check PostGIS extension: %w", err)
	}

	if !exists {
		_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS postgis`)
		if err != nil {
			return fmt.Errorf("failed to create PostGIS extension: %w", err)
		}
	}
	return nil
}

func (s *PostgreStorage) Close() error {
	err := s.Db.Close()
	if err != nil {
		return fmt.Errorf("Failed to disconnect database : %w", err)
	}
	return nil
}
func (s *PostgreStorage) TableExist(tableName string) (bool, error) {
	// Retrieve table names from the database
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`
	var exists bool
	err := s.Db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (s *PostgreStorage) MasterMapExist(masterMap string) (bool, error) {
	// Retrieve table names from the database
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM geometry_columns
			WHERE f_table_schema = 'public'
			AND f_table_name = $1
		)
	`
	var exists bool
	err := s.Db.QueryRow(query, masterMap).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (s *PostgreStorage) MasterMapAttributeExist(masterMap, attr string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2
	)`
	var exists bool
	err := s.Db.QueryRow(query, masterMap, attr).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *PostgreStorage) GetMasterMaps() ([]types.MasterMap, error) {
	// Retrieve table names from the database
	query, err := s.Db.Query(`
		SELECT f_table_name AS name, coord_dimension AS dimension, srid, type 
		FROM geometry_columns 
		WHERE f_table_schema='public'
	`)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	var values []types.MasterMap
	for query.Next() {
		var v types.MasterMap
		err := query.Scan(&v.Name, &v.Dimension, &v.Srid, &v.Category)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}
	if err := query.Err(); err != nil {
		return nil, err
	}
	// Create the response object
	return values, nil
}
func (s *PostgreStorage) GetMasterMapByName(masterMap string) (types.MasterMap, error) {
	// Retrieve table names from the database
	query := `
			SELECT f_table_name, coord_dimension, srid, type
			FROM geometry_columns
			WHERE f_table_schema = 'public'
			AND f_table_name = $1
	`
	//var v map[string]string
	//err := s.Db.QueryRow(query, masterMap).Scan(&v)
	var v types.MasterMap
	err := s.Db.QueryRow(query, masterMap).Scan(&v.Name, &v.Dimension, &v.Srid, &v.Category)
	if err != nil {
		return types.MasterMap{}, err
	}

	return v, nil
}
func (s *PostgreStorage) GetMasterMapAttributes(masterMap string) ([]types.MasterMapAttr, error) {
	query, err := s.Db.Query(`
			SELECT column_name,udt_name
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE table_schema = 'public' and TABLE_NAME=$1 and udt_name != 'geometry'
			ORDER BY column_name ASC
			`, masterMap)

	if err != nil {
		return nil, err
	}
	defer query.Close()
	var values []types.MasterMapAttr
	for query.Next() {
		var v types.MasterMapAttr
		err := query.Scan(&v.Name, &v.Category)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}
	if err := query.Err(); err != nil {
		return nil, err
	}

	// Create the response object
	return values, nil
}

// getGeometryColumn detects the geometry column name of a table via PostGIS's
// geometry_columns view, instead of assuming a fixed name like "geom".
func (s *PostgreStorage) getGeometryColumn(tableName string) (string, error) {
	var col string
	err := s.Db.QueryRow(`
		SELECT f_geometry_column
		FROM geometry_columns
		WHERE f_table_schema = 'public' AND f_table_name = $1
		LIMIT 1
	`, tableName).Scan(&col)
	if err != nil {
		return "", fmt.Errorf("failed to detect geometry column of table %s: %w", tableName, err)
	}
	return col, nil
}

func (s *PostgreStorage) GetExtent(tableName, idSls string, mapType bpsmap.BpsMap) (*types.Extent, error) {

	geomCol, err := s.getGeometryColumn(tableName)
	if err != nil {
		return nil, err
	}
	geomColQuoted := pq.QuoteIdentifier(geomCol)

	// Query to get the bounding box coordinates
	keyName := mapType.GetKeyName()
	query := fmt.Sprintf("SELECT ST_XMin(ST_Extent(%s)), ST_YMin(ST_Extent(%s)), ST_XMax(ST_Extent(%s)), ST_YMax(ST_Extent(%s)) FROM %s WHERE %s = '%s'",
		geomColQuoted, geomColQuoted, geomColQuoted, geomColQuoted, tableName, keyName, idSls)

	var minX, minY, maxX, maxY sql.NullFloat64
	err = s.Db.QueryRow(query).Scan(&minX, &minY, &maxX, &maxY)

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch bounding box from database. Make sure your raster filename is correct and has a polygon in digital master map. Error :%v - %v - %s", tableName, idSls, err.Error())
	}
	if !minX.Valid || !minY.Valid || !maxX.Valid || !maxY.Valid {
		return nil, fmt.Errorf("no polygon found in table %q where %s = '%s' (raster key extracted from filename). Check that the value exists in the %s column and that %s is really the correct key column for this master map", tableName, keyName, idSls, keyName, keyName)
	}

	// Create a BoundingBox object with the coordinates
	extent := types.Extent{
		MinX: minX.Float64,
		MinY: minY.Float64,
		MaxX: maxX.Float64,
		MaxY: maxY.Float64,
	}

	return &extent, nil
}

func (s *PostgreStorage) GetAttributesValue(table string, attrKey string, key string, attributes []string) ([]string, error) {
	selectQuery := strings.Join(attributes, ",")
	query := fmt.Sprintf(`
	SELECT %s
		FROM %s
		WHERE %s = '%s'
	`, selectQuery, table, attrKey, key)

	// columns := make([]string, len(attributes))
	// columnPointers := make([]interface{}, len(attributes))
	// for i, _ := range columns {
	// 	columnPointers[i] = &columns[i]
	// }

	// err := s.Db.QueryRow(query).Scan(columnPointers...)
	columns := make([]string, len(attributes))
	err := s.Db.QueryRow(query).Scan(makeSqlScanFunc(columns)...)
	if err != nil {
		return nil, err
	}

	return columns, nil

}
func (s *PostgreStorage) CreateMasterMaps(tableName string, fileData *[]byte) error {
	tableExist, err := s.TableExist(tableName)
	if err != nil {
		return fmt.Errorf("Error when checking the table existence (%s) in database. %s", tableName, err.Error())
	}
	if tableExist {
		return fmt.Errorf("Layer name or table (%s) already exists in the database. ", tableName)
	}

	var data geojson.FeatureCollection
	if err := json.Unmarshal(*fileData, &data); err != nil {
		return fmt.Errorf("error Unmarshal json ")
	}
	i := 0
	for _, feature := range data.Features {
		// Convert the GeoJSON geometry to WKT
		//prop, err := json.Marshal(feature.Properties)

		if i == 0 {
			geometry, err := geojson.Encode(feature.Geometry)
			if err != nil {
				return fmt.Errorf("error geojson.Encode : %w", err)
			}

			// If the GeoJSON coordinates carry a Z value (elevation), the
			// geometry column must be declared with a Z dimension too --
			// otherwise PostGIS rejects every insert with
			// "Geometry has Z dimension but column does not".
			geomType := geometry.Type
			if feature.Geometry.Layout().ZIndex() >= 0 {
				geomType += "Z"
			}

			propTypes := constructDataTypes(feature.Properties)
			propTypes["geom"] = fmt.Sprintf("geometry(%s, %v)", geomType, 4326)

			_, err = s.createTable(tableName, propTypes)
			if err != nil {
				return fmt.Errorf("error createTable : %v : %w", tableName, err)
			}

			i++
			//return []byte(createQuery), nil
		}

		geom, err := geojson.Marshal(feature.Geometry)
		if err != nil {
			return fmt.Errorf("error geojson.Marshal : %w", err)
		}

		// Build the INSERT statement
		var data map[string]interface{}
		data = feature.Properties
		data["geom"] = []byte(fmt.Sprintf("ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326)", geom))
		_, err = s.insertData(tableName, data)
		if err != nil {
			return fmt.Errorf("error insertData : %w", err)
		}
	}
	return nil
}
func (s *PostgreStorage) DeleteMasterMap(masterMap string) error {
	exist, err := s.MasterMapExist(masterMap)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("Master maps doesnt exist")
	}

	query := fmt.Sprintf("DROP TABLE IF EXISTS %v", masterMap)

	_, err = s.Db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func constructDataTypes(prop map[string]interface{}) map[string]string {
	propTypes := map[string]string{}
	for key, val := range prop {
		if val == nil {
			propTypes[key] = fmt.Sprintf("%T", "")
			continue
		}
		if key == "geom" {
			key = "__geom"
		}
		propTypes[key] = fmt.Sprintf("%T", val)

		// UNTUK MEMBEDAKAN FLOAT64 DAN INT64

		// rawType := fmt.Sprintf("%T", val)
		// if rawType != "float64" {
		// 	propTypes[key] = fmt.Sprintf("%T", val)
		// 	continue
		// }

		// s := fmt.Sprintf("%v", val)
		// i, err := strconv.ParseInt(s, 10, 64)
		// if err == nil {
		// 	propTypes[key] = fmt.Sprintf("%T", i)
		// 	continue
		// }
		// f, err := strconv.ParseFloat(s, 64)
		// if err == nil {
		// 	propTypes[key] = fmt.Sprintf("%T", f)
		// 	continue
		// }

	}
	return propTypes
}
func (s *PostgreStorage) createTable(name string, attr map[string]string) (string, error) {
	// Construct Query statement. Column names come from user-uploaded GeoJSON
	// property names, so they must be quoted as identifiers -- otherwise a
	// property name with a space, hyphen, or other special character breaks
	// the SQL syntax.
	columns := []string{"gid serial primary key"}
	for key, val := range attr {
		var dType string
		switch val {
		case "float64":
			dType = "numeric"
		case "int64":
			dType = "integer"
		case "string":
			dType = "varchar(254)"
		default:
			dType = val
		}
		if key == "gid" {
			key = "__gid"
		}
		columns = append(columns, fmt.Sprintf("%s %s", pq.QuoteIdentifier(key), dType))
	}
	query := fmt.Sprintf("CREATE TABLE %s (\n%s\n);", pq.QuoteIdentifier(name), strings.Join(columns, ",\n"))

	// Execute the CREATE table
	_, err := s.Db.Exec(query)
	if err != nil {
		return "", err
	}
	return query, nil
}
func (s *PostgreStorage) insertData(tableName string, data map[string]interface{}) (string, error) {
	// Column names come from GeoJSON property names (quoted as identifiers,
	// same reasoning as createTable), and values are passed as query
	// parameters ($1, $2, ...) instead of being concatenated into the SQL
	// string -- this avoids breaking on special characters (e.g. an
	// apostrophe in a text value) and closes a SQL injection hole.
	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	argIndex := 1
	for key, val := range data {
		if key == "gid" {
			key = "__gid"
		}
		columns = append(columns, pq.QuoteIdentifier(key))

		if key == "geom" {
			// geom is already a raw SQL expression built by the caller
			// (e.g. ST_SetSRID(ST_GeomFromGeoJSON('...'), 4326)), not a
			// plain value, so it must stay embedded in the query text
			// rather than be passed as a parameter.
			placeholders = append(placeholders, fmt.Sprintf("%s", val))
			continue
		}

		placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
		values = append(values, val)
		argIndex++
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", pq.QuoteIdentifier(tableName), strings.Join(columns, ","), strings.Join(placeholders, ","))

	// Execute the INSERT statement
	_, err := s.Db.Exec(query, values...)
	if err != nil {
		return query, err
	}
	return query, nil
}
