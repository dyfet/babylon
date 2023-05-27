# Copyright (C) 2020-2021 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by a MIT-style license
# that can be found in the included LICENSE.md file.

# Testing paths can be set for debug
ifdef WORKINGDIR
TEST_PREFIX := $(WORKINGDIR)
else
ifdef LOCALSTATEDIR
TEST_PREFIX := $(LOCALSTATEDIR)/lib/$(PROJECT)
else
TEST_PREFIX := $(TESTDIR)
endif
endif

ifdef SYSCONFDIR
TEST_CONFIG := $(SYSCONFDIR)
else
TEST_CONFIG := $(TESTDIR)
endif

ifdef LOGPREFIXDIR
TEST_LOGDIR := $(LOGPREFIX)
else
ifdef LOCALSTATEDIR
TEST_LOGDIR := $(LOCALSTATEDIR)/log
else
TEST_LOGDIR := $(TESTDIR)
endif
endif

ifdef APPDATADIR
TEST_APPDIR := $(APPDATADIR)
else
ifdef DATADIR
TEST_APPDIR := $(DATADIR)/$(PROJECT)
else
ifdef PREFIX
TEST_APPDIR := $(PREFIX)/share/$(PROJECT)
else
TEST_APPDIR := $(TESTDIR)
endif
endif
endif

# Project overrides, starting with prefix install
TAGS =

ifndef	DESTDIR
DESTDIR =
endif

ifndef	PREFIX
PREFIX := /usr/local
endif

ifndef	BINDIR
BINDIR := $(PREFIX)/bin
endif

ifndef	SBINDIR
SBINDIR := $(PREFIX)/sbin
endif

ifndef	LIBDIR
LIBDIR := $(PREFIX)/lib
endif

ifndef	LIBDATADIR
LIBDATADIR := $(PREFIX)/lib
endif

ifndef	DATADIR
DATADIR := $(PREFIX)/share
endif

ifndef	MANDIR
MANDIR := $(PREFIX)/man
endif

ifndef	LOCALSTATEDIR
LOCALSTATEDIR := $(PREFIX)/var
endif

ifndef	SYSCONFDIR
SYSCONFDIR := $(PREFIX)/etc
endif

ifndef	LOGPREFIXDIR
LOGPREFIXDIR := $(LOCALSTATEDIR)/log
endif

ifndef	WORKINGDIR
WORKINGDIR := $(LOCALSTATEDIR)/lib/$(PROJECT)
endif

ifndef	APPDATADIR
APPDATADIR := $(DATADIR)/$(PROJECT)
endif
