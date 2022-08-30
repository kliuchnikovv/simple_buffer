package simple_buffer

import (
	"fmt"
	"strings"

	"github.com/KlyuchnikovV/edicode/core/context"
	"golang.design/x/clipboard"
)

type Buffer struct {
	name string
	data []rune

	ctx context.Context

	KeyEvents   chan KeyboardEvent
	MouseEvents chan MouseEvent
	Selection
}

func New(ctx context.Context, name string, runes ...rune) *Buffer {
	var buffer = &Buffer{
		ctx:         ctx,
		name:        name,
		data:        runes,
		KeyEvents:   make(chan KeyboardEvent, 100),
		MouseEvents: make(chan MouseEvent, 100),
	}

	buffer.Selection = NewSelection(func(s string, i interface{}) {
		buffer.ctx.Emit(s, name, i)
	}, func(i int) int {
		l, err := buffer.LengthOfLine(i)
		if err != nil {
			panic(err)
		}

		return l
	}, func() int {
		return numberOfLines(buffer.data...)
	})

	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	go buffer.listenEvents()

	return buffer
}

func NewFromBytes(ctx context.Context, name string, bytes ...byte) *Buffer {
	return New(ctx, name, []rune(string(bytes))...)
}

func (buffer *Buffer) String() string {
	return string(buffer.data)
}

func (buffer *Buffer) Append(runes ...rune) error {
	if len(runes) == 0 {
		return nil
	}

	defer func() {
		var (
			idx    = strings.LastIndex(string(runes), "\n")
			offset = len(runes)
		)
		if idx != -1 {
			offset = len(runes[idx+1:])
		}

		buffer.Selection.MoveCaret(numberOfLines(runes...), offset)
	}()

	return buffer.Insert(buffer.Selection, runes...)
}

func (buffer *Buffer) Insert(selection Selection, runes ...rune) error {
	if !selection.IsCollapsed() {
		if err := buffer.Remove(selection, false); err != nil {
			return err
		}
	}

	if len(runes) == 0 {
		return nil
	}

	var offset, _ = selection.Linear()

	var data = make([]rune, 0, len(buffer.data)+len(runes))

	data = append(data, buffer.data[:offset]...)
	data = append(data, runes...)
	data = append(data, buffer.data[offset:]...)

	buffer.data = data
	buffer.ctx.Emit("buffer", "changed", buffer.name)

	return nil
}

func (buffer *Buffer) Remove(selection Selection, emitEvent bool) error {
	if selection.IsCollapsed() {
		return nil
	}

	var offset, length = selection.Linear()

	var data = make([]rune, 0, len(buffer.data)-length)

	data = append(data, buffer.data[:offset]...)
	data = append(data, buffer.data[offset+length:]...)

	buffer.data = data
	if emitEvent {
		buffer.ctx.Emit("buffer", "changed", buffer.name)
	}

	buffer.Selection.Collapse()

	return nil
}

func (buffer *Buffer) Delete() error {
	if buffer.Selection.IsCollapsed() {
		buffer.Selection.start.MoveLeft(1)
	}

	defer buffer.Selection.Collapse()

	return buffer.Remove(buffer.Selection, true)
}

func (buffer *Buffer) LengthOfLine(line int) (int, error) {
	var lines = strings.Split(buffer.String(), "\n")

	if line < 0 || line > len(lines) {
		return -1, fmt.Errorf("line not in range [%d : %d)", 0, len(lines))
	}

	return len(lines[line]), nil
}

func (buffer *Buffer) GetSelectedText() string {
	start, end := buffer.GetSelection()

	return string(buffer.data[start.Linear():end.Linear()])
}

func (buffer *Buffer) Copy() {
	clipboard.Write(clipboard.FmtText, []byte(buffer.GetSelectedText()))
}

func (buffer *Buffer) Paste() error {
	if !buffer.Selection.IsCollapsed() {
		if err := buffer.Remove(buffer.Selection, false); err != nil {
			return err
		}
	}

	return buffer.Append([]rune(string(clipboard.Read(clipboard.FmtText)))...)
}

func (buffer *Buffer) Cut() error {
	buffer.Copy()
	return buffer.Remove(buffer.Selection, true)
}

func (buffer *Buffer) SelectAll() error {
	var (
		numberOfLines     = numberOfLines(buffer.data...)
		lengthOfLine, err = buffer.LengthOfLine(numberOfLines - 1)
	)

	if err != nil {
		return err
	}

	buffer.SetSelection(Caret{
		Line: 0, Offset: 0,
	}, Caret{
		Line:   numberOfLines - 1,
		Offset: lengthOfLine,
	})

	return nil
}
