package simple_buffer

import (
	"fmt"
	"strings"
	"time"

	"github.com/KlyuchnikovV/edicode/core/context"
	"github.com/KlyuchnikovV/simple_buffer/selection"
	"github.com/KlyuchnikovV/simple_buffer/stack"
	"golang.design/x/clipboard"
)

type Buffer struct {
	name       string
	data       []rune
	ModifiedAt int64

	ctx context.Context

	KeyEvents   chan KeyboardEvent
	MouseEvents chan selection.MouseEvent
	selection.Selection

	stack *stack.ActionStack
}

func New(ctx context.Context, name string, runes ...rune) *Buffer {
	var buffer = &Buffer{
		ctx:         ctx,
		name:        name,
		data:        runes,
		KeyEvents:   make(chan KeyboardEvent, 100),
		MouseEvents: make(chan selection.MouseEvent, 100),
		ModifiedAt:  time.Now().Unix(),
	}

	buffer.stack = stack.New(50, buffer.insert, buffer.remove, buffer.SetSelection)

	buffer.Selection = selection.NewSelection(func(s string, i interface{}) {
		buffer.ctx.Emit(s, name, i)
	}, buffer.GetLine, func() int {
		return strings.Count(string(runes), "\n") + 1
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

	// defer func() {
	// 	var (
	// 		idx    = strings.LastIndex(string(runes), "\n")
	// 		offset = len(runes)
	// 	)
	// 	if idx != -1 {
	// 		offset = len(runes[idx+1:])
	// 	}

	// 	buffer.Selection.MoveCaret(numberOfLines(runes...), offset)
	// 	buffer.Collapse()
	// }()

	return buffer.Insert(buffer.Selection, runes...)
}

func (buffer *Buffer) Insert(sel selection.Selection, runes ...rune) error {
	if len(runes) == 0 {
		return nil
	}

	// var (
	// 	newSelection = sel.Copy()
	// 	end          = buffer.Selection.End()
	// )

	// end.Offset += len(runes)
	// newSelection.SetSelection(sel.Start(), end)

	// var event stack.Event = stack.NewInsertion(newSelection, runes)
	// if !sel.IsCollapsed() {
	// 	event = stack.NewReplacing(
	// 		*event.(*stack.Insertion),
	// 		*stack.NewDeletion(sel.Copy(), buffer.GetBySelection(sel)),
	// 	)
	// }

	defer func(old []rune) {
		var newSelection = sel.Copy()
		newSelection.SetStartAndEnd(sel.Start(), buffer.Selection.End())

		var event stack.Event = stack.NewInsertion(newSelection, runes)
		if !sel.IsCollapsed() {
			event = stack.NewReplacing(
				*event.(*stack.Insertion),
				*stack.NewDeletion(sel, old),
			)
		}

		buffer.stack.Push(event)
	}(buffer.GetBySelection(sel))

	return buffer.insert(sel, runes...)
}

func (buffer *Buffer) insert(sel selection.Selection, runes ...rune) error {
	if len(runes) == 0 {
		return nil
	}

	if !sel.IsCollapsed() {
		if err := buffer.remove(sel); err != nil {
			return err
		}
	}

	var offset, _ = sel.Linear()

	var data = make([]rune, 0, len(buffer.data)+len(runes))

	data = append(data, buffer.data[:offset]...)
	data = append(data, runes...)
	data = append(data, buffer.data[offset:]...)

	buffer.data = data

	var (
		idx           = strings.LastIndex(string(runes), "\n")
		numberOfLines = strings.Count(string(runes), "\n")
	)

	offset = len(runes)
	if idx != -1 {
		offset = len(runes[idx+1:])
	}

	buffer.Selection.MoveCaret(numberOfLines, offset /*-buffer.Start().Offset*/)
	buffer.Collapse()

	buffer.ctx.Emit("buffer", "changed", buffer.name)

	return nil
}

func (buffer *Buffer) Delete() error {
	if buffer.Selection.IsCollapsed() {
		buffer.Selection.MoveCaret(0, -1)
	}

	return buffer.Remove(buffer.Selection)
}

func (buffer *Buffer) Remove(sel selection.Selection) error {
	if sel.IsCollapsed() {
		return nil
	}

	defer func(s selection.Selection, runes []rune) {
		buffer.stack.Push(stack.NewDeletion(s, runes))
	}(sel.Copy(), buffer.GetBySelection(sel))

	return buffer.remove(sel)
}

func (buffer *Buffer) remove(sel selection.Selection) error {
	if sel.IsCollapsed() {
		return nil
	}

	var offset, length = sel.Linear()

	var data = make([]rune, 0, len(buffer.data)-length)

	data = append(data, buffer.data[:offset]...)
	data = append(data, buffer.data[offset+length:]...)

	buffer.data = data
	sel.Collapse()
	buffer.Selection.SetStartAndEnd(sel.Start(), sel.End())
	buffer.ctx.Emit("buffer", "changed", buffer.name)

	return nil
}

func (buffer *Buffer) LengthOfLine(line int) (int, error) {
	var lines = strings.Split(buffer.String(), "\n")

	if line < 0 || line > len(lines) {
		return -1, fmt.Errorf("line not in range [%d : %d)", 0, len(lines))
	}

	return len(lines[line]), nil
}

func (buffer *Buffer) GetLine(line int) []rune {
	var lines = strings.Split(buffer.String(), "\n")

	return []rune(lines[line])
}

func (buffer *Buffer) GetSelectedText() string {
	return string(buffer.GetBySelection(buffer.Selection))
}

func (buffer *Buffer) GetBySelection(sel selection.Selection) []rune {
	var start, end = sel.GetSelection()

	return buffer.data[start.Linear():end.Linear()]
}

func (buffer *Buffer) Copy() {
	clipboard.Write(clipboard.FmtText, []byte(buffer.GetSelectedText()))
}

func (buffer *Buffer) Paste() error {
	return buffer.Append([]rune(string(clipboard.Read(clipboard.FmtText)))...)
}

func (buffer *Buffer) Cut() error {
	buffer.Copy()
	return buffer.Delete()
}

func (buffer *Buffer) SelectAll() error {
	var (
		numberOfLines     = strings.Count(string(buffer.data), "\n")
		lengthOfLine, err = buffer.LengthOfLine(numberOfLines)
	)

	if err != nil {
		return err
	}

	buffer.SetStartAndEnd(selection.Caret{
		Line: 0, Offset: 0,
	}, selection.Caret{
		Line:   numberOfLines,
		Offset: lengthOfLine,
	})

	return nil
}

func (buffer *Buffer) Name() string {
	return buffer.name
}

func (buffer *Buffer) Undo() {
	buffer.stack.Undo()
}

func (buffer *Buffer) Redo() {
	buffer.stack.Redo()
}
