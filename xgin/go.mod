module github.com/Aoi-hosizora/ahlib-mx/xgin

go 1.16

require (
	github.com/Aoi-hosizora/ahlib v0.0.0-00010101000000-000000000000
	github.com/Aoi-hosizora/ahlib-mx/xvalidator v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.8.2
	github.com/go-playground/validator/v10 v10.11.1
	github.com/sirupsen/logrus v1.9.0
)

replace (
	github.com/Aoi-hosizora/ahlib => ../../ahlib
	github.com/Aoi-hosizora/ahlib-mx/xvalidator => ../xvalidator
)
