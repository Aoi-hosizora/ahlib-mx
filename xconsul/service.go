package xconsul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type ConsulService struct {
	IP   string
	Port int
	Tag  []string
	Name string

	Interval   time.Duration
	Deregister time.Duration

	Namespace string
	Weights   *api.AgentWeights
	Proxy     *api.AgentServiceConnectProxyConfig
	Connect   *api.AgentServiceConnect
}

// RegisterConsulService will connect consul agent, and register a service into it.
func RegisterConsulService(consulAddress string, srv *ConsulService) error {
	// connect consul
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	agent := client.Agent()

	// check some parameter
	if srv.IP == "" || srv.Name == "" {
		return fmt.Errorf("expected non-empty IP and Name")
	}
	if srv.Port <= 0 || srv.Port >= 65536 {
		return fmt.Errorf("invalid port value")
	}
	if srv.Interval == 0 {
		srv.Interval = time.Duration(10) * time.Second
	}
	if srv.Deregister == 0 {
		srv.Deregister = time.Duration(1) * time.Minute
	}

	// set registration
	registration := &api.AgentServiceRegistration{
		ID:      defaultNameHandler.GetConsulID(srv.IP, srv.Port, srv.Name),
		Name:    srv.Name,
		Tags:    srv.Tag,
		Port:    srv.Port,
		Address: srv.IP,
		Check: &api.AgentServiceCheck{
			Interval:                       srv.Interval.String(),
			DeregisterCriticalServiceAfter: srv.Deregister.String(),
			GRPC:                           defaultNameHandler.GetGrpcTarget(srv.IP, srv.Port, srv.Name),
		},

		Namespace: srv.Namespace,
		Weights:   srv.Weights,
		Proxy:     srv.Proxy,
		Connect:   srv.Connect,
	}

	// register to agent
	err = agent.ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}
