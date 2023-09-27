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
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"babylon/internal/service"
)

// representation of an accepted client session when started
type Session struct {
	Remote string
	socket net.Conn
	result chan string
	update time.Time
}

// print into a client session
func (s *Session) Print(args ...interface{}) (int, error) {
	return fmt.Fprint(s.socket, args...)
}

// print a new line into a client session
func (s *Session) Println(args ...interface{}) (int, error) {
	msg := fmt.Sprint(args...)
	msg += "\r\n"
	return fmt.Fprint(s.socket, msg)
}

// post a result to the waiting client session
func (s *Session) Result(text string) error {
	s.result <- text
	return nil
}

// close session, forces created session to exit
func (s *Session) Close() {
	s.socket.Close()
}

// execute client requests in a go routine...
func (s *Session) requests() {
	defer s.Close()
	defer close(s.result)

	input := bufio.NewReader(s.socket)
	for {
		// prompt for and get input
		fmt.Fprint(s.socket, "mml>")
		line, err := input.ReadString('\n')
		if err != nil {
			break
		}

		// process command or send
		line = strings.Trim(line, "\r\n")
		service.Debug(5, "mml request ", line)
		if line == "quit" || line == "bye" {
			break
		}

		// get result after sending command somewhere
		mml.Request(s, line)
		text := <-s.result
		s.update = time.Now()
		if len(text) > 0 {
			service.Error(fmt.Errorf("MML Error on %s %s", s.Remote, text))
		}
	}
	manager.Release(s)
}

// construct client i/o and register
func NewSession(connect net.Conn) *Session {
	s := &Session{
		Remote: fmt.Sprint(connect.RemoteAddr()),
		socket: connect,
		result: make(chan string),
		update: time.Now(),
	}
	manager.Register(s)
	go s.requests()
	return s // usually ignored?
}
