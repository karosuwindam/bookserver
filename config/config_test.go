package config

import (
	"crypto/rand"
	"errors"
	"testing"
)

func TestEnvReadDefult(t *testing.T) {
	t.Log("----------------- EnvRead --------------------------")

	cfg, err := EnvRead()
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	t.Log("----------------- Check Server data --------------------------")
	if cfg.Server.Protocol != "tcp" {
		t.Errorf("Error Hostname %v = %v", cfg.Server.Protocol, "tcp")
		t.FailNow()
	}
	if cfg.Server.Hostname != "" {
		t.Errorf("Error Hostname %v = %v", cfg.Server.Hostname, "")
		t.FailNow()
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("Error Hostname %v = %v", cfg.Server.Port, "8080")
		t.FailNow()
	}
	t.Log("----------------- Check SQL data --------------------------")
	if cfg.Sql.DBNAME != "sqlite3" {
		t.Errorf("Error Hostname %v = %v", cfg.Sql.DBNAME, "sqlite3")
		t.FailNow()
	}
	if cfg.Sql.DBHOST != "127.0.0.1" {
		t.Errorf("Error Hostname %v = %v", cfg.Sql.DBHOST, "127.0.0.1")
		t.FailNow()
	}
	if cfg.Sql.DBPORT != "3306" {
		t.Errorf("Error Hostname %v = %v", cfg.Sql.DBPORT, "3306")
		t.FailNow()
	}
	if cfg.Sql.DBUSER != "" {
		t.Errorf("Error Hostname %v = %v", cfg.Sql.DBUSER, "")
		t.FailNow()
	}
	if cfg.Sql.DBPASS != "" {
		t.Errorf("Error Hostname %v = %v", cfg.Sql.DBPASS, "")
		t.FailNow()
	}
	if cfg.Sql.DBFILE != "test.db" {
		t.Errorf("Error Hostname %v = %v", cfg.Sql.DBFILE, "test.db")
		t.FailNow()
	}
	t.Log("----------------- Check Folder data --------------------------")
	if cfg.Folder.Tmp != "./tmp" {
		t.Errorf("Error Hostname %v = %v", cfg.Folder.Tmp, "./tmp")
		t.FailNow()
	}
	if cfg.Folder.Zip != "./upload/zip" {
		t.Errorf("Error Hostname %v = %v", cfg.Folder.Zip, "./upload/zip")
		t.FailNow()
	}
	if cfg.Folder.Pdf != "./upload/pdf" {
		t.Errorf("Error Hostname %v = %v", cfg.Folder.Pdf, "./upload/pdf")
		t.FailNow()

	}
	if cfg.Folder.Img != "./html/img" {
		t.Errorf("Error Hostname %v = %v", cfg.Folder.Img, "./upload/pdf")
		t.FailNow()

	}
	t.Log("----------------- Seclet key data --------------------------")
	if cfg.SeretKey.JwtKey != "SECRET_KEY" {
		t.Errorf("Error Hostname %v = %v", cfg.SeretKey.JwtKey, "SECRET_KEY")
		t.FailNow()
	}

	t.Log("----------------- EnvRead OK --------------------------")
}

func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func TestEnvReadServer(t *testing.T) {
	setupEnv := map[string]string{}
	aEnv := []string{"PROTOCOL", "WEB_HOST", "WEB_PORT"}
	input := []string{}
	for _, env := range aEnv {
		tmp, _ := makeRandomStr(10)
		setupEnv[env] = tmp
		input = append(input, tmp)
	}
	t.Log("----------------- EnvRead setup server --------------------------")
	for key, value := range setupEnv {
		t.Setenv(key, value)
	}

	cfg, err := EnvRead()
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	ckcfg := []string{cfg.Server.Protocol, cfg.Server.Hostname, cfg.Server.Port}

	for i := 0; i < len(ckcfg); i++ {
		if ckcfg[i] != input[i] {
			t.Errorf("Error %v %v = %v", aEnv[i], ckcfg[i], input[i])
			t.FailNow()
		}

	}
	t.Log(aEnv)
	t.Log(ckcfg)
	t.Log("----------------- EnvRead OK --------------------------")

}

func TestEnvReadSQL(t *testing.T) {
	setupEnv := map[string]string{}
	aEnv := []string{"DB_NAME", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_FILE"}
	input := []string{}
	for _, env := range aEnv {
		tmp, _ := makeRandomStr(10)
		setupEnv[env] = tmp
		input = append(input, tmp)
	}
	t.Log("----------------- EnvRead setup sql --------------------------")
	for key, value := range setupEnv {
		t.Setenv(key, value)
	}

	cfg, err := EnvRead()
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	ckcfg := []string{
		cfg.Sql.DBNAME, cfg.Sql.DBHOST,
		cfg.Sql.DBPORT, cfg.Sql.DBUSER,
		cfg.Sql.DBPASS, cfg.Sql.DBFILE,
	}

	for i := 0; i < len(ckcfg); i++ {
		if ckcfg[i] != input[i] {
			t.Errorf("Error %v %v = %v", aEnv[i], ckcfg[i], input[i])
			t.FailNow()
		}

	}
	t.Log(aEnv)
	t.Log(ckcfg)
	t.Log("----------------- EnvRead OK --------------------------")

}

func TestEnvReadSecletKey(t *testing.T) {
	setupEnv := map[string]string{}
	aEnv := []string{"JWT_KEY"}
	input := []string{}
	for _, env := range aEnv {
		tmp, _ := makeRandomStr(10)
		setupEnv[env] = tmp
		input = append(input, tmp)
	}
	t.Log("----------------- EnvRead setup seclet key --------------------------")
	for key, value := range setupEnv {
		t.Setenv(key, value)
	}

	cfg, err := EnvRead()
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	ckcfg := []string{
		cfg.SeretKey.JwtKey,
	}

	for i := 0; i < len(ckcfg); i++ {
		if ckcfg[i] != input[i] {
			t.Errorf("Error %v %v = %v", aEnv[i], ckcfg[i], input[i])
			t.FailNow()
		}

	}
	t.Log(aEnv)
	t.Log(ckcfg)
	t.Log("----------------- EnvRead OK --------------------------")

}
