package xneo4j

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Get records from run result, see neo4j.Collect.
//
// 1. result: array of neo4j.Record;
//
// 2. record: array of columns value -> Get(key) / GetByIndex(index);
//
// Example:
//	cypher := "MATCH p = ()-[r :FRIEND]->(n) RETURN r, n"
//	rec, _ := xneo4j.GetRecords(session.Run(cypher, nil)) // slice of neo4j.Record
//	for _, r := range rec { // slice of value (interface{})
//		rel := xneo4j.GetRel(r.Values()[0]) // neo4j.Node
//		node := xneo4j.GetNode(r.Values()[1]) // neo4j.Relationship
//		log.Println(rel.Id(), rel.Type(), node.Id(), node.Props())
//	}
func GetRecords(result neo4j.Result, err error) ([]neo4j.Record, error) {
	if err != nil {
		return nil, err
	}

	rec := make([]neo4j.Record, 0)
	for result.Next() {
		rec = append(rec, result.Record())
	}
	if err := result.Err(); err != nil {
		return nil, err
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
