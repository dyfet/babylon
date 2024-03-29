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
#include <osipparser2/osip_port.h>

typedef struct content_type {
    char *ctype;
    char *subtype;
} content_type_t;

void set_option(struct eXosip_t *ctx, int option, int value) {
	eXosip_set_option(ctx, option, &value);
}

char *get_url(osip_uri_t *uri) {
    char *str = NULL;
    osip_uri_to_str(uri, &str);
    return str;
}

char *get_subject(osip_message_t *msg, int index) {
    osip_header_t *header = NULL;
    osip_message_get_subject(msg, index, &header);
    if(header && header->hvalue)
        return header->hvalue;
    return NULL;
}

int get_expires(osip_message_t *msg, int index) {
    osip_header_t *header = NULL;
    osip_message_get_expires(msg, index, &header);
    if(header && header->hvalue)
        return atoi(header->hvalue);
    return -1;
}

osip_message_t *message_response(struct eXosip_t *ctx, int tid, int status) {
    osip_message_t *msg = NULL;
    eXosip_message_build_answer(ctx, tid, status, &msg);
    return msg;
}

osip_message_t *call_response(struct eXosip_t *ctx, int tid, int status) {
    osip_message_t *msg = NULL;
    eXosip_call_build_answer(ctx, tid, status, &msg);
    return msg;
}

content_type_t get_content(osip_message_t *msg) {
    content_type_t res = {NULL, NULL};
    osip_content_type_t *ctype = osip_message_get_content_type(msg);
    if(ctype == NULL || ctype->type == NULL)
        return res;

    res.ctype = ctype->type;
    res.subtype = ctype->subtype;
    return res;
}

osip_body_t *get_body(osip_message_t *msg, int index) {
    osip_body_t *body = NULL;
    osip_message_get_body(msg, index, &body);
    return body;
}

void add_credentials(struct eXosip_t *ctx, const char *user, const char *secret) {
    eXosip_add_authentication_info(ctx, user, user, secret, NULL, NULL);
}

int register_identity(struct eXosip_t *ctx, const char *identity, const char *route, int expires, const char *allow, const char *accept, const char *encoding) {
    osip_message_t *msg = NULL;
    int rid = eXosip_register_build_initial_register(ctx, identity, route, NULL, expires, &msg);
    if(rid > -1) {
        if(allow)
            osip_message_set_header(msg, ALLOW, allow);
        if(accept)
            osip_message_set_header(msg, ACCEPT, accept);
        if(encoding)
            osip_message_set_header(msg, ACCEPT_ENCODING, encoding);
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
    if(osip_free_func)
        osip_free_func(p);
    else
        free(p);
}

