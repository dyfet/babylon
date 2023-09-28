// Copyright (C) 2023 David Sugar, Tycho Softworks
// This code is licensed under MIT license
//go:build !debug

package service

func Logger(level int, path string) {
	openLogger(level, path)
}

func Debug(level int, args ...interface{}) {
}

func IsDebug() bool {
	return false
}
