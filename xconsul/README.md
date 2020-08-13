# xconsul

### Service

+ `type ConsulService struct {}`
+ `RegisterConsulService(ca string, srv *ConsulService) error`

### Client

+ `type ConsulBuilder struct {}`
+ `RegisterConsulResolver()`
+ `NewConsulBuilder() resolver.Builder`

### Name

+ `SetDefaultGetConsulIDHandler(f func(ip string, port int, name string) string)`
+ `SetDefaultGetGrpcTargetHandler(f func(ip string, port int, name string) string)`
+ `SetDefaultParseGrpcTargetHandler(f func(schema string) (ip string, port int, name string, err error))`
