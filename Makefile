local-build:
	go build -o panda

local-run: build
	-./panda
	rm ./panda