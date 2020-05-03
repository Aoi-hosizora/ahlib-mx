package xneo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

// type neo4j.Session interface{xxx}
type Neo4jSession struct {
	session neo4j.Session
	logger  *logrus.Logger
	LogMode bool
}

func NewNeo4jSessionWithLogger(session neo4j.Session, logger *logrus.Logger) *Neo4jSession {
	return &Neo4jSession{session: session, logger: logger}
}

func (n *Neo4jSession) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	result, err := n.session.Run(cypher, params, configurers...)
	if n.LogMode {
		n.print(result, err)
	}
	return result, err
}

func (n *Neo4jSession) print(result neo4j.Result, err error) {
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

func (n *Neo4jSession) Close() error {
	return n.session.Close()
}

func (n *Neo4jSession) LastBookmark() string {
	return n.session.LastBookmark()
}

func (n *Neo4jSession) BeginTransaction(configurers ...func(*neo4j.TransactionConfig)) (neo4j.Transaction, error) {
	return n.session.BeginTransaction(configurers...)
}

func (n *Neo4jSession) ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	return n.session.ReadTransaction(work, configurers...)
}

func (n *Neo4jSession) WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	return n.session.WriteTransaction(work, configurers...)
}
