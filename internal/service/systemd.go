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
//go:build linux

package service

import (
	"fmt"

	"github.com/coreos/go-systemd/daemon"
)

var stopping = false

func DaemonReload(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	if stopping {
		return nil
	}
	_, err := daemon.SdNotify(false, daemon.SdNotifyReloading)
	return err
}

func DaemonLive(args ...interface{}) error {
	if stopping {
		return nil
	}
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	return err
}

func DaemonStatus(status string) error {
	if stopping {
		return nil
	}
	msg := fmt.Sprintf("%s\nSTATUS=%s\n", daemon.SdNotifyReady, status)
	_, err := daemon.SdNotify(false, msg)
	return err
}

func DaemonStop(args ...interface{}) error {
	if stopping {
		return nil
	}
	stopping = true
	_, err := daemon.SdNotify(true, daemon.SdNotifyStopping)
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	return err
}

func Watchdog() error {
	if stopping {
		return nil
	}
	_, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
	return err
}
