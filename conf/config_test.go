package conf

import (
	"encoding/json"
	"testing"
)

func TestInitConfig(t *testing.T) {

	got, err := InitConfig(`../../examples/app/conf.yaml`)

	if err != nil {
		t.Fatal(err)
	}

	j, _ := json.Marshal(got)

	println(string(j))

}
