package simple_buffer

import (
	"fmt"
	"strings"
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

type MouseEvent struct {
	Buffer string `json:"buffer"`

	Start Caret `json:"start"`
	End   Caret `json:"end"`
}

func (e MouseEvent) String() string {
	return fmt.Sprintf("for buffer '%s' start {%d-%d}, end {%d-%d}",
		e.Buffer, e.Start.Line, e.Start.Offset, e.End.Line, e.End.Offset,
	)
}

type SelectionChangedEvent struct {
	Start Caret `json:"start"`
	End   Caret `json:"end"`
}

func numberOfLines(runes ...rune) int {
	return strings.Count(string(runes), "\n")+1
}
