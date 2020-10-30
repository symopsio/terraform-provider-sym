package main

import (
	"encoding/base64"
	"fmt"
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
			Name: "sym:approval:1.0",
			Version: &models.Version{
				Major: 1,
			},
		},
	}
	bytes, _ := proto.Marshal(flow)
	enc := base64.StdEncoding.EncodeToString(bytes)
	tag := "sym.tf.models.Flow;template"
	sep := "---FIELD_SEP---"
	repr := fmt.Sprintf("%s\n%s\n%s", tag, sep, enc)
	fmt.Printf(repr)


	ioutil.WriteFile("/Users/rick/testing.proto", []byte(repr), 0644)
}
