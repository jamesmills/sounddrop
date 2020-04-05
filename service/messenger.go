package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/tuarrep/sounddrop/message"
	"github.com/tuarrep/sounddrop/util"
	"sync"
)

// Messenger internal messenger service
type Messenger struct {
	Message        chan proto.Message
	log            *logrus.Entry
	listeners      map[byte][]message.Receiver
	listenersMutex sync.Mutex
}

// Stop clean service when stopped by supervisor
func (m *Messenger) Stop() {
	m.log.Info("Messenger stopped.")
}

// Serve main service code
func (m *Messenger) Serve() {
	m.log = util.GetContextLogger("service/messenger.go", "Services/Messenger")
	m.log.Info("Messenger starting...")

	for {
		select {
		case msg := <-m.Message:
			opCode, err := message.FindOpCode(msg)
			if err != nil {
				m.log.Warn(fmt.Sprintf("Unable to find opCode for message %#v", msg))
			}

			if listeners, ok := m.listeners[opCode]; ok {
				for _, listener := range listeners {
					listener.GetChan() <- msg
				}
			}
		}
	}
}

// Register add message listener to listener map for one opCode
func (m *Messenger) Register(opCode byte, listener message.Receiver) {
	m.listenersMutex.Lock()
	if m.listeners == nil {
		m.listeners = make(map[byte][]message.Receiver)
	}
	m.listeners[opCode] = append(m.listeners[opCode], listener)
	m.listenersMutex.Unlock()
}

// RegisterSome add message listener to listener map for multiple opCode
func (m *Messenger) RegisterSome(opCodes []byte, receiver message.Receiver) {
	for _, opCode := range opCodes {
		m.Register(opCode, receiver)
	}
}
