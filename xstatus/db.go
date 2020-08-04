package xstatus

type DbStatus uint8

const (
	DbSuccess DbStatus = iota
	DbNotFound
	DbExisted
	DbFailed
)

func (d DbStatus) String() string {
	switch d {
	case DbSuccess:
		return "db-success"
	case DbNotFound:
		return "db-not-found"
	case DbExisted:
		return "db-existed"
	case DbFailed:
		return "db-failed"
	default:
		return "db-?"
	}
}
