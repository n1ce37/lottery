.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/lotteries lotteries/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/detail detail/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/histories histories/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
