package selection

type Caret struct {
	Line   int `json:"line"`
	Offset int `json:"offset"`

	getLine       func(int) []rune
	numberOfLines func() int
}

func NewCaret(line, offset int, getLine func(int) []rune, numberOfLines func() int) Caret {
	return Caret{
		Line:          line,
		Offset:        offset,
		getLine:       getLine,
		numberOfLines: numberOfLines,
	}
}

func (c Caret) Get() (int, int) {
	var offset = c.Offset
	if c.Offset > len(c.getLine(c.Line)) {
		offset = len(c.getLine(c.Line))
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
		c.Offset = len(c.getLine(c.Line))
	}
}

func (c *Caret) Left() {
	if c.Offset > len(c.getLine(c.Line)) {
		c.Offset = len(c.getLine(c.Line))
	}

	if c.Offset > 0 {
		c.Offset--

		return
	}

	if c.Line != 0 {
		c.Line--
		c.Offset = len(c.getLine(c.Line))
	}
}

func (c *Caret) Right() {
	if c.Offset > len(c.getLine(c.Line)) {
		c.Offset = len(c.getLine(c.Line))
	}

	if c.Offset < len(c.getLine(c.Line)) {
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

	if offset > len(c.getLine(line)) {
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
		linear += len(c.getLine(line)) + 1
	}

	_, offset := c.Get()

	return linear + offset
}

type Selection struct {
	start Caret
	end   Caret

	emitter func(string, interface{})
}

func NewSelection(emitter func(string, interface{}), lenOfLine func(int) []rune, numberOfLines func() int) Selection {
	return Selection{
		emitter: emitter,
		start:   NewCaret(0, 0, lenOfLine, numberOfLines),
		end:     NewCaret(0, 0, lenOfLine, numberOfLines),
	}
}

func (s Selection) Start() Caret {
	return s.start
}

func (s Selection) End() Caret {
	return s.end
}

func (s *Selection) MoveCaret(line, offset int) {
	if line != 0 {
		if line-1 > 0 {
			s.start.MoveDown(line)
		} else {
			s.start.MoveUp(-line)
		}
	}

	if offset != 0 {
		if offset > 0 {
			s.start.MoveRight(offset)
		} else {
			s.start.MoveLeft(-offset)
		}
	}
}

func (s *Selection) SetSelection(start, end Caret) {
	s.start.SetAs(start)
	s.end.SetAs(end)

	if s.end.Line < s.start.Line || s.end.Offset < s.start.Offset {
		s.start.Line, s.end.Line = s.end.Line, s.start.Line
		s.start.Offset, s.end.Offset = s.end.Offset, s.start.Offset
	}

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
			Line:          startLine,
			Offset:        startOffset,
			getLine:       s.start.getLine,
			numberOfLines: s.start.numberOfLines,
		}, Caret{
			Line:          endLine,
			Offset:        endOffset,
			getLine:       s.end.getLine,
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

func (s *Selection) SetStart(line, offset int) {
	s.start.Set(line, offset)
	s.Collapse()
}

func (s *Selection) Swap() {
	s.start, s.end = s.end, s.start
}

func (s *Selection) SetEnd(line, offset int) {
	s.end.Set(line, offset)

	if s.end.Line < s.start.Line || (s.end.Line == s.start.Line && s.end.Offset < s.start.Offset) {
		s.Swap()
	}

	s.emitter("cursor_moved", SelectionChangedEvent{
		Start: s.start,
		End:   s.end,
	})
}

func (s Selection) Copy() Selection {
	return Selection{
		start:   s.start,
		end:     s.end,
		emitter: s.emitter,
	}
}
