package simple_buffer

import (
	"fmt"
)

type KeyboardEvent struct {
	Buffer string `json:"buffer"`

	Key   string `json:"key"`
	Alt   bool   `json:"alt"`
	Ctrl  bool   `json:"ctrl"`
	Shift bool   `json:"shift"`
	Meta  bool   `json:"meta"`
}

func (e KeyboardEvent) String() string {
	return fmt.Sprintf("for buffer '%s' key: %s, alt: %t, ctrl: %t, shift: %t",
		e.Buffer, e.Key, e.Alt, e.Ctrl, e.Shift,
	)
}
