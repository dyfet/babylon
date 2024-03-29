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
	"babylon/internal/service"
)

// session manager
type Manager struct {
	sessions map[*Session]bool
	register chan *Session
	release  chan *Session
	shutdown chan bool
}

var (
	// singleton...
	manager = Manager{
		sessions: make(map[*Session]bool),
		register: make(chan *Session),
		release:  make(chan *Session),
		shutdown: make(chan bool),
	}
)

// register a new session with the manager
func (manager *Manager) Register(s *Session) {
	manager.register <- s
}

// release an exiting session from the manager
func (manager *Manager) Release(s *Session) {
	manager.release <- s
}

// shutdown manager
func (manager *Manager) Shutdown() {
	manager.shutdown <- true
}

// process manager api until clean shutdown
func (manager *Manager) Startup() {
	service.Debug(1, "manager running")
	for {
		select {
		case session := <-manager.register:
			manager.sessions[session] = true
			service.Debug(2, "adding session ", session.Remote)
		case session := <-manager.release:
			if _, ok := manager.sessions[session]; ok {
				delete(manager.sessions, session)
				service.Debug(2, "remove session ", session.Remote)
			} else {
				service.Warn("unknown session ", session.Remote)
			}
		case <-manager.shutdown:
			for s := range manager.sessions {
				s.Close()
			}
			return
		}
	}
}
