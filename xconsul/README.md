# xconsul

### Service

+ `type ConsulService struct {}`
+ `RegisterConsulService(ca string, srv *ConsulService) error`

### Client

+ `RegisterConsulResolver()`
+ `NewConsulBuilder() resolver.Builder`

### Name

+ `type NameHandler struct {}`
+ `SetDefaultNameHandler(hdr *NameHandler)`
