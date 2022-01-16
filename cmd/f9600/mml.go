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
	"babylon/lib"
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// mml command request object
type mmlRequest struct {
	command string
	session *Session
}

// representation of f9600 mml serial session
type MML struct {
	requests chan mmlRequest
	port     *serial.Port
}

var (
	// singleton
	mml = MML{
		requests: make(chan mmlRequest),
	}
)

func (mml *MML) framer(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\002')
	if err == nil {
		line, err = reader.ReadString('\003')
	}
	return line, err
}

func (mml *MML) copier(reader *bufio.Reader, session *Session) error {
	var line string
	var err error = nil

	for {
		line, err = mml.framer(reader)
		if err != nil {
			break
		}
		session.Println(line)
		if line[0:5] == " END " {
			break
		}
		if line[0:5] == " ERR-" {
			err = fmt.Errorf("%s", line[6:])
			break
		}
	}
	return err
}

func (mml *MML) password(reader *bufio.Reader, pass string) error {
	var line string
	var err error = nil

	for {
		line, err = mml.framer(reader)
		if err != nil {
			break
		}
		if line[0:11] == " PASSWORD :" {
			fmt.Fprint(mml.port, pass+"\r")
			break
		}
	}
	return err
}

func (mml *MML) Request(s *Session, cmd string) {
	request := mmlRequest{
		command: cmd,
		session: s,
	}
	if !running {
		return
	}
	mml.requests <- request
}

// attempt configuration of mml
func (mml *MML) Configure(config *Config) error {
	parms := &serial.Config{
		Name:        config.Device,
		Baud:        config.Speed,
		ReadTimeout: time.Second * 2,
		Size:        8,
		Parity:      serial.ParityEven,
		StopBits:    serial.Stop1,
	}
	port, err := serial.OpenPort(parms)
	if err != nil {
		return err
	}
	mml.port = port
	lib.Info("opened ", config.Device)
	return nil
}

// start mml session
func (mml *MML) Startup(wg *sync.WaitGroup, config *Config) {
	defer wg.Done()
	active := false
	reader := bufio.NewReader(mml.port)
	lib.Info("mml startup")
	for {
		request, ok := <-mml.requests
		if !ok { // if shutdown, exiting
			break
		}

		if !running { // flush further pending if stopping...
			continue
		}

		session := request.session
		if !active {
			count, err := fmt.Fprint(mml.port, "\r")
			if err != nil || count < 1 {
				session.Println(" ERR-Offline")
				session.Result("offline in login")
				continue
			}
			time.Sleep(time.Second)
			mml.port.Flush()
			count, err = fmt.Fprint(mml.port, "login,"+config.User+"\r")
			if err != nil || count < 1 {
				session.Println(" ERR-Offline")
				session.Result("offline in login")
				continue
			}

			time.Sleep(time.Second)
			if mml.password(reader, config.Pass) != nil {
				session.Println(" ERR-Offline")
				session.Result("pbx login failed")
				continue
			}
			time.Sleep(time.Second)
			mml.port.Flush()
			active = true
		}

		count, err := fmt.Fprint(mml.port, request.command+"\r")
		if err != nil {
			if err.Error() == "EOF" {
				active = false
				err = fmt.Errorf("offline in send")
				session.Println(" ERR-Offline")
			} else {
				session.Println(" ERR-", err.Error())
			}
			session.Result(err.Error())
			continue
		}
		if count < 1 {
			session.Println(" ERR-Offline")
			session.Result("no output")
			continue
		}
		err = mml.copier(reader, session)
		if err == nil {
			session.Result("")
		} else {
			if err.Error() == "EOF" {
				err = fmt.Errorf("offline in recv")
				session.Println(" ERR-Offline")
				active = false
			}
			session.Result(err.Error())
		}
	}
	lib.Info("mml shutdown")
}

// close session, forces mml exit
func (mml *MML) Shutdown() {
	mml.port.Close()
	close(mml.requests)
}
