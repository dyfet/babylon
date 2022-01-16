# Copyright (C) 2020-2021 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by a MIT-style license
# that can be found in the included LICENSE.md file.

.PHONY: dist

dist:	required
	@rm -f $(PROJECT)-*.tar.gz $(PROJECT)-*.tar
	@git archive -o $(PROJECT)-$(VERSION).tar --format tar --prefix=$(PROJECT)-$(VERSION)/ v$(VERSION) 2>/dev/null || git archive -o $(PROJECT)-$(VERSION).tar --format tar --prefix=$(PROJECT)-$(VERSION)/ HEAD
	@if test -f vendor/modules.txt ; then \
		tar --transform s:^:$(PROJECT)-$(VERSION)/: --append --file=$(PROJECT)-$(VERSION).tar vendor ; fi
	@gzip $(PROJECT)-$(VERSION).tar


