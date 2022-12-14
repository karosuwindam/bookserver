package config

import "github.com/caarlos0/env/v6"

type SetupServer struct {
	Protocol string `env:"PROTOCOL" envDefault:"tcp"`
	Hostname string `env:"WEB_HOST" envDefault:""`
	Port     string `env:"WEB_PORT" envDefault:"8080"`
}

type SetupSql struct {
	DBNAME     string `env:"DB_NAME" envDefault:"sqlite3"`
	DBHOST     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPORT     string `env:"DB_PORT" envDefault:"3306"`
	DBUSER     string `env:"DB_USER" envDefault:""`
	DBPASS     string `env:"DB_PASS" envDefault:""`
	DBFILE     string `env:"DB_FILE" envDefault:"test.db"`   //ファイル名
	DBROOTPASS string `env:"DB_ROOTPASS" envDefault:"./db/"` //相対パス
}

type SecretKey struct {
	JwtKey string `env:"JWT_KEY" envDefault:"SECRET_KEY"`
}

type Config struct {
	Server   *SetupServer
	Sql      *SetupSql
	SeretKey *SecretKey
}

//環境設定
func EnvRead() (*Config, error) {
	serverCfg := &SetupServer{}
	if err := env.Parse(serverCfg); err != nil {
		return nil, err
	}
	sqlCfg := &SetupSql{}
	if err := env.Parse(sqlCfg); err != nil {
		return nil, err
	}
	secretCfg := &SecretKey{}
	if err := env.Parse(secretCfg); err != nil {
		return nil, err
	}
	return &Config{
		Server:   serverCfg,
		Sql:      sqlCfg,
		SeretKey: secretCfg,
	}, nil

}
