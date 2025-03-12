package config

// ServerCfg server
type ServerCfg struct {
	ServerInfo ServerConfig `mapstructure:"server" toml:"server"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" toml:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" toml:"consul"`
	RedisInfo  RedisConfig  `mapstructure:"redis" toml:"redis"`
}

type ServerConfig struct {
	Host string `mapstructure:"host" toml:"host"`
	Port int    `mapstructure:"port" toml:"port"`
	Name string `mapstructure:"name" toml:"name"`
}

// MysqlConfig MysqlInfo
type MysqlConfig struct {
	Host     string `mapstructure:"host" toml:"host"`
	Port     int    `mapstructure:"port" toml:"port"`
	Dbname   string `mapstructure:"dbName" toml:"dbname"`
	Username string `mapstructure:"username" toml:"username"`
	Password string `mapstructure:"password" toml:"password"`
}

type ConsulConfig struct {
	Host string   `mapstructure:"host" toml:"host"`
	Port int      `mapstructure:"port" toml:"port"`
	Name string   `mapstructure:"name" toml:"name"`
	Id   string   `mapstructure:"id" toml:"id"`
	Tag  []string `mapstructure:"tag" toml:"tag"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Scheme    string `mapstructure:"scheme"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
