terraform {
  required_providers {
    counter = {
      source = "terraform-provider-counter/counter"
    }
  }
}

resource "counter_semantic_version" "this" {
  minor_triggers = {
    hash = md5(jsonencode(something_else.this))
  }
  patch_triggers = {
    hash = md5(jsonencode(something_else.that))
  }
}

resource "downstream" "this" {
  value = counter_semantic_version.this.value
}
