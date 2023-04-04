default: build

test:
	go test -v ./...

build:
	go build -o bin/yatas-gpt

update:
	go get -u 
	go mod tidy

install: build
	mkdir -p ~/.yatas.d/plugins/github.com/padok-team/yatas-gpt/local/
	mv ./bin/yatas-gpt ~/.yatas.d/plugins/github.com/padok-team/yatas-gpt/local/

release: test
	standard-version
	git push --follow-tags origin main 