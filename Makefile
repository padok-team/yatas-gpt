default: build

test:
	go test ./...

build:
	go build -o bin/yatas-template

update:
	go get -u 
	go mod tidy

install: build
	mkdir -p ~/.yatas.d/plugins/github.com/StanGirard/yatas-template/local/
	mv ./bin/yatas-template ~/.yatas.d/plugins/github.com/StanGirard/yatas-template/local/

release: test
	standard-version
	git push --follow-tags origin main 