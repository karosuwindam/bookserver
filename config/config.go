package config

import (
	"bufio"
	"os"
	"regexp"

	"github.com/caarlos0/env/v6"
)

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

type UploadCfg struct {
	MAX_MULTI_MEMORY string `env:"MAX_MULTI_MEMORY" envDefault:"256M"`
}

type SetupFolder struct {
	Tmp    string `env:"TMP_FILEPASS" envDefault:"./tmp"`        //画像を一時保存するパス
	Img    string `env:"IMG_FILEPASS" envDefault:"./html/img"`   //1ページ目の画像ファイルを保存するフォルダ
	Pdf    string `env:"PDF_FILEPASS" envDefault:"./upload/pdf"` //PDFのアップロード先フォルダ
	Zip    string `env:"ZIP_FILEPASS" envDefault:"./upload/zip"` //ZIPのアップロード先フォルダ
	Public string `env:"PUBLIC_FILEPASS" envDefault:"./public"`  //公開用のフォルダ
}

type SecretKey struct {
	JwtKey string `env:"JWT_KEY" envDefault:"SECRET_KEY"`
}

type Config struct {
	Server   *SetupServer
	Sql      *SetupSql
	SeretKey *SecretKey
	Folder   *SetupFolder
	Upload   *UploadCfg
	Version  string
}

// バージョン読み取り
func versionRead() string {
	f, err := os.Open("version")
	output := "0.0.1"
	if err != nil {
		return output
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rg := regexp.MustCompile(`(.*).(.*).(.*)`)
		tmp := scanner.Text()
		if rg.MatchString(tmp) {
			output = tmp
			break
		}
	}
	return output
}

// 環境設定
func EnvRead() (*Config, error) {
	serverCfg := &SetupServer{}
	if err := env.Parse(serverCfg); err != nil {
		return nil, err
	}
	sqlCfg := &SetupSql{}
	if err := env.Parse(sqlCfg); err != nil {
		return nil, err
	}
	folderCfg := &SetupFolder{}
	if err := env.Parse(folderCfg); err != nil {
		return nil, err
	}
	secretCfg := &SecretKey{}
	if err := env.Parse(secretCfg); err != nil {
		return nil, err
	}
	uploadCfg := &UploadCfg{}
	if err := env.Parse(uploadCfg); err != nil {
		return nil, err
	}
	return &Config{
		Server:   serverCfg,
		Sql:      sqlCfg,
		Folder:   folderCfg,
		SeretKey: secretCfg,
		Upload:   uploadCfg,
		Version:  versionRead(),
	}, nil

}
