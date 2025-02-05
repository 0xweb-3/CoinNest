package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DbName   string `mapstructure:"db_name" json:"db_name"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

//type ConsulConfig struct {
//	Host string `mapstructure:"host" json:"host"`
//	Port int    `mapstructure:"port" json:"port"`
//}
//
//type RedisConfig struct {
//	Host     string `mapstructure:"host" json:"host"`
//	Port     int    `mapstructure:"port" json:"port"`
//	Password string `mapstructure:"password" json:"password"`
//}

type Config struct {
	Name  string      `mapstructure:"name" json:"name"`
	Port  int         `mapstructure:"port" json:"port"`
	Host  string      `mapstructure:"host" json:"host"`
	Mysql MysqlConfig `mapstructure:"mysql" json:"mysql"`
	//Redis  RedisConfig  `mapstructure:"redis" json:"redis"`
	//Consul ConsulConfig `mapstructure:"consul" json:"consul"`
}
