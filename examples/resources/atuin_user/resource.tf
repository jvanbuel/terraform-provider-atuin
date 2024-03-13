resource "atuin_user" "test" {
  username = "twoflower"
  email    = "twoflower@discworld.co.uk"
  password = "swordfish"
}

output "encryption_key" {
  sensitive = true
  value     = atuin_user.test.base64_key
}
