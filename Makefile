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

## dist: Compiles UNIX and Raspberry Pi compatible binaries
dist: dist/hive dist/unix dist/rpi 

## dist/hive: Compiles the Hive GUI
dist/hive: service.hive/
	@echo "	>	Building Hive GUI assets"
	@mkdir -p $(dir $@)
	@cd service.hive && npm run build
	@mv service.hive/dist dist/hive
	@cp service.hive/nginx.conf dist/hive/

## dist/unix: Compiles UNIX compatible binaries
dist/unix: dist/unix/registry dist/unix/shelly dist/unix/hue

## dist/rpi: Compiles Raspberry Pi compatible binaries
dist/rpi: dist/rpi/registry dist/rpi/shelly dist/rpi/hue

## dist/unix/%: Compiles an UNIX compatible binary
dist/unix/%: service.%/
	@echo "	>	Building UNIX binary: $(notdir $@)"
	@mkdir -p $(dir $@)
	@$(GOBUILD) -o $(GOBASE)/$@ $(GOBASE)/$<
	@cp $(GOBASE)/$</README.md $(GOBASE)/$@.md

## dist/rpi/%: Compiles a Raspberry Pi compatible binary
dist/rpi/%: service.%/
	@echo "	>	Building static Raspberry Pi binary: $(notdir $@)"
	@mkdir -p $(dir $@)
	@CGO_ENABLED=0 $(GOBUILD) -ldflags "-extldflags -static" -o $(GOBASE)/$@ $(GOBASE)/$<
	@cp $(GOBASE)/$</README.md $(GOBASE)/$@.md

## help: Display this message
help: Makefile
	@echo " Choose a command run:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
