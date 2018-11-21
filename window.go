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
	data                    map[uint64]*Record
}

// Left returns the left marker.
func (win *Window) Left() time.Time { return win.leftMarker }

// SetLeft defines the left marker for the Window
func (win *Window) SetLeft(marker time.Time) {
	win.leftMarker = marker
	if win.leftMarker.Add(win.width).After(win.leftMarker) {
		win.rightMarker = win.leftMarker.Add(win.width)
	}
}

// Right returns the right marker.
func (win *Window) Right() time.Time { return win.rightMarker }

// SetRight defines the right marker of the Window.
func (win *Window) SetRight(marker time.Time) {
	win.rightMarker = marker
}

// SetIndex provides the integer index of the BaseDateTime field the
// Records stored in the Window.
func (win *Window) SetIndex(index int) {
	win.timeIndex = index
}

// SetWidth provides the block of time coverd by the Window.
func (win *Window) SetWidth(dur time.Duration) {
	win.width = dur
}

// AddRecord appends a new Record to the data in the Window.
func (win *Window) AddRecord(rec Record) {
	if win.data == nil {
		win.data = make(map[uint64]*Record)
	}
	win.data[rec.Hash()] = &rec
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
// Slide also removes any data from the Window that is would no longer return
// true from InWindow for the new left and right markers after the Slide.
func (win *Window) Slide(dur time.Duration) {
	win.SetLeft(win.leftMarker.Add(dur))
	win.SetRight(win.leftMarker.Add(win.width))

	win.validate()
}

// Len returns the lenght of the slice holding the Records in the Window
func (win *Window) Len() int {
	return len(win.data)
}

// String implements the Stringer interface for Window.
func (win Window) String() string {
	var buf bytes.Buffer
	for _, rec := range win.data {
		fmt.Fprintln(&buf, rec)
	}
	return buf.String()
}

// Validate checks the data held by the Window to ensure every Record passes
// the InWindow test.
func (win *Window) validate() error {
	for _, rec := range win.data {
		in, err := win.RecordInWindow(rec)
		if err != nil {
			return err
		}
		if !in {
			// if _, ok := win.data[rec.Hash()]; ok {
			// 	fmt.Println("deleting:", rec)
			// }
			delete(win.data, rec.Hash())
		}
	}
	return nil
}
