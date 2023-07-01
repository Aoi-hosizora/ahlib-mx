# ahlib-mx

[![Build Status](https://app.travis-ci.com/Aoi-hosizora/ahlib-mx.svg?branch=master)](https://app.travis-ci.com/github/Aoi-hosizora/ahlib-mx)
[![Codecov](https://codecov.io/gh/Aoi-hosizora/ahlib-mx/branch/master/graph/badge.svg)](https://codecov.io/gh/Aoi-hosizora/ahlib-mx)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aoi-hosizora/ahlib-mx)](https://goreportcard.com/report/github.com/Aoi-hosizora/ahlib-mx)
[![License](http://img.shields.io/badge/license-mit-blue.svg)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/Aoi-hosizora/ahlib-mx)](https://github.com/Aoi-hosizora/ahlib-mx/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/Aoi-hosizora/ahlib-mx.svg)](https://pkg.go.dev/github.com/Aoi-hosizora/ahlib-mx)

+ A personal golang library for web and db development, which depends on more heavy and complex third-party libraries, requires `Go >= 1.15` or `Go >= 1.16`.
+ This package includes following utilities:
    + [gin](https://github.com/gin-gonic/gin)
    + [validator (v10)](https://github.com/go-playground/validator)
    + [telebot (v2)](https://github.com/tucnak/telebot)
    + [gorm (v1)](https://github.com/jinzhu/gorm) / [gorm (v2)](https://github.com/go-gorm/gorm) (mysql+sqlite+postgres)
    + [go-redis (v8)](https://github.com/go-redis/redis)
    + [neo4j-go-driver (v1)](https://github.com/neo4j/neo4j-go-driver)

### Related libraries

+ [Aoi-hosizora/ahlib](https://github.com/Aoi-hosizora/ahlib)
+ [Aoi-hosizora/ahlib-more](https://github.com/Aoi-hosizora/ahlib-more)
+ [Aoi-hosizora/ahlib-mx](https://github.com/Aoi-hosizora/ahlib-mx)

### Packages

+ xdbutils/*
+ xgin
+ xgorm
+ xgormv2
+ xneo4j
+ xredis
+ xtelebot
+ xvalidator

### Dependencies (for web development)

#### xgin

+ See [go.mod](./xgin/go.mod) and [go.sum](./xgin/go.sum)
+ `github.com/Aoi-hosizora/ahlib v1.6.0`
+ `github.com/Aoi-hosizora/ahlib-mx/xvalidator v1.6.0`
+ `github.com/gin-gonic/gin v1.8.2`
+ `github.com/go-playground/validator/v10 v10.11.1`
+ `github.com/sirupsen/logrus v1.9.0`
+ `golang.org/x/sys v0.3.0`

#### xtelebot

+ See [go.mod](./xtelebot/go.mod) and [go.sum](./xtelebot/go.sum)
+ `gopkg.in/tucnak/telebot.v2 v2.5.0`
+ `github.com/gin-gonic/gin v1.8.2`
+ `github.com/sirupsen/logrus v1.9.0`

#### xvalidator

+ See [go.mod](./xvalidator/go.mod) and [go.sum](./xvalidator/go.sum)
+ `github.com/Aoi-hosizora/ahlib v1.6.0`
+ `github.com/go-playground/validator/v10 v10.11.1`
+ `github.com/go-playground/locales v0.14.0`
+ `github.com/go-playground/universal-translator v0.18.0`

### Dependencies (for database development)

#### xgorm

+ See [go.mod](./xgorm/go.mod) and [go.sum](./xgorm/go.sum)
+ `github.com/Aoi-hosizora/ahlib v1.6.0`
+ `github.com/Aoi-hosizora/ahlib-mx/xdbutils v1.6.0`
+ `github.com/jinzhu/gorm v1.9.16`
+ `github.com/go-sql-driver/mysql v1.5.0`
+ `github.com/VividCortex/mysqlerr v1.0.0`
+ `github.com/mattn/go-sqlite3 v1.14.0`
+ `github.com/lib/pq v1.1.1`
+ `github.com/sirupsen/logrus v1.9.0`

#### xgormv2

+ See [go.mod](./xgormv2/go.mod) and [go.sum](./xgormv2/go.sum)
+ `github.com/Aoi-hosizora/ahlib v1.6.0`
+ `github.com/Aoi-hosizora/ahlib-mx/xdbutils v1.6.0`
+ `gorm.io/gorm v1.22.4`
+ `gorm.io/driver/mysql v1.2.3`
+ `gorm.io/driver/sqlite v1.2.6`
+ `github.com/go-sql-driver/mysql v1.6.0`
+ `github.com/VividCortex/mysqlerr v1.0.0`
+ `github.com/mattn/go-sqlite3 v1.14.9`
+ `github.com/sirupsen/logrus v1.9.0`

#### xneo4j

+ See [go.mod](./xneo4j/go.mod) and [go.sum](./xneo4j/go.sum)
+ `github.com/Aoi-hosizora/ahlib v1.6.0`
+ `github.com/Aoi-hosizora/ahlib-mx/xdbutils v1.6.0`
+ `github.com/neo4j/neo4j-go-driver v1.8.3`
+ `github.com/sirupsen/logrus v1.9.0`

#### xredis

+ See [go.mod](./xredis/go.mod) and [go.sum](./xredis/go.sum)
+ `github.com/Aoi-hosizora/ahlib v1.6.0`
+ `github.com/go-redis/redis/v8 v8.4.11`
+ `github.com/sirupsen/logrus v1.9.0`
