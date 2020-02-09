module github.com/uzhinskiy/PromEL

//replace (
//	gitlab.insitu.co.il/boris.uzhinskiy/PromEL/config => ./config
//	gitlab.insitu.co.il/boris.uzhinskiy/PromEL/driver => ./driver
//	gitlab.insitu.co.il/boris.uzhinskiy/PromEL/es => ./es
//)

require (
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/golang/snappy v0.0.1
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	github.com/olivere/elastic/v7 v7.0.10
	github.com/prometheus/common v0.9.1
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/uzhinskiy/PromEL/config v0.0.1
	github.com/uzhinskiy/PromEL/driver v0.0.1
	github.com/uzhinskiy/PromEL/es v0.0.1
	google.golang.org/genproto v0.0.0-20200205142000-a86caf926a67 // indirect
	gopkg.in/yaml.v2 v2.2.8
)
