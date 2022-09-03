package stack

import (
	"github.com/KlyuchnikovV/simple_buffer/selection"
)

type Event interface {
	Undo(*ActionStack) error
	Redo(*ActionStack) error
}

type ActionStack struct {
	stack   []Event
	pointer int

	insert func(selection.Selection, ...rune) error
	delete func(selection.Selection) error
}

func New(capacity int,
	insert func(selection.Selection, ...rune) error,
	delete func(selection.Selection) error,
) *ActionStack {
	return &ActionStack{
		stack:  make([]Event, 0, capacity),
		insert: insert,
		delete: delete,
	}
}

func (as *ActionStack) Push(event Event) {
	if len(as.stack) != as.pointer {
		as.stack = as.stack[:as.pointer]
	}

	as.stack = append(as.stack, event)
	as.pointer++
}

func (as *ActionStack) Undo() {
	if as.pointer == 0 {
		return
	}

	defer func() {
		as.pointer--
	}()

	if err := as.stack[as.pointer-1].Undo(as); err != nil {
		panic(err)
	}
}

func (as *ActionStack) Redo() {
	if len(as.stack) == as.pointer {
		return
	}

	defer func() {
		as.pointer++
	}()

	if err := as.stack[as.pointer].Redo(as); err != nil {
		panic(err)
	}
}

type Insertion struct {
	selection.Selection
	data []rune
}

func NewInsertion(sel selection.Selection, runes []rune) *Insertion {
	return &Insertion{
		Selection: sel,
		data:      runes,
	}
}

func (i *Insertion) Undo(as *ActionStack) error {
	return as.delete(i.Selection)
}

func (i *Insertion) Redo(as *ActionStack) error {
	return as.insert(i.Selection, i.data...)
}

type Deletion struct {
	selection.Selection
	data []rune
}

func NewDeletion(sel selection.Selection, runes []rune) *Deletion {
	return &Deletion{
		Selection: sel,
		data:      runes,
	}
}

func (d *Deletion) Undo(as *ActionStack) error {
	var newSelection = d.Selection.Copy()
	newSelection.Collapse()

	return as.insert(newSelection, d.data...)
}

func (d *Deletion) Redo(as *ActionStack) error {
	return as.delete(d.Selection)
}

type Replacing struct {
	Insertion
	Deletion
}

func NewReplacing(ins Insertion, del Deletion) *Replacing {
	return &Replacing{
		Insertion: ins,
		Deletion:  del,
	}
}

func (r *Replacing) Undo(as *ActionStack) error {
	if err := r.Insertion.Undo(as); err != nil {
		return err
	}

	return r.Deletion.Undo(as)
}

func (r *Replacing) Redo(as *ActionStack) error {
	if err := r.Deletion.Redo(as); err != nil {
		return err
	}

	return r.Insertion.Redo(as)
}
