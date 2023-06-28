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
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "unsafe"

func (ctx *Context) GetSchema() string {
	if ctx.Tls {
		return "sips:"
	}
	return "sip:"
}

func (ctx *Context) GetIdentity() string {
	ctx.Lock()
	defer ctx.Unlock()

	if ctx.active == -1 {
		return ""
	}

	return ctx.identity
}

func (ctx *Context) IsOpen() bool {
	return !ctx.closed
}

func (ctx *Context) IsActive() bool {
	return ctx.active != -1
}

func (ctx *Context) IsOnline() bool {
	return ctx.online
}

func (ctx *Context) SetRoute(route string) bool {
	ctx.Lock()
	defer ctx.Unlock()

	if ctx.route == nil && len(route) < 1 {
		return false
	}

	cs_route := C.CString(route)
	if ctx.route != nil && C.strcmp(ctx.route, cs_route) == 0 {
		C.free(unsafe.Pointer(cs_route))
		return false
	}

	if ctx.route != nil {
		C.free(unsafe.Pointer(ctx.route))
	}

	if len(route) > 0 {
		ctx.route = cs_route
	} else {
		C.free(unsafe.Pointer(cs_route))
		ctx.route = nil
	}
	return true
}
