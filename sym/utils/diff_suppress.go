package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SuppressEquivalentJsonDiffs is a DiffSuppressFunc that can be passed into
// a schema to ensure meaningless differences in JSON are ignored when reading/writing.
//
// e.g. without this function, TF will see the following two blobs as different, and
// prompt an update. With this function set, the blobs are seen as equivalent:
//
//  {"template": "sym:approval", "name": "access"}
//  {"name": "access", "template": "sym:approval"}
func SuppressEquivalentJsonDiffs(k string, old string, new string, _ *schema.ResourceData) bool {
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

func SuppressFlowDiffs(k string, old string, new string, d *schema.ResourceData) bool {
	suppressJsonDiffs := SuppressEquivalentJsonDiffs(k, old, new, d)

	// allow_revoke defaults to true, so don't show a diff if allow_revoke is not specified
	suppressAllowRevokeDiffs := k == "params.allow_revoke" && old == "true" && new == ""

	// schedule_deescalation defaults to true, so don't show a diff if schedule_deescalation is not specified
	suppressScheduleDeescalationDiffs := k == "params.schedule_deescalation" && old == "true" && new == ""

	return suppressJsonDiffs || suppressAllowRevokeDiffs || suppressScheduleDeescalationDiffs
}

// SuppressNullSettingsDiffs is a DiffSuppressFunc that can be passed into
// a schema to account for differences in settings maps received from the API
// and stored in the Terraform state.
//
// e.g. if a Terraform object's settings are empty:
// 	settings {}
//
// the API may still send the full dictionary of possible settings will null values
// 	{"account_id": null, ...}
//
// where effectively those values are the same, but Terraform's diff will not
// know that.
//
// SuppressNullSettingsDiffs will be called on each item in the settings
// map individually as well as on the settings map as a whole.
func SuppressNullSettingsDiffs(k string, old string, new string, d *schema.ResourceData) bool {
	isTopLevel := strings.HasSuffix(k, ".%")

	if isTopLevel {
		// If we're at the top level of the settings map, all incoming values are null,
		// and there are no values locally, suppress that difference.
		// This is because if all values are null on both sides, suppressing each individual
		// item is not sufficient to suppress the update.
		settingsName := strings.TrimSuffix(k, ".%")
		settings := d.Get(settingsName).(map[string]interface{})
		allSettingsNull := true

		for _, v := range settings {
			if v.(string) != "" {
				allSettingsNull = false
				break
			}
		}

		i, err := strconv.Atoi(new)
		if err != nil {
			log.Printf("Error getting new settings value count from %v for value %v", new, k)
			return false
		}

		if i == 0 && allSettingsNull {
			return true
		}
	} else if !isTopLevel && old == "" && new == "" {
		// If we're not at the top level of the settings map, and both the existing
		// and the incoming data values are all null, suppress that difference.
		// This is sufficient if not all values in the settings map are null on both
		// sides.
		return true
	}

	return false
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
func SuppressEquivalentFileContentDiffs(k string, old string, new string, _ *schema.ResourceData) bool {
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

	return bytes.Equal(newBytes, oldBytes)
}
