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
			"template": {
				Type: schema.TypeString,
				Required: true,
			},
			"handler": {
				Type: schema.TypeString,
				Required: true,
			},
			"strategy_param": {
				Type: schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"group_label": {
							Type:     schema.TypeString,
							Required: true,
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
	version := d.Get("version").(int)
	qualifiedName := qualifyName(c.GetOrg(), name)
	log.Printf("[DEBUG] Qualified name: %s", qualifiedName)

	handler := d.Get("handler").(string)
	body, err := readUTF8(handler)
	if err != nil {
		return diag.FromErr(err)
	}
	template := d.Get("template").(string)

	flow := &models.Flow{
		Name:    qualifiedName,
		Version: &models.Version{
			Major: int32(version),
		},
		Uuid: "bd6b69bd-0d93-463e-b997-b19a8370da6e",
		Template: &models.Template{
			Name: template,
			Version: &models.Version{
				Major: 1,
			},
		},
		Implementation: &models.Source{
			Body:     body,
			Filename: handler,
		},
	}

	id, err := c.CreateFlow(flow)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Created flow with id: %s", id)
	d.SetId(id)

	resourceFlowRead(ctx, d, m)

	return diags
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	var diags diag.Diagnostics

	flow, err := c.GetFlow(d.Id())
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
	// TODO: implement
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: implement
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
