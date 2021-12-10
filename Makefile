build:
	go build -tags "luaa lua54" -o hue

run:
	go run -tags "luaa lua54" ./... test.lua

install:
	ln -sf $(shell pwd)/hue /usr/local/bin/hue

uninstall:
	rm /usr/local/bin/hue
