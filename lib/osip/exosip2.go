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

package osip

/*
#cgo CFLAGS: -I/usr/pkg/include -I/usr/local/include
#cgo LDFLAGS: -L/usr/pkg/lib -L/usr/local/lib -leXosip2
#include "exosip2.h"
*/
import "C"
import (
	"fmt"
	"net"
	"strconv"
	"time"
	"unsafe"
)

func boolToInt(enable bool) C.int {
	if enable {
		return 1
	}
	return 0
}

type Config struct {
	// basic server config
	Agent   string
	Ipv6    bool
	Tcp     bool
	Timeout int
	NoText  bool
	NoMedia bool

	// credentials, refresh set if login
	Refresh int
	Server  string
}

type Context struct {
	Config
	context *C.struct_eXosip_t
	Host    string
	Port    int
	Tls     bool

	// internals...
	closed   bool
	active   int // actively registered (rid)
	online   bool
	route    *C.char
	fails    int
	identity string
	username string
	password string
}

type Event struct {
	Context *Context
	Type    EVT_TYPE
	Status  SIP_STATUS
}

func (ctx *Context) Lock() {
	C.eXosip_lock(ctx.context)
}

func (ctx *Context) Unlock() {
	C.eXosip_unlock(ctx.context)
}

func (ctx *Context) Register(identity, user, secret string) error {
	// if no change, skip...
	if ctx.active != -1 && user == ctx.username && secret == ctx.password {
		return nil
	}

	ctx.Unregister()
	ctx.Lock()
	defer ctx.Unlock()
	C.eXosip_clear_authentication_info(ctx.context)
	if len(secret) > 0 {
		cs_user := C.CString(user)
		cs_secret := C.CString(secret)
		defer C.free(unsafe.Pointer(cs_user))
		defer C.free(unsafe.Pointer(cs_secret))
		C.add_credentials(ctx.context, cs_user, cs_secret)
	}
	cs_identity := C.CString(identity)
	defer C.free(unsafe.Pointer(cs_identity))
	ctx.active = int(C.register_identity(ctx.context, cs_identity, ctx.route, C.int(ctx.Refresh)))
	if ctx.active > -1 {
		ctx.username = user
		ctx.password = secret
		ctx.identity = identity
		return nil
	}
	err := fmt.Errorf("registration failed; code=%d", ctx.active)
	ctx.active = -1
	return err
}

func (ctx *Context) Unregister() {
	if ctx.active == -1 {
		return
	}

	ctx.Lock()
	defer ctx.Unlock()
	C.sip_unregister(ctx.context, C.int(ctx.active))
	ctx.active = -1
	ctx.fails = 0
}

func (ctx *Context) Close() {
	if !ctx.closed {
		ctx.closed = true
		ctx.Unregister()
		for ctx.online {
			time.Sleep(time.Second)
		}
		C.eXosip_quit(ctx.context)
		C.release(unsafe.Pointer(ctx.context))

		if ctx.route != nil {
			C.free(unsafe.Pointer(ctx.route))
		}
	}
}

func (ctx *Context) Automatic() {
	ctx.Lock()
	defer ctx.Unlock()
	C.eXosip_automatic_action(ctx.context)
}

func (ctx *Context) automatic_action(evt *C.eXosip_event_t) {
	ctx.Lock()
	defer ctx.Unlock()
	C.eXosip_default_action(ctx.context, evt)
}

func (ctx *Context) ListenAndServe(address string, out chan<- Event) error {
	host, port, err := net.SplitHostPort(address)

	if err != nil {
		return err
	}
	if host == "*" {
		host = ""
	} else if host == "::*" {
		if !ctx.Ipv6 {
			return fmt.Errorf("IPV6 is not enabled")
		}
		host = ""
	}

	family := C.int(0)
	if ctx.Ipv6 {
		family = 1
	}

	proto := C.int(0)
	if ctx.Tcp {
		proto = 1
	}

	ctx.Port, err = strconv.Atoi(port)
	if err != nil {
		return err
	}
	if ctx.Port == 0 {
		ctx.Port = int(C.find_port(ctx.context, proto, C.int(0)))
	}

	ctx.Host = host
	cs_host := C.CString(host)
	defer C.free(unsafe.Pointer(cs_host))

	result := int(C.sip_listen(ctx.context, cs_host, C.int(ctx.Port), family, proto, C.int(0)))
	if result != 0 {
		return fmt.Errorf("sip error: %d", result)
	}

	var event Event = Event{Context: ctx, Type: EVT_STARTUP, Status: SIP_OK}
	out <- event
	for !ctx.closed {
		event = Event{Context: ctx, Type: EVT_TIMEOUT, Status: SIP_OK}
		evt := C.eXosip_event_wait(ctx.context, C.int(ctx.Timeout/1000), C.int(ctx.Timeout%1000))
		if evt == nil {
			out <- event
			ctx.Automatic()
			continue
		}

		response := evt.response
		switch C.evt_type(evt) {
		case C.EXOSIP_REGISTRATION_SUCCESS:
			ctx.fails = 0
			if ctx.online {
				break
			}
			event.Type = EVT_REGISTER
			event.Status = SIP_OK
			if ctx.active != -1 {
				ctx.online = true
			} else {
				ctx.online = false
			}
			out <- event
		case C.EXOSIP_REGISTRATION_FAILURE:
			event.Type = EVT_REGISTER
			ctx.fails = ctx.fails + 1
			if response == nil {
				event.Status = SIP_UNKNOWN
			} else {
				event.Status = SIP_STATUS(response.status_code)
			}
			if ctx.fails < 2 && ctx.active != -1 && (event.Status == SIP_UNAUTHORIZED || event.Status == SIP_PROXY_AUTH_REQUIRED) {
				ctx.automatic_action(evt)
				break
			}
			ctx.online = false
			out <- event
		default:
			ctx.automatic_action(evt)
		}
		C.eXosip_event_free(evt)
	}

	event = Event{Context: nil, Type: EVT_SHUTDOWN, Status: SIP_OK}
	out <- event
	ctx.active = -1
	ctx.online = false
	return nil
}

// sip := osip.New(...)
func New(config Config) *Context {
	ctx := &Context{Config: config, context: C.eXosip_malloc(), Tls: false, closed: false, active: -1}
	C.eXosip_init(ctx.context)
	C.set_option(ctx.context, C.EXOSIP_OPT_ENABLE_IPV6, boolToInt(config.Ipv6))

	if len(ctx.Server) > 0 {
		ctx.route = C.CString(ctx.Server)
	}

	if ctx.Timeout == 0 {
		ctx.Timeout = 500
	}

	if len(config.Agent) > 0 {
		cs_agent := C.CString(config.Agent)
		defer C.free(unsafe.Pointer(cs_agent))
		C.eXosip_set_user_agent(ctx.context, cs_agent)
	}
	return ctx
}