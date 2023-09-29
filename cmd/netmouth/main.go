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

	"babylon/internal/service"

	"github.com/alexflint/go-arg"
	"github.com/percivalalb/sipuri"
	"gopkg.in/ini.v1"

	osip "babylon/internal/exosip2"
	htgotts "github.com/hegedustibor/htgo-tts"
	handlers "github.com/hegedustibor/htgo-tts/handlers"
	//voices "github.com/hegedustibor/htgo-tts/voices"
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
	Secret   string `ini:"secret"`
	User     string `ini:"user"`

	// tts values
	Proxy    string `ini:"proxy"`
	Native   bool   `ini:"native"`
	Language string `ini:"language"`

	// more internal...
	register string
	route    string
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
	service.Logger(args.Verbose, logPath)
	load()
	err := os.Chdir(args.Prefix)
	if err != nil {
		service.Fail(1, err)
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
		Language: "en",
		Native:   true,
	}

	configs, err := ini.LoadSources(ini.LoadOptions{Loose: true, Insensitive: true}, args.Config, args.Prefix+"/custom.conf")
	if err == nil {
		// map and reset rom args if not default
		configs.Section("sip").MapTo(&new_config)
		configs.Section("tts").MapTo(&new_config)
		configs.Section("netmouth").MapTo(&new_config)
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
		service.Error(err)
	}

	identity, err := sipuri.Parse(new_config.Identity)
	if err == nil && len(identity.User()) < 1 && len(new_config.User) < 1 {
		err = fmt.Errorf("no user for registration identity")
	}
	if err != nil {
		service.Fail(99, err, config.Identity)
	}

	new_config.register = sipuri.New(identity.User(), identity.Host()).String()
	if len(new_config.Secret) < 1 {
		new_config.Secret = identity.Password()
	}
	if len(new_config.User) < 1 {
		new_config.User = identity.User()
	}

	route, err := sipuri.Parse(new_config.Server)
	if err != nil {
		service.Fail(99, err, new_config.Server)
	}
	new_config.route = "sip:" + route.Host()

	// constraints and flags
	if new_config.Host == "*" {
		new_config.Host = ""
	}
	lock.Lock()
	defer lock.Unlock()
	config = &new_config
}

func main() {
	cache := args.Prefix + "/tts"
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	route, err := sipuri.Parse(config.Server)
	if err != nil {
		service.Fail(99, err, config.Server)
	}

	service.Debug(3, "prefix=", args.Prefix, ", bind=", address)
	service.Debug(3, "server=", "sip:"+route.Host(), ", identity=", config.register)
	os.RemoveAll(cache)
	os.Mkdir(cache, 0770)
	defer os.RemoveAll(cache)

	sip := osip.New(osip.Config{
		Agent:   "netmouth/" + version,
		Ipv6:    config.Ipv6,
		Server:  config.route,
		Refresh: config.Refresh,
		Allows:  "OPTIONS,MESSAGE",
	})

	// signal handler...
	signals := make(chan os.Signal, 1)
	texts := make(chan string, 32)
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
				service.Reload("reload service")
				service.LoggerRestart()
				runtime.GC()
				load()
				if sip.SetRoute(config.route) {
					service.Info("changed route to ", config.route)
				}
				sip.Register(config.Identity, config.User, config.Secret)
				texts <- "-reload-"
				service.Live()
			}
		}
	}()

	go func(say <-chan string) {
		var handler handlers.PlayerInterface
		handler = &handlers.Native{}
		if !config.Native {
			handler = &handlers.MPlayer{}
		}

		speach := htgotts.Speech{Folder: args.Prefix + "/tts", Language: config.Language, Handler: handler, Proxy: config.Proxy}
		for {
			text := <-say
			switch text {
			case "-reload-":
				os.RemoveAll(cache)
			case "":
				continue
			default:
				err := speach.Speak(text)
				if err != nil {
					service.Error(err)
				}
			}
		}
	}(texts)

	events := make(chan osip.Event, config.Buffer)
	go func(ch <-chan osip.Event, say chan<- string) {
		defer service.Stop("stop service")
		for {
			event := <-ch
			ctx := event.Context
			service.Debug(3, "event type: ", event.Type)
			switch event.Type {
			case osip.EVT_SHUTDOWN:
				return
			case osip.EVT_STARTUP:
				service.Live("service started ", ctx.Host, ":", ctx.Port)
				ctx.Register(config.Identity, config.User, config.Secret)
			case osip.EVT_REGISTER:
				if event.Status != osip.SIP_OK {
					service.Status("offline")
					service.Error("registration failure; status=", event.Status)
				} else {
					service.Status("online")
					service.Info("service online; id=", ctx.GetIdentity())
				}
			case osip.EVT_MESSAGE:
				if event.Status != osip.SIP_OK {
					break
				}
				if event.Content == "message/imdn+xml" {
					event.Reply(osip.SIP_OK)
					break
				}
				if event.Content != "text/plain" {
					service.Debug(2, "ignored message input ", event.Content)
					event.Reply(osip.SIP_NOT_ACCEPTABLE_HERE)
					break
				}
				text := string(event.Body)
				service.Debug(2, "message from ", event.From, "; text=", text)
				event.Reply(osip.SIP_OK)
				say <- text
			}
		}
	}(events, texts)

	err = sip.ListenAndServe(address, events)
	if err != nil {
		service.Fail(1, err)
	}
}
