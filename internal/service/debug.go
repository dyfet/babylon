// Copyright (C) 2023 David Sugar, Tycho Softworks
// This code is licensed under MIT license
//go:build debug

package service

import (
	"os"
)

func Logger(level int, path string) {
	os.Remove(path)
	openLogger(level, path)
}

// Debug output
func Debug(level int, args ...interface{}) {
	Output(level, args...)
}

func IsDebug() bool {
	return true
}
