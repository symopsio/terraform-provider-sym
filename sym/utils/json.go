package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func JsonBytesEqual(b1 []byte, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

// SuppressEquivalentJsonDiffs is a DiffSuppressFunc that can be passed into
// a schema to ensure meaningless differences in JSON are ignored when reading/writing.
//
// e.g. without this function, TF will see the following two blobs as different, and
// prompt an update. With this function set, the blobs are seen as equivalent:
//
//  {"template": "sym:approval", "name": "access"}
//  {"name": "access", "template": "sym:approval"}
func SuppressEquivalentJsonDiffs(k string, old string, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return JsonBytesEqual(ob.Bytes(), nb.Bytes())
}

// SuppressEquivalentFileContentDiffs is a DiffSuppressFunc that can be passed into
// a schema when the provider receives file content from the API but needs to compare
// against the file content locally, but the terraform specifies the file name vs. content.
//
// e.g. the API response vs. the terraform, we get something like
// implementation = "ZnJvbSBzeW0uYW5ub3RhdGlvbnMgaW1wb3J0"
// implementation = "impl.py"
//
// If the file fails to be read, we will return false to specify that we should not suppress
// the difference.
// If the incoming base64 encoding file fails to decode, we will return false.
func SuppressEquivalentFileContentDiffs(k string, old string, new string, d *schema.ResourceData) bool {
	log.Printf("\n===== HERE SUPPRESS ====\n")
	log.Printf("k: %v", k)
	log.Printf("old: %v", old)
	log.Printf("new: %v", new)

	newBytes, err := ioutil.ReadFile(new)
	if err != nil {
		return false
	}

	oldBytes, err := base64.StdEncoding.DecodeString(old)
	if err != nil {
		return false
	}

	return bytes.Compare(newBytes, oldBytes) == 0
}
