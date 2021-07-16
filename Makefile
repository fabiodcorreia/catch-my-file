#VERSION = $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)
VERSION = $(shell git describe --tags --match=v* 2> /dev/null || echo 0.0.0)

APPID = com.github.fabiodcorreia.catch-my-file
ICON = assets/icons/icon-512.png
NAME = CatchMyFile

format:
	gofmt -s -w main.go 
	gofmt -s -w internal/**/*.go 
	gofmt -s -w cmd/**/*.go

review: format
	@echo "============= Spell Check ============= "
	@misspell .
	
	@echo "============= Ineffectual Assignments Check ============= "
	@ineffassign ./...

	@echo "============= Cyclomatic Complexity Check ============= "
	@gocyclo -total -over 5 -avg .

	@echo "============= Duplication Check ============= "
	@dupl -t 25

	@echo "============= Repeated Strings Check ============= "
	@goconst ./...

	@echo "============= Vet Check ============= "
	@go vet ./...
	
build:
	go mod tidy
	go build -tags release -ldflags="-s -w"  -o $(NAME)

darwin:
	fyne-cross darwin -arch amd64,arm64 -app-id $(APPID) -icon $(ICON) -app-version $(VERSION) -output $(NAME)
	
linux:
	fyne-cross linux -arch amd64,arm64 -app-id $(APPID) -icon $(ICON) -app-version $(VERSION)

windows:
	fyne-cross windows -arch amd64 -app-id $(APPID) -icon $(ICON) -app-version $(VERSION)

bundle:
	rm -fr dist
	mkdir dist

	mv fyne-cross/dist/linux-amd64/$(NAME).tar.gz $(NAME)-$(VERSION)-linux-amd64.tar.gz
	mv fyne-cross/dist/linux-arm64/$(NAME).tar.gz $(NAME)-$(VERSION)-linux-arm64.tar.gz

	(cd fyne-cross/dist/darwin-amd64/ && zip -r $(NAME)-darwin-amd64.zip $(NAME).app/)
	mv fyne-cross/dist/darwin-amd64/$(NAME)-darwin-amd64.zip dist/$(NAME)-$(VERSION)-darwin-amd64.zip

	(cd fyne-cross/dist/darwin-arm64/ && zip -r $(NAME)-darwin-arm64.zip $(NAME).app/)
	mv fyne-cross/dist/darwin-arm64/$(NAME)-darwin-arm64.zip dist/$(NAME)-$(VERSION)-darwin-arm64.zip

	mv fyne-cross/dist/windows-amd64/$(NAME).exe.zip $(NAME)-$(VERSION)-windows-amd64.zip

release: darwin freebsd linux windows bundle

tools:
	go get -u github.com/jgautheron/goconst/cmd/goconst
	go get -u github.com/mdempsky/unconvert
	go get -u github.com/securego/gosec/v2/cmd/gosec
	go get -u github.com/alexkohler/prealloc