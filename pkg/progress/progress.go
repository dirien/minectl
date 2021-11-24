package progress

import (
	"sync"
	"time"

	"github.com/minectl/pkg/logging"
)

type Indicator struct {
	Delay        time.Duration
	StartMessage string
	Message      string
	FinalMessage string
	ErrorMessage string
	logging      *logging.MinectlLogging
	stopChan     chan bool
	mu           *sync.RWMutex
	active       bool
}

func NewIndicator(message string, logging *logging.MinectlLogging) *Indicator {
	return &Indicator{
		Message:  message,
		Delay:    30 * time.Second,
		stopChan: make(chan bool),
		mu:       &sync.RWMutex{},
		active:   false,
		logging:  logging,
	}
}

func (i *Indicator) Start() {
	i.mu.Lock()
	if i.active {
		i.mu.Unlock()
		return
	}
	i.active = true
	i.mu.Unlock()

	go func() {
		for {
			select {
			case <-i.stopChan:
				return
			default:
				i.mu.Lock()
				if !i.active {
					i.mu.Unlock()
					return
				}
				i.logging.RawMessage(i.Message)
			}
			i.mu.Unlock()
			time.Sleep(i.Delay)
		}
	}()
}

func (i *Indicator) StopE(err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.stopChan <- true
	if err != nil {
		i.logging.RawMessage(i.ErrorMessage)
	} else {
		i.logging.RawMessage(i.FinalMessage)
	}
}
