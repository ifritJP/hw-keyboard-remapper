help:
	@echo make all
	@echo make all-win
	@echo make all-raspi

PLATFORM=$(shell uname -m)

all:
ifeq ($(PLATFORM),x86_64)
	go build -o convkey.x86
else
	go build
endif

all-raspi-test:
	CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -v -o convkey.raspi
