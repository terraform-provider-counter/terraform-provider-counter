terraform {
  required_providers {
    counter = {
      source = "terraform-provider-counter/counter"
    }
  }
}

resource "counter_monotonic" "this" {
  step          = 1
  initial_value = 0
  triggers = {
    hash = md5(jsonencode(something_else.this))
  }
}

resource "downstream" "this" {
  value = counter_monotonic.this.value
}
