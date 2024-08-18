local-build:
	go build -o panda

local-run: local-build
	-./panda
	rm ./panda
