package xneo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"strings"
	"time"
)

// logrus

type LogrusNeo4j struct {
	neo4j.Session
	logger  *logrus.Logger
	LogMode bool
}

func NewLogrusNeo4j(session neo4j.Session, logger *logrus.Logger, logMode bool) *LogrusNeo4j {
	return &LogrusNeo4j{Session: session, logger: logger, LogMode: logMode}
}

func (n *LogrusNeo4j) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	s := time.Now()
	result, err := n.Session.Run(cypher, params, configurers...)
	e := time.Now()
	if n.LogMode {
		n.print(result, e.Sub(s).String(), err)
	}
	return result, err
}

func (n *LogrusNeo4j) print(result neo4j.Result, du string, err error) {
	if err != nil { // Failed to run cypher.
		n.logger.WithFields(logrus.Fields{
			"module": "neo4j",
			"error":  err,
		}).Infoln(fmt.Printf("[Neo4j] error: %v", err))
		return
	}

	summary, err := result.Summary()
	if err != nil { // Failed to get summary.
		// Neo.ClientError.Statement.SyntaxError
		// Neo.ClientError.Schema.ConstraintValidationFailed
		// ...
		n.logger.WithFields(logrus.Fields{
			"module": "neo4j",
			"error":  err,
		}).Infoln(fmt.Printf("[Neo4j] error: %v", err))
		return
	}

	keys, err := result.Keys()
	if err != nil { // Failed to get keys.
		n.logger.WithFields(logrus.Fields{
			"module": "neo4j",
			"error":  err,
		}).Infoln(fmt.Printf("[Neo4j] error: %v", err))
		return
	}

	// Success to run cypher and get summary, get information form summary.
	stat := summary.Statement()
	cypher := stat.Text()
	params := stat.Params()
	cypher = render(cypher, params)

	n.logger.WithFields(logrus.Fields{
		"module":   "neo4j",
		"cypher":   cypher,
		"rows":     0,
		"columns":  len(keys),
		"duration": du,
	}).Info(fmt.Sprintf("[Neo4j] #: ?x%d | %10s | %s", len(keys), du, cypher)) // TODO
}

// logger

type LoggerNeo4j struct {
	neo4j.Session
	logger  *log.Logger
	LogMode bool
}

func NewLoggerNeo4j(session neo4j.Session, logger *log.Logger, logMode bool) *LoggerNeo4j {
	return &LoggerNeo4j{Session: session, logger: logger, LogMode: logMode}
}

func (n *LoggerNeo4j) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	s := time.Now()
	result, err := n.Session.Run(cypher, params, configurers...)
	e := time.Now()
	if n.LogMode {
		n.print(result, e.Sub(s).String(), err)
	}
	return result, err
}

func (n *LoggerNeo4j) print(result neo4j.Result, du string, err error) {
	if err != nil {
		n.logger.Printf("[Neo4j] error: %v\n", err)
		return
	}
	summary, err := result.Summary()
	if err != nil {
		n.logger.Printf("[Neo4j] error: %v\n", err)
		return
	}
	keys, err := result.Keys()
	if err != nil {
		n.logger.Printf("[Neo4j] error: %v\n", err)
		return
	}

	stat := summary.Statement()
	cypher := stat.Text()
	params := stat.Params()
	cypher = render(cypher, params)

	n.logger.Printf("[Neo4j] #: ?x%d | %10s | %s", len(keys), du, cypher) // TODO
}

// render

func render(cypher string, params map[string]interface{}) string {
	out := cypher
	for k, v := range params {
		t := reflect.TypeOf(v)
		to := fmt.Sprintf("%v", v)
		if t.Kind() == reflect.String {
			to = "'" + to + "'"
		}
		out = strings.ReplaceAll(out, "$"+k, to)
	}
	return out
}
