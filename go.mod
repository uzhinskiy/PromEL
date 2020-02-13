module github.com/uzhinskiy/PromEL

replace (
	github.com/uzhinskiy/PromEL/modules/config => ./modules/config
	github.com/uzhinskiy/PromEL/modules/driver => ./modules/driver
	github.com/uzhinskiy/PromEL/modules/es => ./modules/es
)

require (
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/golang/snappy v0.0.1
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	github.com/olivere/elastic/v7 v7.0.10
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.9.1
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/uzhinskiy/PromEL/modules/config v0.0.0
	github.com/uzhinskiy/PromEL/modules/driver v0.0.0
	github.com/uzhinskiy/PromEL/modules/es v0.0.0
	google.golang.org/genproto v0.0.0-20200205142000-a86caf926a67 // indirect
	gopkg.in/yaml.v2 v2.2.8
)
