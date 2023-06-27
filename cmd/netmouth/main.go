// Copyright (C) 2023 David Sugar <tychosoft@gmail.com>.
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
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"babylon/lib"
	"babylon/lib/osip"

	"github.com/alexflint/go-arg"
	"github.com/percivalalb/sipuri"
	"gopkg.in/ini.v1"
)

// Argument parser....
type Args struct {
	Config  string `arg:"--config" help:"server config file"`
	Host    string `arg:"--host" help:"server host address" default:""`
	Port    uint16 `arg:"--port" help:"server port" default:"0"`
	Prefix  string `arg:"--prefix" help:"server prefix path"`
	Ipv6    bool   `arg:"-6" help:"enable ipv6 support"`
	Tcp     bool   `arg:"-t" help:"enable tcp sip support"`
	Verbose int    `arg:"-v,--verbose" help:"debugging log level"`
}

// SIP registiry and local config
type Config struct {
	Host     string `ini:"host"`
	Port     uint16 `ini:"port"`
	Ipv6     bool   `ini:"ipv6"`
	Tcp      bool   `ini:"tcp"`
	Refresh  int    `ini:"refresh"`
	Buffer   int    `ini:"events"`
	Timeout  int    `ini:"timeout"`
	Server   string `ini:"server"`
	Identity string `ini:"identity"`
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
	lock   sync.RWMutex
)

func (Args) Version() string {
	return "Version: " + version
}

func (Args) Description() string {
	return "netmouth - TTS speaks sip chat messages"
}

// initialize server and parse arguments
func init() {
	// parse arguments
	for pos, arg := range os.Args {
		switch arg {
		case "--":
			return
		case "-v":
			os.Args[pos] = "--verbose=1"
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
	arg.MustParse(args)

	// setup service
	logPath := logPrefix + "/notmouth.log"
	lib.Logger(args.Verbose, logPath)
	load()
	err := os.Chdir(args.Prefix)
	if err != nil {
		lib.Fail(1, err)
	}
}

// load server config file
func load() {
	// default config
	new_config := Config{
		Host:     args.Host,
		Port:     args.Port,
		Server:   "sip:localhost",
		Identity: "sip:88@localhost",
		Refresh:  300,
		Timeout:  500,
	}

	configs, err := ini.LoadSources(ini.LoadOptions{Loose: true, Insensitive: true}, args.Config, "custom.conf")
	if err == nil {
		// map and reset rom args if not default
		configs.Section("sip").MapTo(&new_config)
		configs.Section("tts").MapTo(&new_config)
		if args.Port != 0 {
			new_config.Port = args.Port
		}

		if len(args.Host) > 0 {
			new_config.Host = args.Host
		}

		if args.Ipv6 {
			new_config.Ipv6 = true
		}

		if args.Tcp {
			new_config.Tcp = true
		}
	} else {
		lib.Error(err)
	}

	// constraints and flags
	if new_config.Host == "*" {
		new_config.Host = ""
	}
	lock.Lock()
	defer lock.Unlock()
	config = &new_config
}

func main() {
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	identity, err := sipuri.Parse(config.Identity)
	if err == nil && len(identity.User()) < 1 {
		err = fmt.Errorf("no user in registration identity")
	}
	if err != nil {
		lib.Fail(99, err, config.Identity)
	}

	lib.Debug(3, "prefix=", args.Prefix, ", bind=", address)
	register := sipuri.New(identity.User(), identity.Host()+":"+identity.Port())
	sip := osip.New(osip.Config{
		Agent:    "netmouth/" + version,
		Ipv6:     config.Ipv6,
		Server:   config.Server,
		Identity: register.User(),
		Secret:   identity.Password(),
		Refresh:  config.Refresh,
		NoMedia:  true,
	})

	// signal handler...
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		defer sip.Close()
		for {
			switch <-signals {
			case os.Interrupt: // sigint/ctrl-c
				fmt.Println()
				return
			case syscall.SIGTERM: // normal exit
				return
			case syscall.SIGHUP: // cleanup
				lib.Info("reload service")
				lib.LoggerRestart()
				runtime.GC()
				load()
			}
		}
	}()

	events := make(chan osip.Event, config.Buffer)
	go func(ch <-chan osip.Event) {
		for {
			event := <-ch
			lib.Debug(2, "event type: ", event.Type)
			if event.Type == osip.EVT_SHUTDOWN {
				return
			}
		}
	}(events)

	lib.Info("start service on ", address)
	err = sip.Listen(address, events)
	if err != nil {
		lib.Fail(1, err)
	}
	lib.Info("stop service")
}
