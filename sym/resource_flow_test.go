package sym

import (
	"testing"

	"github.com/symopsio/protos/go/tf/models"
)

func TestFormatID(t *testing.T) {
	got := formatID("foo", 1)
	if got != "foo:1" {
		t.Errorf("Expected foo:1, got: %s", got)
	}
}

func TestParseNameAndVersion(t *testing.T) {
	name, version, err := parseNameAndVersion("org:foo:1")
	if err != nil {
		t.Errorf("Expected to parse foo:1, got err: %v", err)
	}
	if name != "org:foo" {
		t.Errorf("Expected foo, got: %s", name)
	}
	if version != 1 {
		t.Errorf("Expected 1, got: %v", 1)
	}
}

func TestFlattenHandler(t *testing.T) {
	flow := &models.Flow{
		Implementation: &models.Source{
			Body: "foo",
		},
		Template: &models.Template{
			Name: "bar",
		},
	}

	flowData := flattenHandler(flow)[0].(map[string]interface{})
	if body := flowData["body"]; body != "foo" {
		t.Errorf("Expected foo, got: %s", body)
	}
	if template := flowData["template"]; template != "bar" {
		t.Errorf("Expected bar, got: %s", template)
	}
}
