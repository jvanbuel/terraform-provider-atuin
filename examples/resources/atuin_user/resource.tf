resource "atuin_user" "test" {
  username = "humfrey123"
  email    = "testing12345@yahoo.com"
  password = "password123"
}

output "encryption_key" {
  sensitive = true
  value     = atuin_user.test.base64_key
}
