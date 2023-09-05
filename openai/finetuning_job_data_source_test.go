package openai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFineTuneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccOpenAI(t); testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFineTuneDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.openai_finetuning_job.test", "id"),
				),
			},
		},
	})
}

const testAccFineTuneDataSourceConfig = `
data "openai_finetuning_jobs" "test" {
}

data "openai_finetuning_job" "test" {
	id = data.openai_finetuning_jobs.test.jobs[0].id
}
`
