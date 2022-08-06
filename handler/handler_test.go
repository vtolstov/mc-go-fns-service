package handler

import (
	"testing"

	pb "github.com/vtolstov/mc-go-fns-proto"
	jsonpbcodec "go.unistack.org/micro-codec-jsonpb/v3"
)

func TestResponse(t *testing.T) {
	data := []byte(`{
	}`)
	result := &pb.GetInnRsp{}

	c := jsonpbcodec.NewCodec()

	if err := c.Unmarshal(data, result); err != nil {
		t.Fatal(err)
	}
}
