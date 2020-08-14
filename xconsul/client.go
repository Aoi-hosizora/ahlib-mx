package xconsul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
)

// RegisterConsulResolver will register consul build to grpc resolver, used in client.
func RegisterConsulResolver(doLog bool) {
	resolver.Register(NewConsulBuilder(doLog))
}

type ConsulBuilder struct {
	doLog bool
}

func NewConsulBuilder(doLog bool) resolver.Builder {
	return &ConsulBuilder{doLog: doLog}
}

func (cb *ConsulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 127.0.0.1:8500/xxx
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
	go cr.watcher(cb.doLog)
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
	lastAddrCount        int
}

func (cr *consulResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (cr *consulResolver) Close() {}

func (cr *consulResolver) watcher(doLog bool) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", cr.host, cr.port)
	client, err := api.NewClient(config)
	if err != nil {
		if doLog {
			log.Println("Failed to create consul client:", err)
		}
		return
	}

	for {
		services, metaInfo, err := client.Health().Service(cr.name, cr.name, true, &api.QueryOptions{WaitIndex: cr.lastIndex})
		if err != nil && doLog {
			log.Println("Failed to retrieve instances from consul:", err)
		}

		cr.lastIndex = metaInfo.LastIndex
		addresses := make([]resolver.Address, len(services))
		for idx, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			addresses[idx] = resolver.Address{Addr: addr}
		}

		if l := len(addresses); cr.lastAddrCount != l {
			if doLog {
				log.Printf("Addresses updated: #%d", len(addresses))
			}
			cr.lastAddrCount = l
		}

		cr.cc.UpdateState(resolver.State{Addresses: addresses})
	}
}
