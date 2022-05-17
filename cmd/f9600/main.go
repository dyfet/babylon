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

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"babylon/lib"
	"github.com/alexflint/go-arg"
	"gopkg.in/ini.v1"
)

// Argument parser....
type Args struct {
	Config string `arg:"--config" help:"server config file"`
	Host   string `arg:"--host" help:"server host address" default:""`
	Port   uint16 `arg:"--port" help:"server port" default:"9600"`
	// TODO: future TCP/TLS option
	Prefix  string `arg:"--prefix" help:"server prefix path"`
	Verbose int    `arg:"-v,--verbose" help:"debugging log level"`
}

// F9600 config object
type Config struct {
	Banner  string `ini:"banner"`
	Device  string `ini:"device"`
	Speed   int    `ini:"speed"`
	Host    string `ini:"host"`
	Port    uint16 `ini:"port"`
	User    string `ini:"user"`
	Pass    string `ini:"pass"`
	Address string `ini:"-"`
}

var (
	// bind Makefile config
	prefixPath = "/var/lib/babylon"
	etcPrefix  = "/etc"
	logPrefix  = "/var/log"
	version    = "unknown"

	// globals
	args   *Args   = &Args{Prefix: prefixPath, Config: etcPrefix + "/babylon.conf"}
	config *Config = nil
)

func (Args) Version() string {
	return "Version: " + version
}

func (Args) Description() string {
	return "f9600 - provides Fujitsu F9600 service daemon and command access"
}

// initialize server
func init() {
	for pos, arg := range os.Args {
		switch arg {
		case "--":
			return
		case "-vv":
			os.Args[pos] = "--verbose=2"
		case "-vvv":
			os.Args[pos] = "--verbose=3"
		case "-vvvv":
			os.Args[pos] = "--verbose=4"
		case "-vvvvv":
			os.Args[pos] = "--verbose=5"
		}
	}
}

// load server config file
func load() *Config {
	// default config
	config := Config{
		Banner: "Welcome to F9600 pbx",
		Device: "/dev/ttyUSB0",
		Speed:  9600,
		Host:   args.Host,
		Port:   args.Port,
		User:   "admin",
		Pass:   "admin",
	}

	configs, err := ini.LoadSources(ini.LoadOptions{Loose: true, Insensitive: true}, args.Config, "custom.conf")
	if err == nil {
		// map and reset rom args if not default
		configs.Section("f9600").MapTo(&config)
		if args.Port != 9600 {
			config.Port = args.Port
		}

		if len(args.Host) > 0 {
			config.Host = args.Host
		}
	} else {
		lib.Error(err)
	}

	// constraints and flags
	if config.Host == "*" {
		config.Host = ""
	}
	config.Address = fmt.Sprintf("%s:%v", config.Host, config.Port)
	return &config
}

func main() {
	logPath := logPrefix + "/f9600.log"
	arg.MustParse(args)
	// TODO: constraints on parsed arguments

	lib.Logger(args.Verbose, logPath)
	err := os.Chdir(args.Prefix)
	if err != nil {
		lib.Fail(1, err)
	}

	// config service
	config = load()
	lib.Debug(4, "config=", config)
	tcp, err := net.Listen("tcp", config.Address)
	if err != nil {
		lib.Fail(2, err)
	}
	err = mml.Configure(config)
	if err != nil {
		lib.Fail(3, err)
	}

	// signal handler...
	running := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		switch <-signals {
		case os.Interrupt: // sigint/ctrl-c
			fmt.Println()
			running = false
			tcp.Close()
			return
		case syscall.SIGTERM: // normal exit
			running = false
			tcp.Close()
			return
		case syscall.SIGHUP: // cleanup
			lib.Info("reload service")
			lib.LoggerRestart()
			runtime.GC()
			config = load()
		}
	}()

	// run service
	lib.Info("start service")
	go mml.Startup(config)
	go manager.Startup()
	for {
		client, err := tcp.Accept()
		if err != nil {
			if running {
				lib.Error(err)
			}
			running = false
			break
		}

		fmt.Fprint(client, config.Banner+"\r\n")
		NewSession(client)
	}

	// shutdown sessions
	tcp.Close()
	manager.Shutdown()
	mml.Shutdown()
	lib.Info("stopped service")
}
