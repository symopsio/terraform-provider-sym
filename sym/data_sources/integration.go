package data_sources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func DataSourceIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIntegrationRead,
		Schema: integrationSchema(),
	}
}

func required(valueType schema.ValueType) *schema.Schema {
	return &schema.Schema{
		Type:     valueType,
		Required: true,
	}
}

func integrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type":     required(schema.TypeString),
		"name": 	required(schema.TypeString),
	}
}

func dataSourceIntegrationRead(data *schema.ResourceData, meta interface{}) error {
	//c := meta.(*client.ApiClient)
	log.Printf("DataSourceIntegrationRead id %v", data.Id())
	//repoName := d.Get("repository").(string)
	//branchName := d.Get("branch").(string)
	//branchRefName := "refs/heads/" + branchName
	//
	//log.Printf("[DEBUG] Reading GitHub branch reference %s/%s (%s)",
	//	orgName, repoName, branchRefName)
	//ref, resp, err := client.Git.GetRef(
	//	context.TODO(), orgName, repoName, branchRefName)
	//if err != nil {
	//	return fmt.Errorf("Error reading GitHub branch reference %s/%s (%s): %s",
	//		orgName, repoName, branchRefName, err)
	//}
	//
	//d.SetId(buildTwoPartID(repoName, branchName))
	//d.Set("etag", resp.Header.Get("ETag"))
	//d.Set("ref", *ref.Ref)
	//d.Set("sha", *ref.Object.SHA)

	return nil
}