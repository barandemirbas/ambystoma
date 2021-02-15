format:
	gofmt -l -s -w .

run:
	./bin/ambystoma

build:
	go mod tidy && go mod vendor
	go build -o ./bin/ambystoma && chmod +X ./bin/*
	echo "Ambystoma ready to use in ./bin/ambystoma"

cross-compile:
	env GOOS=linux GOARCH=arm go build -o ./release/ambystoma-linux-arm32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=arm64 go build -o ./release/ambystoma-linux-arm64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=darwin GOARCH=amd64 go build -o ./release/ambystoma-mac-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=386 go build -o ./release/ambystoma-linux-x32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=amd64 go build -o ./release/ambystoma-linux-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=windows GOARCH=386 go build -o ./release/ambystoma-windows-x32.exe -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=windows GOARCH=amd64 go build -o ./release/ambystoma-windows-x64.exe -ldflags "-s -w" -trimpath -mod=readonly

m1:
	env GOOS=darwin GOARCH=arm64 go build -o ./release/ambystoma-mac-arm64 -ldflags "-s -w" -trimpath -mod=readonly

clean:
	rm -rf ./bin/*
	rm -rf ./release/*
	rm -rf ./vendor/*
