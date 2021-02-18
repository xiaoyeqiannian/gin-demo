## go build
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
```
## go http status
https://github.com/golang/go/blob/c04e47f82159f1010d7403276b3dff5ab836fd00/src/net/http/status.go

## api swagger
```
cd cmd/account
swag init -d ../../app/account/api_handler.go -o ../../app/account/docs -g ./
```