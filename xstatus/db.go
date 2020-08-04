package xstatus

type DbStatus uint8

const (
	DbSuccess DbStatus = iota
	DbNotFound
	DbExisted
	DbFailed
	DbTagA
	DbTagB
	DbTagC
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
	case DbTagA:
		return "db-tag-a"
	case DbTagB:
		return "db-tag-b"
	case DbTagC:
		return "db-tag-c"
	default:
		return "db-?"
	}
}
