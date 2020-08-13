package xconsul

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"regexp"
)

type NameHandler struct {
	// Used in api.AgentServiceRegistration `ID` field.
	GetConsulID func(ip string, port int, name string) string

	// Used in api.AgentServiceCheck `GRPC` field.
	GetGrpcTarget func(ip string, port int, name string) string

	// Used to parse api.AgentServiceCheck `GRPC` field in ConsulBuilder.Build.
	ParseGrpcTarget func(schema string) (ip string, port int, name string, err error)
}

var defaultNameHandler = &NameHandler{
	GetConsulID:     getConsulIDDefault,
	GetGrpcTarget:   getGrpcTargetDefault,
	ParseGrpcTarget: parseGrpcTargetDefault,
}

// Change default NameHandler.
func SetDefaultNameHandler(hdr *NameHandler) {
	if hdr == nil || hdr.GetConsulID == nil || hdr.GetGrpcTarget == nil || hdr.ParseGrpcTarget == nil {
		panic("invalid name handler, could not be nil")
	}

	defaultNameHandler = hdr
}

func getConsulIDDefault(ip string, port int, name string) string {
	return fmt.Sprintf("%s-%d-%s", ip, port, name)
}

func getGrpcTargetDefault(ip string, port int, name string) string {
	return fmt.Sprintf("%s:%d/%s", ip, port, name)
}

func parseGrpcTargetDefault(target string) (host string, port int, name string, err error) {
	if target == "" {
		return "", 0, "", fmt.Errorf("consul resolver: missing address")
	}

	// localhost:8500/xxx
	regexConsul, err := regexp.Compile(`^([A-z0-9.]+)(?::([0-9]{1,5}))?/([A-z_-]+)$`)
	if err != nil {
		return "", 0, "", err
	}
	if !regexConsul.MatchString(target) {
		return "", 0, "", fmt.Errorf("consul resolver: invalid uri")
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]                    // localhost
	name = groups[3]                    // xxx
	port, err = xnumber.Atoi(groups[2]) // 8500
	if err != nil {
		return "", 0, "", err
	}
	if port == 0 {
		port = 8500
	}

	return host, port, name, nil
}
