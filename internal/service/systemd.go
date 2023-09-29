// Copyright (C) 2021-2023 David Sugar, Tycho Softworks
// This code is licensed under MIT license
//go:build linux

package service

import (
	"fmt"

	"github.com/coreos/go-systemd/daemon"
)

var stopping = false

func Reload(args ...interface{}) error {
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

func Live(args ...interface{}) error {
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

func Status(status string) error {
	if stopping {
		return nil
	}
	msg := fmt.Sprintf("%s\nSTATUS=%s\n", daemon.SdNotifyReady, status)
	_, err := daemon.SdNotify(false, msg)
	return err
}

func Stop(args ...interface{}) error {
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
