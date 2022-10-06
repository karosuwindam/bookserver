package webserver

import (
	"bookserver/config"
	"bookserver/message"
	"bookserver/webserver/common"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"
)

// SetupServer
// サーバ動作の設定
type SetupServer struct {
	protocol string // Webサーバーのプロトコル
	hostname string //Webサーバのホスト名
	port     string //Webサーバの解放ポート

	routefunc map[string]func(
		interface{}, http.ResponseWriter, *http.Request,
	) // /v1/{data}による実行関数
	routeinterface map[string]interface{}     // /v1/{data}による実行関数の入力データ
	accessmpa      map[string]common.UserType // /v1/{data}によるアクセス制御

	mux *http.ServeMux //webサーバのmux
}

// Server
// Webサーバの管理情報
type Server struct {
	// Webサーバの管理関数
	Srv *http.Server
	// 解放の管理関数
	L net.Listener
}

// Status
// ToDo
type Status struct {
	// 	Status string `json:status`
}

// NewSetup(*config.Config) = *SetupServer,error
//
// Webサーバ設定の初期化関数
//
// data(*config.Config) : Env設定で読みだした設定
func NewSetup(data *config.Config) (*SetupServer, error) {
	cfg := &SetupServer{
		protocol: data.Server.Protocol,
		hostname: data.Server.Hostname,
		port:     data.Server.Port,
	}
	cfg.mux = http.NewServeMux()
	cfg.routefunc = map[string]func(interface{}, http.ResponseWriter, *http.Request){}
	cfg.routeinterface = map[string]interface{}{}
	cfg.accessmpa = map[string]common.UserType{}
	cfg.mux.HandleFunc("/v1/", cfg.v1)

	cfg.route(data)

	return cfg, nil
}

// (*SetupServer) NewServer() = *Server,error
//
// Webサーバの開始設定
func (t *SetupServer) NewServer() (*Server, error) {
	message.Println("Setupserver", t.protocol, t.hostname+":"+t.port)
	l, err := net.Listen(t.protocol, t.hostname+":"+t.port)
	if err != nil {
		return nil, err
	}
	return &Server{
		Srv: &http.Server{Handler: t.muxHandler()},
		L:   l,
	}, nil
}

// (*SetupServer) AddV1(route, sdata, funcdata) = error
//
// /v1/{route}による処理関数を紐づける紐づけられない場合はエラーが返る
//
// route(string) : URLルートパス /v1/として紐づける
// sdata(interface{}) : 関数に引き渡すポインタ情報
// funcdata(func(interface{}, http.ResponseWriter, *http.Request)) : 処理実行関数
func (t *SetupServer) AddV1(usertype common.UserType, route string, sdata interface{}, funcdata func(interface{}, http.ResponseWriter, *http.Request)) error {
	if reflect.TypeOf(sdata).Kind() != reflect.Ptr {
		return errors.New("sdata is not pointer")
	}
	var v1route string
	if route[0] != "/"[0] {
		v1route = "/v1" + "/" + route
	} else {
		v1route = "/v1" + route
	}

	if t.routefunc[v1route] != nil {
		return errors.New("")
	}
	t.routefunc[v1route] = funcdata
	t.routeinterface[v1route] = sdata
	t.accessmpa[v1route] |= usertype
	return nil
}

// (*SetupServer) AddV1(route, handler)
//
// http上に定義されたハンドラ関数を紐づける
//
// route(string) : ホームからのURLルートパス
// handler(func(http.ResponseWriter, *http.Request)) : httpの関数処理
func (t *SetupServer) Add(route string, handler func(http.ResponseWriter, *http.Request)) {
	t.mux.HandleFunc(route, handler)
}

// (*SetupServer) muxHandler()
// SetupServer内のmuxhandlerを返す関数
func (t *SetupServer) muxHandler() http.Handler { return t.mux }

// v1OtherBuck(http.ResponseWriter, *http.Request)
//
// /v1/によるURLパスで定義されていないときの返却処理関数
func v1OtherBuck(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	methd := r.Method
	msg := message.Message{Name: "v1", Status: methd + ":" + url, Code: http.StatusOK}
	rst := message.Result{Name: "v1", Code: http.StatusOK, Option: "", Date: time.Now(), Result: msg}
	fmt.Fprintf(w, "%v", rst.Output())
}

// (*SetupServer) v1(http.ResponseWriter, *http.Request)
//
// /v1/にアクセスした時に判断する処理
func (t *SetupServer) v1(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	flag := false
	user := common.ADMIN | common.GUEST | common.USER //ユーザの値
	var check string
	for turl := range t.routeinterface {
		if len(url) >= len(turl) {
			if turl == url[:len(turl)] {
				check = turl
				flag = true
				break
			}
		}
	}
	if flag && t.accessmpa[check]|user > 0 { //登録済み
		t.routefunc[check](t.routeinterface[check], w, r)
	} else { //日登録
		v1OtherBuck(w, r)
	}

}
