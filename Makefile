unit-test:
	go test -v -tags=unit -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

int-test:
	go test -v -tags=integration -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

test:
	go test -v -tags="integration unit" -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

mod-up:
	go get github.com/jschneider98/jgocache
	go get github.com/jschneider98/jgovalidator
