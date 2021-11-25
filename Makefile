golden-update:
	GOLDEN_UPDATE=1 go test .

generate:
	go generate

test:
	go test .

.PHONY: golden-update generate test
