package ghiccup

import (
	"log"
	"os"
	"path"
	"testing"
)

func createBuffer() (*MessageBuffer, string, error) {
	dbFile := path.Join(os.TempDir(), "testing.db")
	os.Remove(dbFile)
	buf, err := NewBuffer(dbFile)
	return buf, dbFile, err
}
func TestInitializeBuffer(t *testing.T) {
	buf, dbFile, err := createBuffer()
	if err != nil {
		t.Error(err)
	}
	defer buf.Close()
	_, err = os.Open(dbFile)
	if err != nil {
		t.Error(err)
	}
}

func TestAddSingleItem(t *testing.T) {
	buf, _, _ := createBuffer()
	defer buf.Close()
	err := buf.Add("Hello")
	if err != nil {
		t.Error(err)
	}
	count := 0
	buf.ForEach(func(o interface{}) error {
		log.Println("Received Message")
		count++
		return nil
	}, func() interface{} {
		var r string
		return &r
	})
	if count != 1 {
		t.Error("Didn't pop message")
	}
}
