package sym

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symopsio/protos/go/tf/models"
	"github.com/symopsio/terraform-provider-sym/sym/client"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowCreate,
		ReadContext:   resourceFlowRead,
		UpdateContext: resourceFlowUpdate,
		DeleteContext: resourceFlowDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// You currently can't represent nested structures (like handler) without
			// wrapping in a single-element list:
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/155
			"handler": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	qualifiedName := qualifyName(c.GetOrg(), name)

	version := uint32(d.Get("version").(int))

	handlers := d.Get("handler").([]interface{})
	handler := handlers[0].(map[string]interface{})

	source := handler["source"].(string)
	body, err := readUTF8(source)
	if err != nil {
		return diag.FromErr(err)
	}
	handler["body"] = body

	template := handler["template"].(string)

	flow := &models.Flow{
		Name:    qualifiedName,
		Version: version,
		Template: &models.Template{
			Name: template,
		},
		Implementation: &models.Source{
			Body:     body,
			Filename: source,
		},
	}

	err = c.CreateFlow(flow)
	if err != nil {
		return diag.FromErr(err)
	}

	id := formatID(qualifiedName, version)

	log.Printf("[DEBUG] Created flow with id: %s", id)

	d.SetId(id)

	resourceFlowRead(ctx, d, m)

	return diags
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	var diags diag.Diagnostics

	name, version, err := parseNameAndVersion(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flow, err := c.GetFlow(name, version)
	if err != nil {
		return diag.FromErr(err)
	}

	handlerData := flattenHandler(flow)
	if err := d.Set("handler", handlerData); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func flattenHandler(flow *models.Flow) []interface{} {
	h := make(map[string]interface{})
	if flow.Implementation != nil {
		h["source"] = flow.Implementation.Filename
		h["body"] = flow.Implementation.Body
	}
	if flow.Template != nil {
		h["template"] = flow.Template.Name
	}

	return []interface{}{h}
}

func qualifyName(org string, name string) string {
	return fmt.Sprintf("%s:%s", org, name)
}

func parseVersion(s string) (uint32, error) {
	version, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint32(version), nil
}

func formatID(name string, version uint32) string {
	return fmt.Sprintf("%s:%v", name, version)
}

func parseNameAndVersion(id string) (string, uint32, error) {
	split := strings.SplitN(id, ":", 3)
	if len(split) < 3 {
		return "", 0, fmt.Errorf("Unsupported id: %s", id)
	}
	version, err := parseVersion(split[2])
	if err != nil {
		return "", 0, err
	}
	return qualifyName(split[0], split[1]), version, nil
}

func readUTF8(path string) (string, error) {
	src, err := readFileBytes(".", path)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(src) {
		return "", fmt.Errorf("contents of %s are not valid UTF-8; use the filebase64 function to obtain the Base64 encoded contents or the other file functions (e.g. filemd5, filesha256) to obtain file hashing results instead", path)
	}
	return string(src), nil
}

// Reusing the file function's implementation from here
// https://github.com/hashicorp/terraform/blob/master/lang/funcs/filesystem.go#L355
func readFileBytes(baseDir, path string) ([]byte, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand ~: %s", err)
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}

	// Ensure that the path is canonical for the host OS
	path = filepath.Clean(path)

	src, err := ioutil.ReadFile(path)
	if err != nil {
		// ReadFile does not return Terraform-user-friendly error
		// messages, so we'll provide our own.
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no file exists at %s; this function works only with files that are distributed as part of the configuration source code, so if this file will be created by a resource in this configuration you must instead obtain this result from an attribute of that resource", path)
		}
		return nil, fmt.Errorf("failed to read %s", path)
	}

	return src, nil
}
