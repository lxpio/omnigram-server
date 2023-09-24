package schema_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nexptr/omnigram-server/epub/schema"
	"github.com/nexptr/omnigram-server/store"
)

func TestGetReadProcessBooks(t *testing.T) {
	//获取当前测试文件路径
	_, filename, _, _ := runtime.Caller(0)

	db, _ := store.OpenDB(&store.Opt{
		Driver: store.DRPostgres,
		Host:   filepath.Join(filename, "../"),
	})

	data, err := schema.ReadingBooks(db, 1, 10)

	if err != nil {
		t.Error(err)
	}

	// println(data.Total)

	println(data[0].Title)

}
