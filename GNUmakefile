default: testacc

# Run acceptance tests
.PHONY: testacc docs
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

docs:
	tfplugindocs generate --provider-name terraform-provider-atuin

install:
	go install .
