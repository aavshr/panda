build:
	go build -o panda --tags "fts5"

run: build
	-./panda
	rm ./panda
