test: fmtcheck
	go test ./... -v -timeout 120s

testacc: fmtcheck
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120s

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
