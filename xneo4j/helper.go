package xneo4j

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// get records from run result
// result: array of record
// record: array of columns value -> Get(key) / GetByIndex(index)
// noinspection GoUnusedExportedFunction
func GetRecords(result neo4j.Result) ([]neo4j.Record, error) {
	rec := make([]neo4j.Record, 0)
	for result.Next() {
		if result.Err() != nil {
			// ??? Get result record error
			continue
		}
		rec = append(rec, result.Record())
	}

	return rec, nil
}

// noinspection GoUnusedExportedFunction
func GetInteger(data interface{}) int64 {
	return data.(int64)
}

// noinspection GoUnusedExportedFunction
func GetFloat(data interface{}) float64 {
	return data.(float64)
}

// noinspection GoUnusedExportedFunction
func GetString(data interface{}) string {
	return data.(string)
}

// noinspection GoUnusedExportedFunction
func GetBoolean(data interface{}) bool {
	return data.(bool)
}

// noinspection GoUnusedExportedFunction
func GetByteArray(data interface{}) []byte {
	return data.([]byte)
}

// noinspection GoUnusedExportedFunction
func GetList(data interface{}) []interface{} {
	return data.([]interface{})
}

// noinspection GoUnusedExportedFunction
func GetMap(data interface{}) map[string]interface{} {
	return data.(map[string]interface{})
}

// noinspection GoUnusedExportedFunction
func GetNode(data interface{}) neo4j.Node {
	return data.(neo4j.Node)
}

// noinspection GoUnusedExportedFunction
func GetRel(data interface{}) neo4j.Relationship {
	return data.(neo4j.Relationship)
}

// noinspection GoUnusedExportedFunction
func GetPath(data interface{}) neo4j.Path {
	return data.(neo4j.Path)
}
