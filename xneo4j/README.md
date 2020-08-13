# xneo4j

### Logger Functions

+ `type LogrusNeo4j struct {}`
+ `NewLogrusNeo4j(session neo4j.Session, logger *logrus.Logger, logMode bool) *LogrusNeo4j`
+ `type LoggerNeo4j struct {}`
+ `NewLoggerNeo4j(session neo4j.Session, logger *log.Logger, logMode bool) *LoggerNeo4j`

### Functions

+ `GetRecords(result neo4j.Result) ([]neo4j.Record, error)`
+ `GetInteger(data interface{}) int64`
+ `GetFloat(data interface{}) float64`
+ `GetString(data interface{}) string`
+ `GetBoolean(data interface{}) bool`
+ `GetByteArray(data interface{}) []byte`
+ `GetList(data interface{}) []interface{}`
+ `GetMap(data interface{}) map[string]interface{}`
+ `GetNode(data interface{}) neo4j.Node`
+ `GetRel(data interface{}) neo4j.Relationship`
+ `GetPath(data interface{}) neo4j.Path`
