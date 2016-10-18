default:
	$(MAKE) all
test:
	bash -c "./scripts/test.sh"
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"