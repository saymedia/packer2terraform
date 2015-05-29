
# Testing
go test ./...

# Linting
go vet */*.go
go vet *.go
golint */*.go
golint *.go
gocyclo -over 10 */*.go
gocyclo -over 10 *.go
