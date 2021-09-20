// Package game stores game state.
package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/inahga/aightreader/device"
	"github.com/inahga/aightreader/ui"
)

// Being tick based avoids the condition where user presses a bunch of keys at
// the same time, and the incoming MIDI messages coincidentally arrive in the
// correct order.
const tick = 200 * time.Millisecond

type state struct {
	u         ui.UI
	dev       device.Device
	ticker    *time.Ticker
	state     map[uint8]struct{}
	stateLock sync.Mutex
}

func Start(ctx context.Context, u ui.UI) error {
	state := &state{
		u:         u,
		ticker:    time.NewTicker(tick),
		state:     map[uint8]struct{}{},
		stateLock: sync.Mutex{},
	}
	state.chooseDevice()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-state.ticker.C:
				if result := state.Check([]uint8{21, 23}); result {
					fmt.Println("yay")
				}
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}

func (s *state) Set(key uint8) {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()
	s.state[key] = struct{}{}
}

func (s *state) Unset(key uint8) {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()
	delete(s.state, key)
}

func (s *state) Check(desired []uint8) bool {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()
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

func (s *state) setDevice(d device.Device) error {
	if s.dev != nil {
		if err := s.dev.Close(); err != nil {
			return err
		}
	}

	d.NoteOn(func(key uint8) {
		s.Set(key)
		// Reset the ticker so that there is a constant time between the user playing
		// the correct keys and the game recognizing it.
		s.ticker.Reset(tick)
	})
	d.NoteOff(func(key uint8) {
		s.Unset(key)
		s.ticker.Reset(tick)
	})
	if err := d.Listen(); err != nil {
		return fmt.Errorf("failed to listen on device: %w", err)
	}

	return nil
}

func (s *state) chooseDevice() {
	// TODO: additional edge cases to consider:
	//   - User plugs in a keyboard on this screen
	//   - Only one device is available, should automatically choose that one
	//     - But how to surface errors if there is one?
	//   - User wants to remember device choice from the last session
	for {
		dev := s.u.ChooseDevice()
		if err := s.setDevice(dev); err != nil {
			s.u.ChooseDeviceError(err)
		} else {
			break
		}
	}
}
