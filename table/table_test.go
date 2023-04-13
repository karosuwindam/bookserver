package table

import (
	"bookserver/config"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestTableSetup(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	sql, err := Setup(cfg)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	defer sql.Close()
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)

}

func TestTableReadAll(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	t.Setenv("DB_FILE", "test-read.db")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	if jsond, err := sql.ReadAll(BOOKNAME); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
	if jsond, err := sql.ReadAll(COPYFILE); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
	if jsond, err := sql.ReadAll(FILELIST); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
}

func TestTableName(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	t.Setenv("DB_FILE", "test-read.db")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	t.Log("check: test")
	if jsond, err := sql.ReadName(BOOKNAME, "test"); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
	t.Log("check: tt")
	if jsond, err := sql.ReadName(BOOKNAME, "tt"); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
}

func TestTableSerch(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	t.Setenv("DB_FILE", "test-read.db")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	if jsond, err := sql.Search(BOOKNAME, "t"); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
	if jsond, err := sql.Search(COPYFILE, "bb"); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
	if jsond, err := sql.Search(FILELIST, "bbb"); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(jsond)
	}
}

func TestTableReadID(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)
	writebook := []Booknames{}
	for i := 0; i < 20; i++ {
		tmp := Booknames{}
		if str, err := makeRandomStr(10); err == nil {
			tmp.Title = str
		}
		writebook = append(writebook, tmp)
	}
	if err := sql.Add(BOOKNAME, &writebook); err != nil {
		t.Log((err.Error()))
		t.FailNow()
	}
	rand.Seed(time.Now().UnixNano())
	tid := rand.Intn(21)
	for i := 0; i < 3; i++ {

		if jsond, err := sql.ReadID(BOOKNAME, tid); err != nil {
			t.Fatal(err.Error())
			t.FailNow()
		} else {
			t.Log(jsond)
			tjson, _ := json.Marshal(writebook[tid-1])
			if jsond != "["+string(tjson)+"]" {
				t.Fatal("err write log")
			}
		}
		tid = rand.Intn(20)
		if tid == 0 {
			tid = 1
		}
	}

}

func TestTableWrite(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)
	jsondata := readfile("./jsonsample/booknames.json")
	writebook := Booknames{}
	json.Unmarshal([]byte(jsondata), &writebook)
	if err := sql.Add(BOOKNAME, &writebook); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	jsondata = readfile("./jsonsample/copyfile.json")
	writeCopy := Copyfile{}
	json.Unmarshal([]byte(jsondata), &writeCopy)
	if err := sql.Add(COPYFILE, &writeCopy); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	jsondata = readfile("./jsonsample/filelist.json")
	writeFile := Filelists{}
	json.Unmarshal([]byte(jsondata), &writeFile)
	if err := sql.Add(FILELIST, &writeFile); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

}

func readfile(path string) string {
	var output string
	fp, err := os.Open(path)
	if err != nil {
		log.Panic(err)
		return ""
	}
	defer fp.Close()
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		output += string(buf[:n])
	}
	return output

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
