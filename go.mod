module github.com/Aoi-hosizora/ahlib-web

go 1.15

require (
	github.com/Aoi-hosizora/ahlib v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/tucnak/telebot.v2 v2.5.0
)

replace github.com/Aoi-hosizora/ahlib => ../ahlib
