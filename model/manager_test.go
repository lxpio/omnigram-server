package model_test

import (
	"path"
	"runtime"
	"testing"

	"github.com/nexptr/llmchain/llms/openai"
	"github.com/nexptr/omnigram-server/conf"
)

func TestManager_LoadConfig(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)

	cf, err := conf.InitConfig(path.Join(path.Dir(filename), "../conf/conf.yaml"))

	if err != nil {
		t.Fatal(err)
	}

	for _, o := range cf.ModelOptions {

		if o.Name == `vicuna-13b-v1.1` {

			openai, err := openai.FromYaml(o)

			if err != nil {
				t.Fatal(err)
			}

			println(openai.String())

		}

	}

	// openai :=

}
