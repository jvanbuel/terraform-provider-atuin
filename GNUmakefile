default: testacc

# Run acceptance tests

ATUIN_PORT := 8888
.PHONY: testacc docs
testacc:
	TF_ACC=1 ATUIN_HOST=http://localhost:$(ATUIN_PORT) go test ./... -v $(TESTARGS) -timeout 120m

docs:
	tfplugindocs generate --provider-name terraform-provider-atuin

install:
	go install .
