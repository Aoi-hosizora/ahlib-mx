package ahlib_gin_gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

func HookGormDeleteAtField(db *gorm.DB, defaultDeleteAtTimeStamp string) {
	db.Callback().Query().Before("gorm:query").
		Register("new_deleted_at_before_query_callback", newBeforeQueryUpdateCallback(defaultDeleteAtTimeStamp))
	db.Callback().RowQuery().Before("gorm:row_query").
		Register("new_deleted_at_before_row_query_callback", newBeforeQueryUpdateCallback(defaultDeleteAtTimeStamp))
	db.Callback().Update().Before("gorm:update").
		Register("new_deleted_at_before_update_callback", newBeforeQueryUpdateCallback(defaultDeleteAtTimeStamp))
	db.Callback().Delete().Replace("gorm:delete", newDeleteCallback(defaultDeleteAtTimeStamp))
}

func newBeforeQueryUpdateCallback(defaultDeleteAtTimeStamp string) func(scope *gorm.Scope) {
	// https://qiita.com/touyu/items/f1ac43b186cd6b26b8c7
	return func(scope *gorm.Scope) {
		var (
			quotedTableName                   = scope.QuotedTableName()
			deletedAtField, hasDeletedAtField = scope.FieldByName("DeletedAt")
			defaultTimeStamp                  = defaultDeleteAtTimeStamp
		)

		if !scope.HasError() && !scope.Search.Unscoped && hasDeletedAtField {
			scope.Search.Unscoped = true
			sql := fmt.Sprintf("%v.%v = '%v'", quotedTableName, scope.Quote(deletedAtField.DBName), defaultTimeStamp)
			scope.Search.Where(sql)
		}
	}
}

func newDeleteCallback(defaultDeleteAtTimeStamp string) func(scope *gorm.Scope) {
	// https://github.com/jinzhu/gorm/blob/master/callback_delete.go
	return func(scope *gorm.Scope) {
		if scope.HasError() {
			return
		}
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		var (
			quotedTableName                   = scope.QuotedTableName()
			deletedAtField, hasDeletedAtField = scope.FieldByName("DeletedAt")
			defaultTimeStamp                  = defaultDeleteAtTimeStamp
		)

		addExtraSpaceIfExist := func(str string) string {
			if str != "" {
				return " " + str
			}
			return ""
		}

		if !scope.Search.Unscoped && hasDeletedAtField {
			var (
				comb = scope.CombinedConditionSql()
				from = fmt.Sprintf("%s IS NULL", scope.Quote(deletedAtField.DBName))
				to   = fmt.Sprintf("%s = '%s'", scope.Quote(deletedAtField.DBName), defaultTimeStamp)
				now  = time.Now().Format("2006-01-02 15:04:05")
			)
			comb = strings.Replace(comb, from, to, 1)

			sql := fmt.Sprintf(
				"UPDATE %v SET %v='%v'%v%v",
				quotedTableName,
				scope.Quote(deletedAtField.DBName), now,
				addExtraSpaceIfExist(comb),
				addExtraSpaceIfExist(extraOption),
			)
			scope.Raw(sql).Exec()
		} else {
			sql := fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)
			scope.Raw(sql).Exec()
		}
	}
}
