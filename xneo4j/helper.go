package xneo4j

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// get records from run result
// result: array of record
// record: array of columns value -> Get(key) / GetByIndex(index)
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

func GetInteger(data interface{}) int64 {
	return data.(int64)
}

func GetFloat(data interface{}) float64 {
	return data.(float64)
}

func GetString(data interface{}) string {
	return data.(string)
}

func GetBoolean(data interface{}) bool {
	return data.(bool)
}

func GetByteArray(data interface{}) []byte {
	return data.([]byte)
}

func GetList(data interface{}) []interface{} {
	return data.([]interface{})
}

func GetMap(data interface{}) map[string]interface{} {
	return data.(map[string]interface{})
}

func GetNode(data interface{}) neo4j.Node {
	return data.(neo4j.Node)
}

func GetRel(data interface{}) neo4j.Relationship {
	return data.(neo4j.Relationship)
}

func GetPath(data interface{}) neo4j.Path {
	return data.(neo4j.Path)
}
