package xneo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
	"log"
)

// logrus.Logger

type Neo4jLogrus struct {
	neo4j.Session
	logger  *logrus.Logger
	LogMode bool
}

func NewNeo4jLogrus(session neo4j.Session, logger *logrus.Logger, logMode bool) *Neo4jLogrus {
	return &Neo4jLogrus{Session: session, logger: logger, LogMode: logMode}
}

func (n *Neo4jLogrus) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	result, err := n.Session.Run(cypher, params, configurers...)
	if n.LogMode {
		n.print(result, err)
	}
	return result, err
}

func (n *Neo4jLogrus) print(result neo4j.Result, err error) {
	if err != nil {
		// failed to run cypher
		n.logger.WithFields(logrus.Fields{
			"module": "neo4j", "error": err,
		}).Infoln(fmt.Printf("[Neo4j] error: %v\n", err))
		return
	}
	summary, err := result.Summary()
	if err != nil {
		// success to run cypher but get failed to get summary
		// Neo.ClientError.Statement.SyntaxError
		// Neo.ClientError.Schema.ConstraintValidationFailed
		// ...
		n.logger.WithFields(logrus.Fields{
			"module": "neo4j", "error": err,
		}).Infoln(fmt.Printf("[Neo4j] error: %v\n", err))
		return
	}

	// success to run cypher and get summary, get information form summary
	stat := summary.Statement()
	cypher := stat.Text()
	params := stat.Params()
	field := n.logger.WithFields(logrus.Fields{
		"module": "neo4j",
		"cypher": cypher,
		"params": params,
	})
	field.Infof("[Neo4j] %s <- %v\n", cypher, params)
}

// log.Logger

type Neo4jLogger struct {
	neo4j.Session
	logger  *log.Logger
	LogMode bool
}

func NewNeo4jLogger(session neo4j.Session, logger *log.Logger, logMode bool) *Neo4jLogger {
	return &Neo4jLogger{Session: session, logger: logger, LogMode: logMode}
}

func (n *Neo4jLogger) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	result, err := n.Session.Run(cypher, params, configurers...)
	if n.LogMode {
		n.print(result, err)
	}
	return result, err
}

func (n *Neo4jLogger) print(result neo4j.Result, err error) {
	if err != nil {
		n.logger.Printf("[Neo4j] error: %v\n", err)
		return
	}
	summary, err := result.Summary()
	if err != nil {
		n.logger.Printf("[Neo4j] error: %v\n", err)
		return
	}

	stat := summary.Statement()
	cypher := stat.Text()
	params := stat.Params()
	n.logger.Printf("[Neo4j] %s <- %v\n", cypher, params)
}
