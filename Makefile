releases:
	mkdir -p dist
	env GOOS=darwin go build && tar -czf dist/loopit-osx.tar.gz loopit
	env GOOS=linux go build && tar -czf dist/loopit-linux.tar.gz loopit
	rm loopit
