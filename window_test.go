package ais

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("errorReader used for testing")
}

func TestNewWindow(t *testing.T) {
	testSet, _ := OpenRecordSet("testdata/ten.csv")
	defer testSet.Close()
	badSet := &RecordSet{h: badHeaders}
	badSet2 := &RecordSet{
		h: goodHeaders,
		r: csv.NewReader(&errorReader{}),
	}
	badSet3, _ := OpenRecordSet("testdata/bad.csv")
	defer badSet3.Close()

	type args struct {
		rs    *RecordSet
		width time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Window
		wantErr bool
	}{
		{
			name: "window from ten.csv",
			args: args{
				rs:    testSet,
				width: 5 * time.Second,
			},
			want:    &testWindow,
			wantErr: false,
		},
		{
			name: "bad headers should trigger error",
			args: args{
				rs:    badSet,
				width: 5 * time.Second,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "readFirst error",
			args: args{
				rs:    badSet2,
				width: 5 * time.Second,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "time parse error",
			args: args{
				rs:    badSet3,
				width: 5 * time.Second,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWindow(tt.args.rs, tt.args.width)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWindow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_Left(t *testing.T) {

	tests := []struct {
		name string
		win  *Window
		want time.Time
	}{
		{
			name: "from testWindow",
			win:  &testWindow,
			want: time.Date(2017, time.December, 1, 00, 00, 01, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.win.Left(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Window.Left() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_SetLeft(t *testing.T) {

	type args struct {
		marker time.Time
	}
	tests := []struct {
		name string
		win  *Window
		args args
		want time.Time
	}{
		{
			name: "from testWindow",
			win:  &testWindow,
			args: args{
				marker: time.Date(2016, time.November, 30, 00, 00, 02, 0, time.UTC),
			},
			want: time.Date(2016, time.November, 30, 00, 00, 02, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			win := tt.win
			win.SetLeft(tt.args.marker)
			got := win.Left()
			if got != tt.want {
				t.Errorf("Window.SetLeft() where Left()= %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_Right(t *testing.T) {
	tests := []struct {
		name string
		win  *Window
		want time.Time
	}{
		{
			name: "test window right",
			win:  &testWindow,
			want: time.Date(2017, time.December, 1, 00, 00, 06, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.win.Right(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Window.Right() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow_AddRecord(t *testing.T) {

	type args struct {
		rec Record
	}
	tests := []struct {
		name string
		win  *Window
		args args
		want int
	}{
		{
			name: "add to empty window",
			win:  &testWindow,
			args: args{testRec0},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.win.AddRecord(tt.args.rec)
			got := len(tt.win.Data)
			if got != tt.want {
				t.Errorf("Window.AddRecord() error: got len=%v, want len=%v", got, tt.want)
			}
		})
	}
}
