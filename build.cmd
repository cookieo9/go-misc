set GOPATH=%BuildFolder%/..
go env
go get -d -v ./...
go test -v ./...
