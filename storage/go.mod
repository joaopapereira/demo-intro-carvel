module storage

go 1.20

require (
	github.com/bufbuild/connect-go v1.5.2
	github.com/rs/cors v1.8.3
	backend v0.0.0
)

require google.golang.org/protobuf v1.28.1 // indirect

replace "backend" v0.0.0 => "../backend"

