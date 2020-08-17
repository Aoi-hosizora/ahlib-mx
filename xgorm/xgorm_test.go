package xgorm

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
)

type TblTest struct {
	Id   uint64
	Name string
}

func TestLogrus(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123@tcp(localhost:3306)/db_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}

	db.SingularTable(true)
	db.LogMode(true)

	db.SetLogger(NewGormLogrus(logrus.New()))
	HookDeleteAtField(db, DefaultDeleteAtTimeStamp)

	test := &TblTest{}
	db.Model(&TblTest{}).First(test)
	log.Println(test)
	tests := make([]*TblTest, 0)
	db.Model(&TblTest{}).Find(&tests)
	log.Println(tests)
	db.Model(test).Related(test)
	log.Println(tests)
}

func TestLogger(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123@tcp(localhost:3306)/db_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}

	db.SingularTable(true)
	db.LogMode(true)

	// logger
	db.SetLogger(NewGormLogger(log.New(os.Stderr, "", log.LstdFlags)))
	HookDeleteAtField(db, DefaultDeleteAtTimeStamp)

	test := &TblTest{}
	db.Model(&TblTest{}).First(test)
	log.Println(test)
	tests := make([]*TblTest, 0)
	db.Model(&TblTest{}).Find(&tests)
	log.Println(tests)
}
