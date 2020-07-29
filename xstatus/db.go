package xstatus

type DbStatus uint8

// noinspection GoUnusedConst
const (
	DbSuccess DbStatus = iota
	DbNotFound
	DbExisted
	DbFailed
)
