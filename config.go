package main

func newConfig(name, version string) *Config {
	return &Config{
		App:    &AppConfig{},
		Logger: &LoggerConfig{},
		Core:   &CoreConfig{},
		Server: &ServerConfig{Name: name, Version: version},
		Consul: &ConsulConfig{
			Path: "product-domains/mc-go-fns/service",
		},
		Vault:    &VaultConfig{},
		Metric:   &MetricConfig{},
		Database: &DatabaseConfig{},
		Broker:   &BrokerConfig{},
	}
}

type Config struct {
	App      *AppConfig      `json:"app"`
	Logger   *LoggerConfig   `json:"logger"`
	Core     *CoreConfig     `json:"core"`
	Server   *ServerConfig   `json:"server"`
	Consul   *ConsulConfig   `json:"consul"`
	Vault    *VaultConfig    `json:"vault"`
	Metric   *MetricConfig   `json:"metric"`
	Database *DatabaseConfig `json:"database"`
	Broker   *BrokerConfig   `json:"broker"`
}

type LoggerConfig struct {
	Level string `json:"level"`
}

type CoreConfig struct {
	Profile bool `json:"profile"`
}

type ServerConfig struct {
	Name    string `json:"name"`
	Version string `json:"-"`
	Address string `json:"address" default:":9090"`
}

type AppConfig struct {
	FnsToken   string `json:"fns_token"`   // jwt token
	FnsAddress string `json:"fns_address"` // address schema://host:port
	MainTopic  string `json:"main_topic" default:"main_topic"`
	ErrorTopic string `json:"error_topic" default:"error_topic"`
}

type ConsulConfig struct {
	Address string `json:"address" default:"http://localhost:12345"`
	Path    string `json:"path"`
	Token   string `env:"CONSUL_TOKEN" json:"-"`
}

type VaultConfig struct {
	Address string `json:"address" default:"http://localhost:54321"`
	Path    string `json:"path"`
	Token   string `env:"VAULT_TOKEN" json:"-"`
}

type MetricConfig struct {
	Address string `json:"address" default:":8080"`
}

type DatabaseConfig struct {
	Address         string `json:"address"`
	Login           string `json:"login"`
	Passw           string `json:"passw"`
	Name            string `json:"name" default:"file:database.db?cache=shared&mode=memory"`
	ConnMax         int64  `json:"conn_max"`
	ConnLifetime    int64  `json:"conn_lifetime"`
	ConnMaxIdleTime int64  `json:"conn_max_idletime"`
	MigrateUp       bool   `json:"migrate_up" flag:"name=migrate-up,desc='database migrations up',default='false'"`
	MigrateDown     bool   `json:"migrate_down" flag:"name=migrate-down,desc='database migrations down',default='false'"`
}

type BrokerConfig struct {
	Address []string `json:"address"`
	Login   string   `json:"login"`
	Passw   string   `json:"passw"`
}
