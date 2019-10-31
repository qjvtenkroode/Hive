# Parameters
GOCMD=go
GOBASE=$(shell pwd)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

all:

## clean: Removes all compiled data
clean:
	@echo "	>	Cleaning dist"
	@rm -rf $(GOBASE)/dist

## dist: Compiles OSX and Raspberry Pi compatible binaries
dist: dist/osx dist/rpi

## dist/osx: Compiles OSX compatible binaries
dist/osx: dist/osx/registry dist/osx/shelly

## dist/rpi: Compiles Raspberry Pi compatible binaries
dist/rpi: dist/rpi/registry dist/rpi/shelly

## dist/osx/%: Compiles an OSX compatible binary
dist/osx/%: service.%/
	@echo "	>	Building OSX binary: $(notdir $@)"
	@mkdir -p $(dir $@)
	@$(GOBUILD) -o $(GOBASE)/$@ $(GOBASE)/$<

## dist/rpi/%: Compiles a Raspberry Pi compatible binary
dist/rpi/%: service.%/
	@echo "	>	Building static Raspberry Pi: binary $(notdir $@)"
	@mkdir -p $(dir $@)
	@CGO_ENABLED=0 $(GOBUILD) -ldflags "-extldflags -static" -o $(GOBASE)/$@ $(GOBASE)/$<

## help: Display this message
help: Makefile
	@echo " Choose a command run:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
