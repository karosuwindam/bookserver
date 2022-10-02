package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname string `env:"WEB_HOST" envDefault:""`
	Port     string `env:"WEB_PORT" envDefault:"8080"`
}

type SetupSql struct {
	DBNAME string `env:"DB_NAME" envDefault:"mysql"`
	DBHOST string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPORT string `env:"DB_PORT" envDefault:"3306"`
	DBUSER string `env:"DB_USER" envDefault:""`
	DBPASS string `env:"DB_PASS" envDefault:""`
	DBFILE string `env:"DB_FILE" envDefault:"test.db"`
}

type Config struct {
	Server *SetupServer
	Sql    *SetupSql
}

//環境設定
func Envread() (*Config, error) {
	server_cfg := &SetupServer{}
	if err := env.Parse(server_cfg); err != nil {
		return nil, err
	}
	sql_cfg := &SetupSql{}
	if err := env.Parse(sql_cfg); err != nil {
		return nil, err
	}
	return &Config{
		Server: server_cfg,
		Sql:    sql_cfg,
	}, nil

}
