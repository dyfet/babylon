// Copyright (C) 2021-2023 David Sugar <tychosoft@gmail.com>.
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
//go:build !linux

package lib

import "fmt"

var stopping = false

func DaemonReload(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	return nil
}

func DaemonLive(args ...interface{}) error {
	if stopping {
		return nil
	}
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	return nil
}

func DaemonStatus(string) error {
	return nil
}

func DaemonStop(args ...interface{}) error {
	if stopping {
		return nil
	}
	stopping = true
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	return nil
}

func Watchdog() error {
	return nil
}
