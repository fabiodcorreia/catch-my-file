#VERSION = $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)
VERSION = $(shell git describe --tags --match=v* 2> /dev/null || echo 0.1.0)

APPID = com.github.fabiodcorreia.catch-my-file
ICON = assets/icons/icon-512.png
NAME = CatchMyFile

TARGET = pkg/**/*.go

format:
	@gofmt -s -w $(TARGET)

review: format
	@echo "============= Spell Check ============= "
	@misspell .
	
	@echo "============= Ineffectual Assignments Check ============= "
	@ineffassign ./...

	@echo "============= Duplication Check ============= "
	find ./pkg -not -name '*_test.go' -name '*.go' | dupl -t 30 -files

	@echo "============= Repeated Strings Check ============= "
	@goconst $(TARGET)

	@echo "============= Security Check ============= "
	@gosec ./...

	@echo "============= Vet Check ============= "
	@go vet --all .

	@echo "============= Preallocation Check ============= "
	@prealloc -forloops -set_exit_status -simple -rangeloops ./...

	@echo "============= Shadow Variables Check ============= "
	@shadow -strict ./...

	@echo "============= Cyclomatic Complexity Check ============= "
	@gocyclo -total -ignore "_test" -over 8 -avg $(TARGET)
	
test: 
	@go test -cover ./pkg/...

cover: 
	@go test -coverprofile=coverage.out ./pkg/...
	@go tool cover -func=coverage.out

cover-html: 
	@go test -coverprofile=coverage.out ./pkg/...
	@go tool cover -html=coverage.out


bench:
	go test -benchtime=1s -count=5 -benchmem -bench . ./pkg/...

pre-build: review
	go mod tidy

build: pre-build
	go build -tags release -ldflags="-s -w"  -o $(NAME)

darwin: pre-build
	fyne-cross darwin -arch amd64,arm64 -app-id $(APPID) -icon $(ICON) -app-version $(VERSION) -output $(NAME)
	
linux: pre-build
	fyne-cross linux -arch amd64,arm64 -app-id $(APPID) -icon $(ICON) -app-version $(VERSION) -output $(NAME)

windows: pre-build
	fyne-cross windows -arch amd64 -app-id $(APPID) -icon $(ICON) -app-version $(VERSION) -output $(NAME)

bundle-linux: linux
	mv fyne-cross/dist/linux-amd64/$(NAME).tar.gz dist/$(NAME)-$(VERSION)-linux-amd64.tar.gz
	mv fyne-cross/dist/linux-arm64/$(NAME).tar.gz dist/$(NAME)-$(VERSION)-linux-arm64.tar.gz

bundle-darwin: darwin
	(cd fyne-cross/dist/darwin-amd64/ && zip -r $(NAME)-darwin-amd64.zip $(NAME).app/)
	mv fyne-cross/dist/darwin-amd64/$(NAME)-darwin-amd64.zip dist/$(NAME)-$(VERSION)-darwin-amd64.zip

	#(cd fyne-cross/dist/darwin-arm64/ && zip -r $(NAME)-darwin-arm64.zip $(NAME).app/)
	#mv fyne-cross/dist/darwin-arm64/$(NAME)-darwin-arm64.zip dist/$(NAME)-$(VERSION)-darwin-arm64.zip

bundle-windows: windows
	mv fyne-cross/dist/windows-amd64/$(NAME).zip dist/$(NAME)-$(VERSION)-windows-amd64.zip

release: bundle-linux bundle-darwin bundle-windows

tools:
	go get -u github.com/jgautheron/goconst/cmd/goconst
	go get -u github.com/mdempsky/unconvert
	go get -u github.com/securego/gosec/v2/cmd/gosec
	go get -u github.com/alexkohler/prealloc
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest