#!/usr/bin/make -f
# Copyright (C) 2020-2021 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by a MIT-style license
# that can be found in the included LICENSE.md file.

# This is needed because GO lacks proper project level build support.  Sure,
# they could have used go.mod to do some things such as store project metadata
# like Cargo.toml does for rust, but they didn't. So we have a make to get past
# stupid go project limitations.

# Project constants
PROJECT := babylon
VERSION := 0.0.1
PATH := $(PWD)/target/debug:${PATH}

# Project overrides, starting with prefix install
DESTDIR =
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
SYSCONFDIR = $(PREFIX)/etc
LOGPREFIXDIR = $(PREFIX)/var/log
SHAREDSTATEDIR = $(PREFIX)/var/lib/babylon
TESTDIR = $(PWD)/test
TAGS =

.PHONY: all required version build release install clean

all:            build           # default target debug
required:       vendor          # required to build

# Define or override custom env
sinclude custom.mk

build:  required
	@mkdir -p target/debug
	@go build -v -tags debug,$(TAGS) -ldflags '-X main.version=$(VERSION) -X main.etcPrefix=$(TESTDIR) -X main.varPrefix=$(TESTDIR) -X main.logPrefix=$(TESTDIR) -X main.buildType=debug' -mod vendor -o target/debug ./...

release:        required
	@mkdir -p target/release
	@CGO_ENABLED=0 go build -v -mod vendor -tags release,static,$(TAGS) -ldflags '-s -extldflags -static -X main.version=$(VERSION) -X main.etcPrefix=$(SYSCONFDIR) -X main.varPrefix=$(SHAREDSTATEDIR) -X main.logPrefix=$(LOGPREFIXDIR)' -o target/release ./...

install:        release
	@install -d -m 755 $(DESTDIR)$(BINDIR)
	@install -d -m 755 $(DESTDIR)$(SBINDIR)
	@install -d -m 755 $(DESTDIR)$(SYSCONFDIR)
	@install -d -m 755 $(DESTDIR)$(LOGPREFIXDIR)
	@install -d -m 755 $(DESTDIR)$(SHAREDSTATEDIR)
	@install -D -m 644 test/babylon.conf $(DESTDIR)$(SHAREDSTATEDIR)/babylon.conf
	@install -m 755 target/release/f9600-server $(DESTDIR)$(SBINDIR)
	@install -m 755 target/release/f9600-mml $(DESTDIR)$(BINDIR)

clean:
	@go clean -cache ./...
	@rm -rf target *.out vendor
	@rm -f $(PROJECT)-*.tar.gz $(PROJECT)-*.tar

version:
	@echo $(VERSION)

docs:
	@go doc -all switchroom/lib

# Optional make components we add
sinclude .make/*.mk

