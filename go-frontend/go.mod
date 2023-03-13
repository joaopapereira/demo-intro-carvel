module frontend

go 1.20

require (
	backend v0.0.0
	github.com/bufbuild/connect-go v1.5.2
	github.com/rs/cors v1.8.3
)

require google.golang.org/protobuf v1.28.1 // indirect

replace backend v0.0.0 => ../backend
