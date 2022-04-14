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

package lib

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"path"
)

var (
	logger  *syslog.Writer = nil
	logpath                = ""
	console                = log.New(io.Discard, "", log.LstdFlags)
	verbose                = 0
	argv0                  = path.Base(os.Args[0])
)

// internal specify logging level and path
func openLogger(level int, path string) {
	var err error
	verbose = level
	logpath = path
	LoggerRestart()
	logger, err = syslog.New(syslog.LOG_SYSLOG, argv0)
	if err != nil {
		log.Println(err)
		logger = nil
	}
}

// Reset Logger such as from sighup
func LoggerRestart() {
	if len(logpath) > 0 && logpath != "none" && logpath != "no" && logpath != "/dev/nul" {
		logfile, err := os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
		if err != nil {
			Error(err)
			return
		}
		log.SetOutput(logfile)
		console.SetOutput(os.Stderr)
		console.SetFlags(0) // log.Ltime?
		Notice("logger restart")
	}
}

// Log errors
func Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if logger != nil {
		logger.Err(msg)
	}
	if verbose > 0 {
		console.Println("error:", msg)
	}
	log.Println(msg)
}

// Log failure and exit
func Fail(code int, args ...interface{}) {
	msg := fmt.Sprint(args...)
	if logger != nil {
		logger.Crit(msg)
	}
	if verbose > 0 {
		console.Println("fail:", msg)
	}
	log.Println(msg)
	os.Exit(code)
}

// Log warnings
func Warn(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if logger != nil {
		logger.Warning(msg)
	}
	if verbose > 0 {
		console.Println("warn:", msg)
	}
	log.Println(msg)
}

// Log notices
func Notice(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if logger != nil {
		logger.Notice(msg)
	}
	if verbose > 1 {
		console.Println("notice:", msg)
	}
	log.Println(msg)
}

// Log info
func Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if logger != nil {
		logger.Info(msg)
	}
	if verbose > 1 {
		console.Println("info:", msg)
	}
	log.Println(msg)
}

// Verbose output
func Output(level int, args ...interface{}) {
	if level > verbose {
		return
	}

	msg := fmt.Sprint(args...)
	console.Println("debug:", msg)
	log.Println(msg)
}
