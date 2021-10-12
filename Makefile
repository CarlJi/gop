all:
	cd cmd; go install -v ./...

release:
	cd cmd; go install -v ./...
	ls -ltr ${GOPATH}/bin
	which gop

