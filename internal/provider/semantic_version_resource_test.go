package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func step1Semantic() string {
	return `
		resource counter_semantic_version this {
			patch_triggers = {
				hash = "potatoes"
			}
		}
	`
}

func step2Semantic() string {
	return `
		resource counter_semantic_version this {
			patch_triggers = {
				hash = "eggs"
			}
		}
	`
}

func step3Semantic() string {
	return `
		resource counter_semantic_version this {
			minor_triggers = {
				hash = "potatoes"
			}
			patch_triggers = {
				hash = "bacon"
			}
		}
	`
}
func step4Semantic() string {
	return `
		resource counter_semantic_version this {
			major_triggers = {
				hash = "potatoes"
			}
			minor_triggers = {
				hash = "eggs"
			}
			patch_triggers = {
				hash = "butter"
			}
		}
	`
}

func TestAccSemanticVersionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test initial release
			{
				Config: step1Semantic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_semantic_version.this", "value", "1.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.value", "1.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.major_value", "1"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.minor_value", "0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.patch_value", "0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.patch_triggers.hash", "potatoes"),
				),
			},
			// Test patch release
			{
				Config: step2Semantic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_semantic_version.this", "value", "1.0.1"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.value", "1.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.patch_triggers.hash", "potatoes"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.1.value", "1.0.1"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.1.patch_triggers.hash", "eggs"),
				),
			},
			// Test minor release
			{
				Config: step3Semantic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_semantic_version.this", "value", "1.1.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.value", "1.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.patch_triggers.hash", "potatoes"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.1.value", "1.0.1"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.1.patch_triggers.hash", "eggs"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.2.value", "1.1.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.2.patch_triggers.hash", "bacon"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.2.minor_triggers.hash", "potatoes"),
				),
			},
			// Test major release
			{
				Config: step4Semantic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("counter_semantic_version.this", "value", "2.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.value", "1.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.0.patch_triggers.hash", "potatoes"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.1.value", "1.0.1"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.1.patch_triggers.hash", "eggs"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.2.value", "1.1.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.2.patch_triggers.hash", "bacon"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.2.minor_triggers.hash", "potatoes"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.3.value", "2.0.0"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.3.patch_triggers.hash", "butter"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.3.minor_triggers.hash", "eggs"),
					resource.TestCheckResourceAttr("counter_semantic_version.this", "history.3.major_triggers.hash", "potatoes"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
