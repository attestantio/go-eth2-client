fmt:
	gofmt -s -w .
	gofumpt -extra -w .
	gci write .
	go mod tidy