# Copyright (C) 2020-2021 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by a MIT-style license
# that can be found in the included LICENSE.md file.

.PHONY: lint vet fix test cover release

ifndef	RELEASE_TYPE
RELEASE_TYPE = shared
endif

lint:	required
	@go fmt ./...
	@go mod tidy
	@staticcheck ./...

vet:	required
	@go vet ./...

fix:	required
	@go fix ./...

test:	vet
	@go test ./...

cover:	vet
	@go test -coverprofile=coverage.out ./...

go.sum:	go.mod
	@go mod tidy

# select release by type (static or shared)
release: release-$(RELEASE_TYPE)
	@echo "releasing as $(RELEASE_TYPE)"

# if no vendor directory (clean) or old in git checkouts
vendor:	go.sum
	@if test -d .git ; then \
		rm -rf vendor ;\
		go mod vendor ;\
	elif test ! -d vendor ; then \
		go mod vendor ;\
	else \
		touch vendor ;\
	fi
