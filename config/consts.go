package config

const (
	VERSION      = "0.0.11"
	PORT         = 3444
	SERVICE_NAME = "ps-chartdata"

	WORKSTATION = "workstation"
	DEV         = "dev"
	STAGE       = "stage"
	PROD        = "prod"

	CONSUL_KV           = "ps"
	CONSUL_HOST_DEV     = "localhost"
	CONSUL_PORT_DEV     = "8500"
	CONSUL_HOST_CLUSTER = "consul-server"
	CONSUL_PORT_CLUSTER = "8500"
)

var (
	MONGOHOSTS_WORKSTATION = []string{"localhost:27017"}
)
