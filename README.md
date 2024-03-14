<img src="image/README/1710340059105.png" alt="drawing" width="500"/>

**⚠️ This repository is not associated with the Atuin project, and should be considered highly experimental. Use at your own risk. ⚠️**

This repository contains a Terraform provider for Atuin. It allows you to create and manage Atuin users.

The original motivation behind this provider was to manage the lifecycle of Atuin users for online IDE environments, like e.g. Gitpod, as part of your IaC codebase. In addition to managing the lifecycle Atuin users for individuals, which is mostly feasible for smaller organizations, you could also create project- and/or repository-specific users. By having a single shell history per project/repository, you have access to a richer history of shell commands relevant to that project/repository.

The provider creates the key that is used for encrypting the shell history automatically, and provides it as a base64 and bip39 formatted attribute of the resource. You can use it then to configure any other environments in which you want to use Atuin to sync shell history. As with any other sensitive information that is managed by Terraform: make sure to use a remote state!

Example usage:

```hcl
terraform {
  required_providers {
    atuin = {
      source = "lightcone/atuin"
    }
  }
}

provider "atuin" {
}


resource "atuin_user" "user" {
  username = "twoflower"
  password = "swordfish"
  email    = "twoflower@discworld.co.uk"
}

output "bip39" {
  sensitive = true
  value     = atuin_user.user.bip39_key
}
```

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-atuin
```

## Test provider

To test the provider, start an Atuin server via docker-compose:

```shell
$ cd docker-compose && docker-compose up
```

Then run the acceptance tests:

```shell
$ make testacc
```
