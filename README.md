A terraform provider for generating sequential values based on generic triggers.

---

## Provider

This provider is available from the [Terraform Registry](https://registry.terraform.io/providers/terraform-provider-counter/counter/latest/docs).

```terraform 
terraform {
    required_providers = {
        counter = {
            source = "terraform-provider-counter/counter"
        }
    }
}
```

---

## Resources

Examples of supported resources are provided below.

- [Monotonic](#monotonic)
- [Semantic Version](#semantic-version)

---

#### Monotonic

Use this to produce a number which increments by step each time there's a change to any triggers.

```terraform
resource counter_monotonic this {
    step = 1
    initial_value = 0
    triggers = {
        hash = md5(jsonencode(something_else.this))
    }
}

resource downstream this {
    value = counter_monotonic.this.value
}
```

---

#### Semantic Version

Use this to produce a semantic version which increments each time there's a change to any triggers of the relevant
version component. When the major version changes, the minor and patch versions start over at zero. When the minor
changes, the patch version starts over at zero.

```terraform
resource counter_semantic_version this {
    minor_triggers = {
        hash = md5(jsonencode(something_else.this))
    }
    patch_triggers = {
        hash = md5(jsonencode(something_else.that))
    }
}

resource downstream this {
    value = counter_semantic_version.this.value
}
```

---

## License

This project is licensed under [MIT license](http://opensource.org/licenses/MIT).
