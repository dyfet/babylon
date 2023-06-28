// Copyright (C) 2021-2022 David Sugar <tychosoft@gmail.com>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>

#include <eXosip2/eXosip.h>

void set_option(struct eXosip_t *ctx, int option, int value) {
	eXosip_set_option(ctx, option, &value);
}

void add_credentials(struct eXosip_t *ctx, const char *user, const char *secret) {
    eXosip_add_authentication_info(ctx, user, user, secret, NULL, NULL);
}

int register_identity(struct eXosip_t *ctx, const char *identity, const char *route, int expires) {
    osip_message_t *msg = NULL;
    int rid = eXosip_register_build_initial_register(ctx, identity, route, NULL, expires, &msg);
    if(rid > -1) {
        eXosip_register_send_register(ctx, rid, msg);
    }
    return rid;
}

int find_port(struct eXosip_t *ctx, int proto, int tls) {
	if(proto)
		proto = IPPROTO_TCP;
	else
		proto = IPPROTO_UDP;

	int port = 5060;
	for(;;) {
		port = eXosip_find_free_port(ctx, port, proto);
		if(tls && !(port & 0x01))
			++port;
		else if(!tls && (port & 0x01))
			++port;
		else
			return port;
	}
}

int sip_listen(struct eXosip_t *ctx, const char *host, int port, int family, int proto, int tls) {
	if(family)
		family = AF_INET6;
	else
		family = AF_INET;

	if(proto)
		proto = IPPROTO_TCP;
	else
		proto = IPPROTO_UDP;

	if(!host || !host[0])
		host = NULL;

	return eXosip_listen_addr(ctx, proto, host, port, family, tls);
}

int evt_type(eXosip_event_t *evt) {
    return evt->type;
}

void sip_unregister(struct eXosip_t *ctx, int rid) {
    osip_message_t *msg = NULL;
    int res = eXosip_register_build_register(ctx, rid, 0, &msg);
    if(res > -1)
        eXosip_register_send_register(ctx, rid, msg);
}

void release(void *p) {
	if(p != NULL) {
        free(p);
/*
		if(osip_free_func)
			osip_free_func(p);
		else
			free(p);
*/
	}
}

