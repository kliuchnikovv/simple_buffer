package simple_buffer

import (
	"fmt"
	"strings"
)

type Buffer struct {
	data []rune

	line   int
	offset int
	cursor int
}

func New(runes ...rune) *Buffer {
	return &Buffer{
		data: runes,
	}
}

func NewFromBytes(bytes ...byte) *Buffer {
	return &Buffer{
		data: []rune(string(bytes)),
	}
}

func (buffer *Buffer) String() string {
	return string(buffer.data)
}

func (buffer *Buffer) Insert(index int, runes ...rune) error {
	if err := buffer.checkIndex(index); err != nil {
		return err
	}

	if len(runes) == 0 {
		return nil
	}

	var data = make([]rune, 0, len(buffer.data)+len(runes))

	data = append(data, buffer.data[:index]...)
	data = append(data, runes...)
	data = append(data, buffer.data[index:]...)

	buffer.data = data

	return nil
}

func (buffer *Buffer) Append(runes ...rune) error {
	if len(runes) == 0 {
		return nil
	}

	defer func() {
		buffer.cursor += len(runes)
		buffer.restoreCurrentLine()
	}()

	return buffer.Insert(buffer.cursor, runes...)
}

func (buffer *Buffer) Remove(index, symbols int) error {
	if err := buffer.checkIndex(index); err != nil {
		return err
	}

	if err := buffer.checkIndex(index + symbols); err != nil {
		return err
	}

	var data = make([]rune, 0, len(buffer.data)-symbols)

	data = append(data, buffer.data[:index]...)
	data = append(data, buffer.data[index+symbols:]...)

	buffer.data = data

	return nil
}

func (buffer *Buffer) Delete(symbols int) error {
	if err := buffer.checkIndex(buffer.cursor - symbols); err != nil {
		return err
	}

	buffer.cursor -= symbols
	defer buffer.restoreCurrentLine()

	return buffer.Remove(buffer.cursor, symbols)
}

func (buffer *Buffer) LengthOfLine(line int) (int, error) {
	var lines = strings.Split(buffer.String(), "\n")

	if line < 0 || line > len(lines) {
		return -1, fmt.Errorf("line not in range [%d : %d)", 0, len(lines))
	}

	return len(lines[line]), nil
}

func (buffer *Buffer) GetRange(index, symbols int) (string, error) {
	if err := buffer.checkIndex(index); err != nil {
		return "", err
	}

	if err := buffer.checkIndex(index + symbols); err != nil {
		return "", err
	}

	return string(buffer.data[index : index+symbols]), nil
}

func (buffer Buffer) checkIndex(index int) error {
	if index < 0 || index >= len(buffer.data) {
		return fmt.Errorf("index not in range [%d : %d)", 0, len(buffer.data))
	}

	return nil
}

func (buffer *Buffer) CursorUp() {
	var lines = strings.Split(buffer.String(), "\n")

	if buffer.line == 0 {
		buffer.cursor = 0
		return
	}

	buffer.line--

	buffer.cursor -= buffer.offset

	if len(lines[buffer.line]) < buffer.offset {
		buffer.offset = len(lines[buffer.line])
	}

	buffer.cursor -= len(lines[buffer.line]) - buffer.offset + 1

	buffer.restoreCurrentLine()
}

func (buffer *Buffer) CursorDown() {
	var lines = strings.Split(buffer.String(), "\n")

	if buffer.line >= len(lines)-1 {
		buffer.cursor = len(buffer.data) - 1
		return
	}

	buffer.line++

	buffer.cursor += len(lines[buffer.line-1]) - buffer.offset + 1

	if len(lines[buffer.line]) < buffer.offset {
		buffer.offset = len(lines[buffer.line])
	}

	buffer.cursor += buffer.offset

	buffer.restoreCurrentLine()
}

func (buffer *Buffer) CursorLeft() {
	if buffer.cursor > 0 {
		buffer.cursor--
	}
	buffer.restoreCurrentLine()
}

func (buffer *Buffer) CursorRight() {
	if buffer.cursor < len(buffer.data) {
		buffer.cursor++
	}
	buffer.restoreCurrentLine()
}

func (buffer *Buffer) restoreCurrentLine() {
	var (
		line   = 0
		offset = 0
		cursor = buffer.cursor
	)

	for _, symbol := range buffer.data {
		if cursor == 0 {
			break
		}

		if symbol == '\n' {
			line++
			offset = -1
		}

		cursor--
		offset++
	}

	buffer.line = line
	buffer.offset = offset
}

func (buffer *Buffer) restoreCursor() {
	var (
		line   = buffer.line
		offset = buffer.offset
		cursor = 0
	)

	for _, symbol := range buffer.data {
		if line == 0 {
			if offset == 0 {
				break
			}

			offset--
		}

		if symbol == '\n' {
			line--
		}

		cursor++
	}

	buffer.cursor = cursor
}

func (buffer *Buffer) GetCursor() (int, int, int) {
	return buffer.line, buffer.offset, buffer.cursor
}

func (buffer *Buffer) SetCursor(line, offset int) {
	buffer.line = line

	if offset == -1 {
		offset, _ = buffer.LengthOfLine(line)
	}
	buffer.offset = offset

	buffer.restoreCursor()
}
