package config

import (
	"bufio"
	"os"
	"regexp"

	"github.com/caarlos0/env/v6"
)

// Webサーバの設定
type WebConfig struct {
	Protocol   string `env:"WEB_PROTOCOL" envDefault:"tcp"`  //接続プロトコル
	Hostname   string `env:"WEB_HOST" envDefault:""`         //接続DNS名
	Port       string `env:"WEB_PORT" envDefault:"8080"`     //接続ポート
	StaticPage string `env:"WEB_FOLDER" envDefault:"./html"` //静的ページの参照先
}

type SetupSql struct {
	DBNAME     string `env:"DB_NAME" envDefault:"sqlite3"`   //SQLの種類
	DBHOST     string `env:"DB_HOST" envDefault:"127.0.0.1"` //接続先のip
	DBPORT     string `env:"DB_PORT" envDefault:"3306"`      //接続ポート
	DBUSER     string `env:"DB_USER" envDefault:""`          //接続ユーザ名
	DBPASS     string `env:"DB_PASS" envDefault:""`          //接続のパス
	DBFILE     string `env:"DB_FILE" envDefault:"test.db"`   //ファイル名
	DBROOTPASS string `env:"DB_ROOTPASS" envDefault:"./db/"` //相対パス
}

// 図書サーバで使用するフォルダ設定
type BookserverConfig struct {
	Tmp              string `env:"TMP_FILEPASS" envDefault:"./tmp"`        //画像を一時保存するパス
	Img              string `env:"IMG_FILEPASS" envDefault:"./html/img"`   //1ページ目の画像ファイルを保存するフォルダ
	Pdf              string `env:"PDF_FILEPASS" envDefault:"./upload/pdf"` //PDFのアップロード先フォルダ
	Zip              string `env:"ZIP_FILEPASS" envDefault:"./upload/zip"` //ZIPのアップロード先フォルダ
	Public           string `env:"PUBLIC_FILEPASS" envDefault:"./public"`  //ファイル共有で使用するフォルダ
	MAX_MULTI_MEMORY string `env:"MAX_MULTI_MEMORY" envDefault:"512M"`     //アップロード時のメモリ制限
	RandamRead       int    `env:"RAND_COUNT" envDefault:"5"`              //ランダムの読み取り最大数
	ConvertCountMax  int    `env:"MAX_CONVERT_COUNT" envDefault:"3"`       //アップロードされたファイルの失敗回数
	HistoryMax       int    `env:"MAX_HISTORY" envDefault:"100"`           //履歴を読み取る際の最大数
}

// 認証関連の設定ファイル
type AuthConfig struct {
	JwtKey string `env:"JWT_KEY" envDefault:"SECRET_KEY"`
}

type LogConfig struct {
	Colors bool `env:"LOG_COLORS" envDefault:"true"`
}

type TracerData struct {
	// GrpcURL string `env:"TRACER_GRPC_URL" envDefault:"localhost:4317"`
	GrpcURL     string `env:"TRACER_GRPC_URL" envDefault:"otel-grpc.bookserver.home:4317"`
	ServiceName string `env:"TRACER_SERVICE_URL" envDefault:"bookserver-test"`
	TracerUse   bool   `env:"TRACER_ON" envDefault:"true"`
}

// versionファイルによるバージョン読み取り
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

var Web WebConfig
var DB SetupSql
var BScfg BookserverConfig
var Version string
var Auth AuthConfig
var Log LogConfig
var TraData TracerData

// 環境設定
func Init() error {
	if err := env.Parse(&Web); err != nil {
		return err
	}
	if err := env.Parse(&DB); err != nil {
		return err
	}
	if err := env.Parse(&BScfg); err != nil {
		return err
	}
	if err := env.Parse(&Auth); err != nil {
		return err
	}
	if err := env.Parse(&Log); err != nil {
		return err
	}
	if err := env.Parse(&TraData); err != nil {
		return err
	}
	Version = versionRead()
	return nil
}
