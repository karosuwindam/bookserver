package table

import (
	"encoding/json"
	"errors"
	"math/rand"
	"testing"
)

func TestTablelistSetup(t *testing.T) {
	tablelistsetup()
	if !ckType(&Booknames{}) {
		t.FailNow()
	}
	if !ckType(&Copyfile{}) {
		t.FailNow()
	}
	if !ckType(&Filelists{}) {
		t.FailNow()
	}
	for tname, _ := range tablelist {
		readBaseCreate(tname)
	}

}

func TestCreateSerchKeyword(t *testing.T) {
	tablelistsetup()
	booksrdata := createSerchText(BOOKNAME, "a")
	ckbookdata := []string{"booktype", "burand", "ext", "name", "title", "writer"}
	for i := 1; i < len(ckbookdata); i++ {
		if booksrdata[ckbookdata[i]] != booksrdata[ckbookdata[i-1]] {
			t.Fatalf("%v != %v", booksrdata[ckbookdata[i]], booksrdata[ckbookdata[i-1]])
			t.FailNow()
		}
	}
	t.Log(booksrdata)
	copysrdata := createSerchText(COPYFILE, "b")
	ckcopydata := []string{"zippass"}
	for i := 1; i < len(ckcopydata); i++ {
		if copysrdata[ckcopydata[i]] != copysrdata[ckcopydata[i-1]] {
			t.Fatalf("%v != %v", copysrdata[ckcopydata[i]], copysrdata[ckcopydata[i-1]])
			t.FailNow()
		}
	}
	t.Log(copysrdata)
	filesrdata := createSerchText(FILELIST, "c")
	ckfiledata := []string{"name", "pdfpass", "tag", "zippass"}
	for i := 1; i < len(ckfiledata); i++ {
		if filesrdata[ckfiledata[i]] != filesrdata[ckfiledata[i-1]] {
			t.Fatalf("%v != %v", filesrdata[ckfiledata[i]], filesrdata[ckfiledata[i-1]])
			t.FailNow()
		}
	}
	t.Log(filesrdata)

}

func TestJsonToStruct(t *testing.T) {
	jsondata := readfile("./jsonsample/booknames.json")
	if tmp := JsonToStruct(BOOKNAME, []byte(jsondata)); tmp == nil {
		t.FailNow()
	} else {
		_, ok := tmp.(Booknames)
		if !ok {
			t.FailNow()
		}
	}
	jsondata = readfile("./jsonsample/copyfile.json")
	if tmp := JsonToStruct(COPYFILE, []byte(jsondata)); tmp == nil {
		t.FailNow()
	} else {
		_, ok := tmp.(Copyfile)
		if !ok {
			t.FailNow()
		}
	}
	jsondata = readfile("./jsonsample/filelist.json")
	if tmp := JsonToStruct(FILELIST, []byte(jsondata)); tmp == nil {
		t.FailNow()
	} else {
		_, ok := tmp.(Filelists)
		if !ok {
			t.FailNow()
		}
	}
}

func TestRandGenerate(t *testing.T) {
	tmp := []Booknames{}
	t.Log("--------------------- rond check max down -------------------------")
	for i := 0; i < RAND_MAX; i++ {
		str, _ := MakeRandomStr(10)
		rtmp := Booknames{
			Id:   i + 1,
			Name: str,
		}
		tmp = append(tmp, rtmp)
	}
	jout1 := RandGenerate(tmp)
	if jout, _ := json.Marshal(tmp); string(jout) != jout1 {
		t.Fatalf("%v!=%v", string(jout), jout1)
		t.FailNow()
	}
	tmpf := []Filelists{}
	for i := 0; i < RAND_MAX+20; i++ {
		str, _ := MakeRandomStr(10)
		rtmp := Filelists{
			Id:   i + 1,
			Name: str,
		}
		tmpf = append(tmpf, rtmp)
	}
	t.Log("------------------- rond check max down OK---------------------------")
	t.Log("------------------- rond check max over ---------------------------")
	jout1 = RandGenerate(tmpf)
	if jt, ok := JsonToStruct(FILELIST, []byte(jout1)).([]Filelists); ok {
		if len(jt) != RAND_MAX {
			t.Fatalf("%v,%v", len(jt), jout1)
			t.FailNow()
		} else {
			for _, arr := range jt {
				flag := true
				for _, arr2 := range tmpf {
					if arr.Id == arr2.Id && arr.Name == arr2.Name {
						flag = false
						break
					}
				}
				if flag {
					t.Fatalf("not data:%v+", arr)
					t.FailNow()
				}
			}
		}
	} else {
		t.FailNow()
	}
	t.Log("------------------ rond check max over OK ----------------------------")

}

func MakeRandomStr(digit uint32) (string, error) {
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
