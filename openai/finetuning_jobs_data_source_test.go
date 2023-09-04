package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuningJobsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFineTunesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify ID has any value set
					resource.TestCheckResourceAttrSet("data.openai_finetuning_jobs.test", "id"),
				),
			},
		},
	})
}

const testAccFineTunesDataSourceConfig = `
data "openai_finetuning_jobs" "test" {}
`
