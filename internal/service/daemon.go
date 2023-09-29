// Copyright (C) 2021-2023 David Sugar, Tycho Softworks
// This code is licensed under MIT license
//go:build !linux

package service

import "fmt"

var stopping = false

func Reload(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	return nil
}

func Live(args ...interface{}) error {
	if stopping {
		return nil
	}
	msg := fmt.Sprint(args...)
	if len(msg) > 0 {
		Info(msg)
	}
	return nil
}

func Status(string) error {
	if stopping {
		return fmt.Errorf("already exiting")
	}
	return nil
}

func Stop(args ...interface{}) error {
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
