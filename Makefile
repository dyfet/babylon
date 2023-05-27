#!/usr/bin/make -f
# Copyright (C) 2020-2023 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by the terms of the GNU GPL
# v3 or later as found in the included LICENSE.md file.
#
# This is needed because GO lacks proper project level build support.
# Sure, they could have used go.mod to do some things such as store
# project metadata like Cargo.toml does for rust, but they didn't. So
# we have a make to get past stupid go project limitations.

# Project constants
PROJECT := babylon
VERSION := 0.0.6
PATH := $(PWD)/target/debug:${PATH}
TESTDIR := $(PWD)/test

.PHONY: all required version build debug release install clean

all:		build		# default target debug
required:       vendor          # required to build

# Define or override custom env
sinclude custom.mk

build:  lint
	@mkdir -p target/debug
	@go build -v -tags debug,$(TAGS) -ldflags '-X main.version=$(VERSION) -X main.etcPrefix=$(TEST_CONFIG) -X main.prefixPath=$(TEST_PREFIX) -X main.logPrefix=$(TEST_LOGDIR)' -mod vendor -o target/debug ./...

release:	required
	@mkdir -p target/release
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build --buildmode=$(BUILD_MODE) -v -mod vendor -tags release,$(TAGS) -ldflags '-s -w -X main.version=$(VERSION) -X main.etcPrefix=$(SYSCONFDIR) -X main.prefixPath=$(WORKINGDIR) -X main.logPrefix=$(LOGPREFIXDIR)' -o target/release ./...

debug:	build

install:	release
	@install -d -m 755 $(DESTDIR)$(BINDIR)
	@install -d -m 755 $(DESTDIR)$(SBINDIR)
	@install -d -m 755 $(DESTDIR)$(SYSCONFDIR)
	@install -d -m 755 $(DESTDIR)$(LOGPREFIXDIR)
	@install -d -m 755 $(DESTDIR)$(LOCALSTATEDIR)
	@install -D -m 644 etc/babylon.conf $(DESTDIR)$(SYSCONFDIR)/babylon.conf
	@install -s -m 755 target/release/f9600 $(DESTDIR)$(SBINDIR)

clean:
	@go clean ./...
	@rm -rf target *.out
	@rm -f $(PROJECT)-*.tar.gz $(PROJECT)-*.tar

version:
	@echo $(VERSION)

# Optional make components we add
sinclude .make/*.mk

