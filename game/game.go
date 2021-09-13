// Package game stores game state.
package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/inahga/aightreader/device"
)

// Being tick based avoids the condition where user presses a bunch of keys at
// the same time, and the incoming MIDI messages coincidentally arrive in the
// correct order.
const tick = 200 * time.Millisecond

type state struct {
	state map[uint8]struct{}
	lock  sync.Mutex
}

func Start(ctx context.Context, dev device.Device) error {
	state := newState()
	ticker := time.NewTicker(tick)

	dev.NoteOn(func(key uint8) {
		state.Set(key)
		// Reset the ticker so that there is a constant time between the user playing
		// the correct keys and the game recognizing it.
		ticker.Reset(tick)
	})
	dev.NoteOff(func(key uint8) {
		state.Unset(key)
		ticker.Reset(tick)
	})
	if err := dev.Listen(ctx); err != nil {
		return fmt.Errorf("failed to listen on device: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if result := state.Check([]uint8{21, 23}); result {
					fmt.Println("yay")
				}
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}

func newState() *state {
	return &state{
		state: map[uint8]struct{}{},
		lock:  sync.Mutex{},
	}
}

func (s *state) Set(key uint8) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.state[key] = struct{}{}
}

func (s *state) Unset(key uint8) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.state, key)
}

func (s *state) Check(desired []uint8) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.state) != len(desired) {
		return false
	}
	for _, key := range desired {
		if _, ok := s.state[key]; !ok {
			return false
		}
	}
	return true
}
