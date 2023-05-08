# Copyright (C) 2020-2021 David Sugar <tychosoft@gmail.com>.
# Use of this source code is governed by a MIT-style license
# that can be found in the included LICENSE.md file.

.PHONY:	certs

certs:
	@rm -f test/server.key test/server.crt
	@openssl ecparam -genkey -name secp384r1 -out test/server.key
	@openssl req -new -x509 -sha256 -key test/server.key -out test/server.crt -days 3650

