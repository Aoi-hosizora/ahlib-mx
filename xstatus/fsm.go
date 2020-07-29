package xstatus

type FsmStatus uint8

// noinspection GoUnusedConst
const (
	FsmNone FsmStatus = iota
	FsmInState
	FsmFinal
)
