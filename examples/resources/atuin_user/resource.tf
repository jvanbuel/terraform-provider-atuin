terraform {
  required_providers {
    atuin = {
      source = "lightcone/atuin"
    }
  }
}

provider "atuin" {
  # Configuration options
}


resource "atuin_user" "test" {
  username = "humfrey123"
  email    = "testing12345@yahoo.com"
  password = "password123"
}

output "encryption_key" {
  sensitive = true
  value     = atuin_user.test.key
}
