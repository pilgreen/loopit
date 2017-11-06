releases:
	go build && tar -czf ~/Downloads/loopit-osx.tar.gz loopit
	env GOOS=linux go build && tar -czf ~/Downloads/loopit-linux.tar.gz loopit
	rm loopit
