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
	"sync"
	"syscall"

	"babylon/lib"
	"github.com/alexflint/go-arg"
	"gopkg.in/ini.v1"
)

// Argument parser....
type Args struct {
	Host string `arg:"--host" help:"server host address" default:""`
	Port uint16 `arg:"--port" help:"server port" default:"9600"`
	// TODO: future TCP/TLS option
	LogLevel    int  `arg:"-x,--debug" help:"debugging log level"`
	ShowVersion bool `arg:"--version" help:"display version and exit"`
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
	etcPrefix string = "/etc"
	varPrefix string = "/var/lib/babylon"
	logPrefix string = "/var/log"
	version   string = "unknown"
	buildType string = "release"

	// globals
	args    *Args   = &Args{}
	config  *Config = nil
	running         = true
	exiting         = 0
)

func (Args) Version() string {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		return "Version: " + version + "\nConfig: " + etcPrefix + "/babylon.conf\nPrefix: " + varPrefix
	}
	return ""
}

func (Args) Description() string {
	return "f9600 " + version + "\nDavid Sugar <tychosoft@gmail.com>\nProvides Fujitsu F9600 service daemon and command access\n"
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

	configs, err := ini.LoadSources(ini.LoadOptions{Loose: true, Insensitive: true}, etcPrefix+"/babylon.conf", "custom.conf")
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

	if buildType == "debug" {
		os.Remove(logPath)
	}

	lib.Logger(args.LogLevel, logPath)
	err := os.Chdir(varPrefix)
	if err != nil {
		lib.Fail(1, err)
	}

	config = load()
	lib.Debug(4, "config=", config)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	lib.Info("start service")

	// bind sockets and connections
	tcp, err := net.Listen("tcp", config.Address)
	if err != nil {
		lib.Fail(2, err)
	}

	// action group
	ag := new(sync.WaitGroup)
	ag.Add(1)
	err = mml.Configure(config)
	if err != nil {
		lib.Fail(3, err)
	}
	go mml.Startup(ag, config)

	// session group
	sg := new(sync.WaitGroup)
	sg.Add(2)
	go manager.Startup(sg)
	go func() { // lambda for tcp local...
		defer sg.Done()
		lib.Info("server started")
		for {
			client, err := tcp.Accept()
			if !running {
				break
			}
			if err != nil {
				lib.Error(err)
				continue
			}
			fmt.Fprint(client, config.Banner+"\r\n")
			NewSession(client)
		}
		lib.Info("server stopped")
	}()

	for running {
		signal := <-signals
		switch signal {
		case os.Interrupt: // sigint/ctrl-c
			exiting = 1
			running = false
			fmt.Println()
		case syscall.SIGTERM: // normal exit
			running = false
		case syscall.SIGHUP: // cleanup
			lib.Info("reload service")
			lib.LoggerRestart()
			runtime.GC()
			config = load()
		}
	}

	// shutdown sessions
	tcp.Close()
	manager.Shutdown()
	sg.Wait()

	// shutdown actions
	mml.Shutdown()
	ag.Wait()
	lib.Info("stopped service; reason=", exiting)
	os.Exit(exiting)
}
