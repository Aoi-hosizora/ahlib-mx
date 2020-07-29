package xgorm

import (
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/jinzhu/gorm"
)

type Helper struct {
	db *gorm.DB
}

// noinspection GoUnusedExportedFunction
func WithDB(db *gorm.DB) *Helper {
	return &Helper{db: db}
}

func (h *Helper) Pagination(limit int32, page int32) *gorm.DB {
	return h.db.Limit(limit).Offset((page - 1) * limit)
}

func (h *Helper) Count(model interface{}, where interface{}) uint64 {
	cnt := 0
	h.db.Model(model).Where(where).Count(&cnt)
	return uint64(cnt)
}

func (h *Helper) Exist(model interface{}, where interface{}) bool {
	return h.Count(model, where) > 0
}

func (h *Helper) Insert(model interface{}, object interface{}) xstatus.DbStatus {
	rdb := h.db.Model(model).Create(object)

	if IsMySqlDuplicateEntryError(rdb.Error) {
		return xstatus.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return xstatus.DbFailed
	}
	return xstatus.DbSuccess
}

func (h *Helper) Update(model interface{}, where interface{}, object interface{}) xstatus.DbStatus {
	if where == nil {
		where = object
	}
	rdb := h.db.Model(model).Where(where).Update(object)

	if IsMySqlDuplicateEntryError(rdb.Error) {
		return xstatus.DbExisted
	} else if rdb.Error != nil {
		return xstatus.DbFailed
	} else if rdb.RowsAffected == 0 {
		return xstatus.DbNotFound
	}
	return xstatus.DbSuccess
}

func (h *Helper) Delete(model interface{}, where interface{}, object interface{}) xstatus.DbStatus {
	if where == nil {
		where = object
	}
	rdb := h.db.Model(model).Where(where).Delete(object)

	if rdb.Error != nil {
		return xstatus.DbFailed
	} else if rdb.RowsAffected == 0 {
		return xstatus.DbNotFound
	}
	return xstatus.DbSuccess
}
