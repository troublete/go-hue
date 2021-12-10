build:
	go build -tags "luaa lua54" -o hue

run:
	go run -tags "luaa lua54" ./... test.lua