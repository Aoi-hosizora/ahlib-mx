package xconsul

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"log"
	"regexp"
	"testing"
)

func TestConsul(t *testing.T) {
	RegisterConsulResolver(true)

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

	SetDefaultGetConsulIDHandler(func(ip string, port int, name string) string {
		return fmt.Sprintf("%s:%d:%s", ip, port, name)
	})
	SetDefaultGetGrpcTargetHandler(func(ip string, port int, name string) string {
		return fmt.Sprintf("%s:%d//%s", ip, port, name)
	})
	SetDefaultParseGrpcTargetHandler(func(target string) (host string, port int, name string, err error) {
		regexConsul, _ := regexp.Compile(`^([A-z0-9.]+)(?::([0-9]{1,5}))?//([A-z_-]+)$`)
		if !regexConsul.MatchString(target) {
			return "", 0, "", fmt.Errorf("consul resolver: invalid uri")
		}
		groups := regexConsul.FindStringSubmatch(target)
		host = groups[1]                  // localhost
		name = groups[3]                  // xxx
		port, _ = xnumber.Atoi(groups[2]) // 8500
		return host, port, name, nil
	})

	for i := 0; i < 5; i++ {
		err := RegisterConsulService("127.0.0.1:8500", &ConsulService{
			IP:   "127.0.0.1",
			Port: 5555 + i,
			Name: "test",
		})
		log.Println(err)
	}
}
