all:
	cd cmd; go install -v ./...

release:
	cd cmd; go install -v ./...
	ls -ltr ${GOPATH}/bin
	ls -ltr ${HOME}/go/bin
	which gop

