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
  username = "test123456789"
  email    = "weifjweijw@gmail.com"
  password = "password1234"
}
