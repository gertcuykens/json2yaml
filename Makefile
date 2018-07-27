build: version=$(shell date "+(%T %d %b %Y)")
build:
	go build -ldflags "-X \"main.versionInfo=$(version)\""
