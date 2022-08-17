golden-update:
	GOLDEN_UPDATE=1 go test ./...

generate:
	go generate ./...

test:
	go test ./...

cov-html:
	go test -coverprofile=cov.prof ./...
	go tool cover -html=cov.prof -o cov.html
	rm cov.prof

.PHONY: golden-update generate test cov-html
