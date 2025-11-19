build:
	go build -o panda --tags "fts5"

test:
	go test -v -tags "fts5" ./...

run: build
	-./panda
	rm ./panda
