package xconsul

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"regexp"
)

// nameHandler Include some handler used for grpc and consul.
type nameHandler struct {
	// Used in api.AgentServiceRegistration `ID` field.
	GetConsulID func(ip string, port int, name string) string

	// Used in api.AgentServiceCheck `GRPC` field.
	GetGrpcTarget func(ip string, port int, name string) string

	// Used to parse api.AgentServiceCheck `GRPC` field in ConsulBuilder.Build.
	ParseGrpcTarget func(target string) (host string, port int, name string, err error)
}

// Default nameHandler used for inner.
var defaultNameHandler = &nameHandler{
	GetConsulID:     _defaultGetConsulID,
	GetGrpcTarget:   _defaultGetGrpcTarget,
	ParseGrpcTarget: _defaultParseGrpcTarget,
}

// Change default implement of nameHandler.GetConsulID.
func SetDefaultGetConsulIDHandler(f func(ip string, port int, name string) string) {
	defaultNameHandler.GetConsulID = f
}

// Change default implement of nameHandler.GetGrpcTarget.
func SetDefaultGetGrpcTargetHandler(f func(ip string, port int, name string) string) {
	defaultNameHandler.GetGrpcTarget = f
}

// Change default implement of nameHandler.ParseGrpcTarget.
func SetDefaultParseGrpcTargetHandler(f func(target string) (host string, port int, name string, err error)) {
	defaultNameHandler.ParseGrpcTarget = f
}

// Default nameHandler.GetConsulID.
func _defaultGetConsulID(ip string, port int, name string) string {
	return fmt.Sprintf("%s-%d-%s", ip, port, name)
}

// Default nameHandler.GetGrpcTarget.
func _defaultGetGrpcTarget(ip string, port int, name string) string {
	return fmt.Sprintf("%s:%d/%s", ip, port, name)
}

// Default nameHandler.ParseGrpcTarget.
func _defaultParseGrpcTarget(target string) (host string, port int, name string, err error) {
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
