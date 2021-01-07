package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"log"
)

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
		log.Printf("Error decoding json %v for value %v", old, k)
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		log.Printf("Error decoding json %v for value %v", new, k)
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
//
// NOTE: if true is returned, the diff is suppressed and the field whose diff is being suppressed
// will have a value matching the OLD input. This means that if the diff is suppressed, we need to account
// for already having a base64 encoded string as the field's value vs. a file path.
func SuppressEquivalentFileContentDiffs(k string, old string, new string, d *schema.ResourceData) bool {
	if old == "" {
		return false
	}

	newBytes, err := ioutil.ReadFile(new)
	if err != nil {
		log.Printf("Error reading file %v for value %v", new, k)
		return false
	}

	oldBytes, err := base64.StdEncoding.DecodeString(old)
	if err != nil {
		log.Printf("Error decoding file content from API %v for value %v", old, k)
		return false
	}

	return bytes.Compare(newBytes, oldBytes) == 0
}
