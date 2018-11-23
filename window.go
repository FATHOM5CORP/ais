package ais

import (
	"bytes"
	"fmt"
	"time"
)

// Window is used to create a convolution algorithm that slides down a RecordSet
// and performs analysis on Records that are within the a time window.
type Window struct {
	leftMarker, rightMarker time.Time
	timeIndex               int
	width                   time.Duration
	Data                    map[uint64]*Record
}

// NewWindow returns a *Window with the left marker set to the time in
// the next record read from the RecordSet. The Window width will be set from
// the argument provided and the righ marker will be derived from left and width.
// When creating a Window right after opening a RecordSet then the Window
// will be set to first Record in the set, but that first record will still be
// available to the client's first call to rs.Read(). For any non-nil error
// NewWindow returns nil and the error.
func NewWindow(rs *RecordSet, width time.Duration) (*Window, error) {
	win := new(Window)
	timeIndex, ok := rs.Headers().Contains("BaseDateTime")
	if !ok {
		return nil, fmt.Errorf("newwindow: headers does not contain BaseDateTime")
	}
	win.SetIndex(timeIndex)
	rec, err := rs.readFirst()
	if err != nil {
		return nil, fmt.Errorf("newwindow: %v", err)
	}
	t, err := time.Parse(TimeLayout, (*rec)[timeIndex])
	if err != nil {
		return nil, fmt.Errorf("newwindow: %v", err)
	}
	win.SetLeft(t)
	win.SetWidth(width)
	win.SetRight(win.Left().Add(win.Width()))
	return win, nil
}

// Left returns the left marker.
func (win *Window) Left() time.Time { return win.leftMarker }

// SetLeft defines the left marker for the Window
func (win *Window) SetLeft(marker time.Time) {
	win.leftMarker = marker
}

// Right returns the right marker.
func (win *Window) Right() time.Time { return win.rightMarker }

// SetRight defines the right marker of the Window.
func (win *Window) SetRight(marker time.Time) {
	win.rightMarker = marker
}

// Width returns the width of the Window.
func (win *Window) Width() time.Duration { return win.width }

// SetWidth provides the block of time coverd by the Window.
func (win *Window) SetWidth(dur time.Duration) {
	win.width = dur
}

// SetIndex provides the integer index of the BaseDateTime field the
// Records stored in the Window.
func (win *Window) SetIndex(index int) {
	win.timeIndex = index
}

// AddRecord appends a new Record to the data in the Window.
func (win *Window) AddRecord(rec Record) {
	if win.Data == nil {
		win.Data = make(map[uint64]*Record)
	}
	h := rec.Hash()
	// fmt.Println("hash is: ", h)
	win.Data[h] = &rec
}

// InWindow tests if a time is in the Window.
func (win *Window) InWindow(t time.Time) bool {
	if win.leftMarker.Equal(t) {
		return true
	}
	return win.leftMarker.Before(t) && t.Before(win.rightMarker)
}

// RecordInWindow returns true if the record is in the Window.
// Errors are possible from parsing the BaseDateTime field of the
// Record.
func (win *Window) RecordInWindow(rec *Record) (bool, error) {
	t, err := rec.ParseTime(win.timeIndex)
	if err != nil {
		return false, fmt.Errorf("recordinwindow: %v", err)
	}
	return win.InWindow(t), nil
}

// Slide moves the window down by the time provided in the arugment dur.
// Slide also removes any data from the Window that would no longer return
// true from InWindow for the new left and right markers after the Slide.
func (win *Window) Slide(dur time.Duration) {
	win.SetLeft(win.leftMarker.Add(dur))
	win.SetRight(win.leftMarker.Add(win.Width()))

	win.validate()
}

// Len returns the lenght of the slice holding the Records in the Window
func (win *Window) Len() int {
	return len(win.Data)
}

// Validate checks the data held by the Window to ensure every Record passes
// the InWindow test.
func (win *Window) validate() error {
	for hash, rec := range win.Data {
		in, err := win.RecordInWindow(rec)
		if err != nil {
			return err
		}
		if !in {
			delete(win.Data, hash)
		}
	}
	return nil
}

// String implements the Stringer interface for Window.
func (win Window) String() string {
	var buf bytes.Buffer
	for _, rec := range win.Data {
		fmt.Fprintln(&buf, rec)
	}
	return buf.String()
}
