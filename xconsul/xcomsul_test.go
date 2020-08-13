package xconsul

import (
	"log"
	"testing"
)

func TestConsul(t *testing.T) {
	RegisterConsulResolver()

	for i := 0; i < 5; i++ {
		err := RegisterConsulService("127.0.0.1:8500", &ConsulService{
			IP:   "127.0.0.1",
			Port: 5555 + i,
			Name: "test",
		})
		log.Println(err)
	}
}

func TestDefaultNameHandler(t *testing.T) {
	log.Println(defaultNameHandler.GetConsulID("127.0.0.1", 8500, "aaa"))
	log.Println(defaultNameHandler.GetGrpcTarget("127.0.0.1", 8500, "aaa"))
	log.Println(defaultNameHandler.ParseGrpcTarget("127.0.0.1:8500/aaa"))

	func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		SetDefaultNameHandler(nil)
	}()
}
