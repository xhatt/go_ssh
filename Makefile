
release:
	echo "发版"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/go_ssh-mac-x86.Intel
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/go_ssh-mac-arm64.m1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/go_ssh-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/go_ssh-linux-arm64



