GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o
BIN=tests/better-av-tool
SOURCE=.

docker:
	$(GOBUILD) $(BIN) $(SOURCE)