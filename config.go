package nebulaorm

import (
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"time"
)

// Config  config
type Config struct {
	// Username to connect to nebula graph server
	Username string `json:"username" yaml:"username"`

	// Password to connect to nebula graph server
	Password string `json:"password" yaml:"password"`

	// SpaceName name of the graph space used for the current connection
	SpaceName string `json:"space_name" yaml:"space_name"`

	// Addresses server address list，host:port
	Addresses []string `json:"addresses" yaml:"addresses"`

	// Timeout connection dail read-write timeout
	ConnTimeout time.Duration `json:"conn_timeout" yaml:"conn_timeout"`

	// ConnMaxIdleTime connection max idle time
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`

	// MaxOpenConns max number of connections in the connection pool
	MaxOpenConns int `json:"max_open_conns" yaml:"max_open_conns"`

	// MinOpenConns min number of connections in the connection pool
	MinOpenConns int `json:"min_open_conns" yaml:"min_open_conns"`

	// TimezoneName time zone name, default is Local, if the nebula graph server is configured in a time zone different from Local,
	// you need to change it to the same configuration as the nebula graph server.
	TimezoneName string `json:"timezone_name" yaml:"timezone_name"`

	// nebulaSessionOpts nebula session pool config
	nebulaSessionOpts []nebula.SessionPoolConfOption

	timezone *time.Location
}

type ConfigOption interface {
	apply(*Config)
}

type funcConfigOption func(*Config)

func (f funcConfigOption) apply(conf *Config) {
	f(conf)
}

// WithNebulaSessionPoolOptions customizing nebula session pool parameters
func WithNebulaSessionPoolOptions(opts []nebula.SessionPoolConfOption) ConfigOption {
	return funcConfigOption(func(config *Config) {
		config.nebulaSessionOpts = append(config.nebulaSessionOpts, opts...)
	})
}
