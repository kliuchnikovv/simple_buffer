package selection

import "fmt"

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
