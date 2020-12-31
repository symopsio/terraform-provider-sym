package utils

import (
	"bytes"
	"encoding/json"
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
