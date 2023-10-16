package m4t_test

import (
	context "context"
	"io/ioutil"
	"testing"

	"github.com/nexptr/omnigram-server/api/m4t"
	"github.com/nexptr/omnigram-server/log"
	"go.uber.org/zap/zapcore"
	grpc "google.golang.org/grpc"
)

func TestManager_LoadConfig(t *testing.T) {

	log.Init(`stdout`, zapcore.DebugLevel)

	// _, filename, _, _ := runtime.Caller(0)

	// cf, err := conf.InitConfig(path.Join(path.Dir(filename), "../conf/conf.yaml"))

	// if err != nil {
	// 	t.Fatal(err)
	// }

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := m4t.NewTextToAudioClient(conn)

	// Send a request to the server
	req := &m4t.TextRequest{
		Text: "安如磐石~收！丝线--交织，古华深秘，裁！断！收！嗨！咻~走~咻~走~嗨~咻~吃饱喝饱一路走好！",
		// Text:    `超过一天的攻击记录不要使用此接口`,
		Lang:    `zh-cn`,
		AudioId: `female_001`,
	}

	response, err := client.ConvertTextToAudio(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to call SayHello: %v", err)
	}

	// write the whole body at once
	ioutil.WriteFile("output.wav", response.AudioData, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// // Print the response
	// log.I("Response: %s\n", response.Message)

	// openai :=

}
