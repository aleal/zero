test:
	go test ./... -coverprofile=coverage.txt -covermode=atomic -v

coverage:
	make test
	go tool cover -html=coverage.txt

bench:
	go test -bench=. ./... | tee bench.new.txt

bench-reg:
	go test -bench=. ./... | tee bench.reg.txt

bench-cmp:
	~/go/bin/benchcmp bench.new.txt bench.reg.txt

bench-stat:
	~/go/bin/benchstat bench.new.txt bench.reg.txt

build-example-server:
	go build -o build/example-server examples/server/cmd/main.go

run-example-server:
	go run examples/server/cmd/main.go