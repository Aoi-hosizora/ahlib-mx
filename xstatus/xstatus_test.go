package xstatus

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"testing"
)

func TestDbStatus(t *testing.T) {
	xtesting.Equal(t, fmt.Sprintf("%v", DbSuccess), "db-success")
	xtesting.Equal(t, DbSuccess.String(), "db-success")
	xtesting.Equal(t, DbNotFound.String(), "db-not-found")
	xtesting.Equal(t, DbExisted.String(), "db-existed")
	xtesting.Equal(t, DbFailed.String(), "db-failed")
	xtesting.Equal(t, DbTagA.String(), "db-tag-a")
	xtesting.Equal(t, DbTagB.String(), "db-tag-b")
	xtesting.Equal(t, DbTagC.String(), "db-tag-c")
	xtesting.Equal(t, DbStatus(20).String(), "db-?")
}

func TestFsmStatus(t *testing.T) {
	xtesting.Equal(t, fmt.Sprintf("%v", FsmNone), "fsm-none")
	xtesting.Equal(t, FsmNone.String(), "fsm-none")
	xtesting.Equal(t, FsmInState.String(), "fsm-in-state")
	xtesting.Equal(t, FsmFinal.String(), "fsm-final")
	xtesting.Equal(t, FsmStatus(20).String(), "fsm-?")
}
