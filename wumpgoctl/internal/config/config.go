package config

type Config struct {
	Meta      Meta
	Features  Features
	GoVersion string
}

type Features struct {
	Codegen bool
	Gateway Gateway
	HTTP    HTTP
}

type Gateway struct {
	Enabled bool
	NATS    bool
	Redis   bool
	Local   bool
}

type HTTP struct {
	Enabled bool
}

type Meta struct {
	Name        string
	Package     string
	Summary     string
	Description string
}
