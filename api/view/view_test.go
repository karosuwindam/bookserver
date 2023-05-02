package view

import (
	"bookserver/config"
	"bookserver/webserver"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestViewList(t *testing.T) {

	t.Setenv("ZIP_FILEPASS", "./zip")
	cfg, _ := config.EnvRead()
	ss, _ := webserver.NewSetup(cfg)
	web, _ := Setup(cfg)
	webserver.Config(ss, web, "")
	s, _ := ss.NewServer()
	t.Log("----------------- View List Start --------------------------")

	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return s.Run(ctx)
	})
	getUrldata("http://localhost:8080/view/1", t)
	cancel()
	t.Log("----------------- View List End --------------------------")
}

/*---------------確認用関数----------------*/
func getUrldata(url string, t *testing.T) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	if resp.StatusCode != http.StatusOK {
		t.Fail()
	}
	fmt.Printf("%s\n", b)

}
