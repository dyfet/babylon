# Copyright (C) 2021 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by a MIT-style license
# that can be found in the included LICENSE.md file.

# cross build overrides
ifdef ARCH
GOARCH:= $(shell .make/cross.arch $(ARCH))
endif

