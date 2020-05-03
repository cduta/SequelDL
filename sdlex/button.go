package sdlex

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
  BUTTON_PRESSED, BUTTON_RELEASED = 1,0
)

func IsRepeatButtonPress(event *sdl.KeyboardEvent) bool {
	return event.Repeat != 0
}