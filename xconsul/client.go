package xconsul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
)

// RegisterConsulResolver will register consul build to grpc resolver, used in client.
func RegisterConsulResolver() {
	resolver.Register(NewConsulBuilder())
}

type ConsulBuilder struct{}

func NewConsulBuilder() resolver.Builder {
	return &ConsulBuilder{}
}

func (cb *ConsulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ta := fmt.Sprintf("%s/%s", target.Authority, target.Endpoint)
	host, port, name, err := defaultNameHandler.ParseGrpcTarget(ta)
	if err != nil {
		return nil, err
	}

	cr := &consulResolver{
		host:                 host,
		port:                 port,
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	cr.wg.Add(1)
	go cr.watcher()
	return cr, nil
}

func (cb *ConsulBuilder) Scheme() string {
	return "consul"
}

type consulResolver struct {
	host                 string
	port                 int
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	lastIndex            uint64
}

func (cr *consulResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (cr *consulResolver) Close() {}

func (cr *consulResolver) watcher() {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", cr.host, cr.port)
	client, err := api.NewClient(config)
	if err != nil {
		log.Println("Failed to create consul client:", err)
		return
	}

	for {
		services, metaInfo, err := client.Health().Service(cr.name, cr.name, true, &api.QueryOptions{WaitIndex: cr.lastIndex})
		if err != nil {
			log.Println("Failed to retrieve instances from consul:", err)
		}

		cr.lastIndex = metaInfo.LastIndex
		addresses := make([]resolver.Address, len(services))
		for idx, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			addresses[idx] = resolver.Address{Addr: addr}
		}

		log.Println("Update service addresses:", len(addresses))
		cr.cc.UpdateState(resolver.State{Addresses: addresses})
	}
}
