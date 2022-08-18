package sbuffer

import "fmt"

type Buffer struct {
	data []rune
}

func New(runes ...rune) *Buffer {
	return &Buffer{
		data: runes,
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

func (buffer *Buffer) Delete(index, symbols int) error {
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
