package createzip

import (
	"os"
	"testing"
)

func TestCreateZip(t *testing.T) {
	t.Log("-------------- pdf to image ----------------")
	c, err := Setup(
		"./pdf", "testout.pdf",
		"./zip", "testout.zip",
		"/tmp",
	)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	if err := c.Pdftoimage(); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	os.Remove(c.Zippass + "/" + c.ZipName)
	t.Log("-------------- OK ----------------")
}
