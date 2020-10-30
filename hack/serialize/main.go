package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/symopsio/protos/go/tf/models"
	"io/ioutil"
)

func main() {
	flow := &models.Flow{
		Name: "test:flow",
		Version: &models.Version{
			Major: 1,
		},
		Uuid: "bd6b69bd-0d93-463e-b997-b19a8370da6e",
		Template: &models.Template{
			Name: "template",
		},
	}
	bytes, _ := proto.Marshal(flow)
	ioutil.WriteFile("/Users/rick/testing.proto", bytes, 0644)
}
