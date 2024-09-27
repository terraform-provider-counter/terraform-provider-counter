package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func step1() string {
	return `
		resource counter_monotonic this {
			initial_value = 35
			triggers = {
				hash = "potatoes"
			}
		}
	`
}

func step2() string {
	return `
		resource counter_monotonic this {
			initial_value = 35
			triggers = {
				hash = "eggs"
			}
		}
	`
}
func step3() string {
	return `
		resource counter_monotonic this {
			initial_value = 35
			step = 2
			triggers = {
				hash = "bacon"
			}
		}
	`
}

func TestAccMonotonicResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: step1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_monotonic.this", "value", "35"),
					resource.TestCheckResourceAttr("counter_monotonic.this", "history.0.value", "35"),
				),
			},
			// Update and Read testing
			{
				Config: step2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_monotonic.this", "value", "36"),
					resource.TestCheckResourceAttr("counter_monotonic.this", "history.0.value", "35"),
					resource.TestCheckResourceAttr("counter_monotonic.this", "history.1.value", "36"),
				),
			},
			// Test change of step value
			{
				Config: step3(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_monotonic.this", "value", "38"),
					resource.TestCheckResourceAttr("counter_monotonic.this", "history.0.value", "35"),
					resource.TestCheckResourceAttr("counter_monotonic.this", "history.1.value", "36"),
					resource.TestCheckResourceAttr("counter_monotonic.this", "history.2.value", "38"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
