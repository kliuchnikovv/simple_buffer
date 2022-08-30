package simple_buffer

type Caret struct {
	Line   int `json:"line"`
	Offset int `json:"offset"`

	lenOfLine     func(int) int
	numberOfLines func() int
}

func NewCaret(line, offset int, lenOfLine func(int) int, numberOfLines func() int) Caret {
	return Caret{
		Line:          line,
		Offset:        offset,
		lenOfLine:     lenOfLine,
		numberOfLines: numberOfLines,
	}
}

func (c Caret) Get() (int, int) {
	var offset = c.Offset
	if c.Offset > c.lenOfLine(c.Line) {
		offset = c.lenOfLine(c.Line)
	}

	return c.Line, offset
}

func (c *Caret) Up() {
	if c.Line > 0 {
		c.Line--
	} else {
		c.Offset = 0
	}
}

func (c *Caret) Down() {
	if c.Line < c.numberOfLines() {
		c.Line++
	} else {
		c.Offset = c.lenOfLine(c.Line)
	}
}

func (c *Caret) Left() {
	if c.Offset > c.lenOfLine(c.Line) {
		c.Offset = c.lenOfLine(c.Line)
	}

	if c.Offset > 0 {
		c.Offset--

		return
	}

	if c.Line != 0 {
		c.Line--
		c.Offset = c.lenOfLine(c.Line)
	}
}

func (c *Caret) Right() {
	if c.Offset > c.lenOfLine(c.Line) {
		c.Offset = c.lenOfLine(c.Line)
	}

	if c.Offset < c.lenOfLine(c.Line) {
		c.Offset++

		return
	}

	if c.Line < c.numberOfLines() {
		c.Line++
		c.Offset = 0
	}
}

func (c *Caret) MoveUp(up int) {
	for i := 0; i < up; i++ {
		c.Up()
	}
}

func (c *Caret) MoveDown(down int) {
	for i := 0; i < down; i++ {
		c.Down()
	}
}

func (c *Caret) MoveLeft(left int) {
	for i := 0; i < left; i++ {
		c.Left()
	}
}

func (c *Caret) MoveRight(right int) {
	for i := 0; i < right; i++ {
		c.Right()
	}
}

func (c *Caret) Set(line, offset int) {
	if line > c.numberOfLines() {
		return
	}
	c.Line = line

	if offset > c.lenOfLine(line) {
		return
	}
	c.Offset = offset
}

func (c *Caret) SetAs(caret Caret) {
	c.Set(caret.Line, caret.Offset)
}

func (c Caret) Equal(second Caret) bool {
	return c.Line == second.Line && c.Offset == second.Offset
}

func (c Caret) Linear() int {
	var linear = 0

	for line := 0; line < c.Line; line++ {
		linear += c.lenOfLine(line) + 1
	}

	_, offset := c.Get()

	return linear + offset
}

type Selection struct {
	start Caret
	end   Caret

	emitter func(string, interface{})
}

func NewSelection(emitter func(string, interface{}), lenOfLine func(int) int, numberOfLines func() int) Selection {
	return Selection{
		emitter: emitter,
		start:   NewCaret(0, 0, lenOfLine, numberOfLines),
		end:     NewCaret(0, 0, lenOfLine, numberOfLines),
	}
}

func (s *Selection) MoveCaret(line, offset int) {
	if line-1 > 0 {
		s.start.MoveDown(line-1)
	} else {
		s.start.MoveUp(line-1)
	}

	if offset > 0 {
		s.start.MoveRight(offset)
	} else {
		s.start.MoveLeft(offset)
	}

	s.Collapse()
}

func (s *Selection) SetSelection(start, end Caret) {
	s.start.SetAs(start)
	s.end.SetAs(end)

	s.emitter("cursor_moved", SelectionChangedEvent{
		Start: s.start,
		End:   s.end,
	})
}

func (s Selection) GetSelection() (Caret, Caret) {
	var (
		startLine, startOffset = s.start.Get()
		endLine, endOffset     = s.end.Get()
	)

	return Caret{
			Line:   startLine,
			Offset: startOffset,
			lenOfLine: s.start.lenOfLine,
			numberOfLines: s.start.numberOfLines,
		}, Caret{
			Line:   endLine,
			Offset: endOffset,
			lenOfLine: s.end.lenOfLine,
			numberOfLines: s.end.numberOfLines,
		}
}

func (s *Selection) CursorUp(isShifted bool) {
	s.start.Up()

	if !isShifted {
		s.Collapse()
	} else {
		s.emitter("cursor_moved", SelectionChangedEvent{
			Start: s.start,
			End:   s.end,
		})
	}

}

func (s *Selection) CursorDown(isShifted bool) {
	if !isShifted {
		s.start.Down()
		s.Collapse()
	} else {
		s.end.Down()
		s.emitter("cursor_moved", SelectionChangedEvent{
			Start: s.start,
			End:   s.end,
		})
	}
}

func (s *Selection) CursorLeft(isShifted bool) {
	s.start.Left()

	if !isShifted {
		s.Collapse()
	} else {
		s.emitter("cursor_moved", SelectionChangedEvent{
			Start: s.start,
			End:   s.end,
		})
	}
}

func (s *Selection) CursorRight(isShifted bool) {
	if !isShifted {
		s.start.Right()
		s.Collapse()
	} else {
		s.end.Right()
		s.emitter("cursor_moved", SelectionChangedEvent{
			Start: s.start,
			End:   s.end,
		})
	}
}

func (s Selection) IsCollapsed() bool {
	return s.start.Equal(s.end)
}

func (s Selection) Linear() (int, int) {
	var (
		start = s.start.Linear()
		end   = s.end.Linear()
	)

	return start, end - start
}

func (s *Selection) Collapse() {
	s.end.SetAs(s.start)

	s.emitter("cursor_moved", SelectionChangedEvent{
		Start: s.start,
		End:   s.end,
	})
}
