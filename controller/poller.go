package controller

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"
)

type Poller struct {
	statusProtocol    StatusProtocol
	runner            *Runner
	currentStatus     *ServerStatus
	currentStatusJson []byte
	lock              sync.RWMutex
}

func NewPoller(statusProtocol StatusProtocol, runner *Runner) * Poller{
	p := Poller{
		statusProtocol: statusProtocol,
		runner:         runner,
	}
	return &p
}

func (s *Poller) GetStatus() *ServerStatus {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.currentStatus
}

func (s *Poller) GetStatusJson() []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.currentStatusJson
}

func (s *Poller) StartPolling() {
	go func() {
		for {
			var stJson []byte
			var st *ServerStatus
			// server definitely not running, so don't bother getting status
			if s.runner.Status == "OFFLINE" {
				st = nil
				stJson = SERVER_OFFLINE
			} else { // we're running, yay!
				log.Printf("Polling %s at %s", s.runner.Name, s.statusProtocol)
				st := s.statusProtocol.GetStatus()
				if st != nil {
					stJson, _ = json.Marshal(st)
				} else { // possibly still starting
					stJson = []byte(strings.Replace(SERVER_OFFLINE_STR, "OFFLINE", s.runner.Status, 1))
				}
			}
			s.lock.Lock()
			s.currentStatus = st
			s.currentStatusJson = stJson
			s.lock.Unlock()
			if s.runner.Status == "STARTING" {
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}
