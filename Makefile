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
VERSION := 0.0.4
PATH := $(PWD)/target/debug:${PATH}

# Project overrides, starting with prefix install
DESTDIR =
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
SBINDIR = $(PREFIX)/sbin
SYSCONFDIR = $(PREFIX)/etc
LOCALSTATEDIR = $(PREFIX)/var/
LOGPREFIXDIR = $(LOCaLSTATEDIR)/log
PREFIXPATH = $(LOCALSTATEDIR)/lib/babylon
TESTDIR = $(PWD)/test
TAGS =

.PHONY: all required version build release install clean

all:            build           # default target debug
required:       vendor          # required to build

# Define or override custom env
sinclude custom.mk

build:  required
	@mkdir -p target/debug
	@go build -v -tags debug,$(TAGS) -ldflags '-X main.version=$(VERSION) -X main.etcPrefix=$(TESTDIR) -X main.prefixPath=$(TESTDIR) -X main.logPrefix=$(TESTDIR)' -mod vendor -o target/debug ./...

release-static:	required
	@mkdir -p target/release
	@CGO_ENABLED=0 go build -v -mod vendor -tags release,static,$(TAGS) -ldflags '-s -w -extldflags -static -X main.version=$(VERSION) -X main.etcPrefix=$(SYSCONFDIR) -X main.prefixPath=$(PREFIXPATH) -X main.logPrefix=$(LOGPREFIXDIR)' -o target/release ./...

release-shared:	required
	@mkdir -p target/release
	@go build -v -mod vendor -tags release,$(TAGS) -ldflags '-s -w -X main.version=$(VERSION) -X main.etcPrefix=$(SYSCONFDIR) -X main.prefixPath=$(PREFIXPATH) -X main.logPrefix=$(LOGPREFIXDIR)' -o target/release ./...

install:	release
	@install -d -m 755 $(DESTDIR)$(BINDIR)
	@install -d -m 755 $(DESTDIR)$(SBINDIR)
	@install -d -m 755 $(DESTDIR)$(SYSCONFDIR)
	@install -d -m 755 $(DESTDIR)$(LOGPREFIXDIR)
	@install -d -m 755 $(DESTDIR)$(LOCALSTATEDIR)
	@install -D -m 644 etc/babylon.conf $(DESTDIR)$(SYSCONFDIR)/babylon.conf
	@install -m 755 target/release/f9600 $(DESTDIR)$(SBINDIR)

clean:
	@go clean -cache ./...
	@rm -rf target *.out
	@rm -f $(PROJECT)-*.tar.gz $(PROJECT)-*.tar

version:
	@echo $(VERSION)

docs:
	@go doc -all babylon/lib

# Optional make components we add
sinclude .make/*.mk

