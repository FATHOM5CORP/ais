package ais

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// TestString is a string that serves as a proxy for a csv file and is used in
// a RecordSet test by intantiating the csv.Reader as csv.NewReader(strings.NewReader(testsTring))
var testString = `# test data corresponding to the first few rows of ten.csv
MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo
477307901,2017-12-01T00:00:01,31.90512,-76.32652,0.0,131.0,352.0,FIRST,IMO9739666,VRPJ6,1004,moored,337,,,
338029922,2017-12-01T00:00:02,42.83931,-73.74403,37.7,110.6,511.0,SECOND,,,,,,,,
369080003,2017-12-01T00:00:03,43.60792,-74.20417,4.1,1.0,5.0,THIRD,IMO9795933,WDI7248,1025,under way using engine,,,,
`

// TestStringBadHeader1 lacks the canonical capitalization for MMSI
var testStringBadHeader1 = `# test data corresponding to the first few rows of ten.csv
mmsi,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo
477307901,2017-12-01T00:00:01,31.90512,-76.32652,0.0,131.0,352.0,FIRST,IMO9739666,VRPJ6,1004,moored,337,,,
338029922,2017-12-01T00:00:02,42.83931,-73.74403,37.7,110.6,511.0,SECOND,,,,,,,,
369080003,2017-12-01T00:00:03,43.60792,-74.20417,4.1,1.0,5.0,THIRD,IMO9795933,WDI7248,1025,under way using engine,,,,
`

// TestStringBadHeader2 is missing the VesselName field
var testStringBadHeader2 = `# test data corresponding to the first few rows of ten.csv
MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo
477307901,2017-12-01T00:00:01,31.90512,-76.32652,0.0,131.0,352.0,IMO9739666,VRPJ6,1004,moored,337,,,
338029922,2017-12-01T00:00:02,42.83931,-73.74403,37.7,110.6,511.0,,,,,,,,
369080003,2017-12-01T00:00:03,43.60792,-74.20417,4.1,1.0,5.0,IMO9795933,WDI7248,1025,under way using engine,,,,
`

var testRec0 = Record{"376494000", "2017-12-01T00:00:00", "30.28963", "-110.73522", "9.4", "158.2", "511.0"}
var testRec1 = Record{"376494001", "2017-12-01T00:00:01", "31.28963", "-111.73522", "9.4", "158.2", "511.0"}
var testRec2 = Record{"376494002", "2017-12-01T00:00:02", "32.28963", "-112.73522", "9.4", "158.2", "511.0"}

// These records are the string slices from the ten.csv file in testdata.
var firstRec = []string{"477307901", "2017-12-01T00:00:01", "31.90512", "-76.32652", "0.0", "131.0", "352.0", "FIRST", "IMO9739666", "VRPJ6", "1004", "moored", "337", "", "", ""}

// These records are the three track records in track.csv
var track1 = []string{"477307901", "2017-12-01T00:00:01", "31.90512", "-76.32652", "0.0", "131.0", "352.0", "FIRST", "IMO9739666", "VRPJ6", "1004", "underway using engines", "337", "", "", ""}
var track2 = []string{"477307901", "2017-12-01T00:01:01", "31.80512", "-76.42652", "0.0", "131.0", "352.0", "FIRST", "IMO9739666", "VRPJ6", "1004", "underway using engines", "337", "", "", ""}
var track3 = []string{"477307901", "2017-12-01T00:02:01", "31.70512", "-76.52652", "0.0", "131.0", "352.0", "FIRST", "IMO9739666", "VRPJ6", "1004", "underway using engines", "337", "", "", ""}

// Bad record with malformed time
var testRec4 = Record{"376494002", "2017-12-01Txx:00:02", "32.28963", "-112.73522", "9.4", "158.2", "511.0"}

const testRec0String = `376494000 2017-12-01T00:00:00 30.28963 -110.73522 9.4 158.2 511.0]`

var testWindow = Window{
	leftMarker:  time.Date(2017, time.December, 1, 00, 00, 01, 0, time.UTC),
	width:       5 * time.Second,
	rightMarker: time.Date(2017, time.December, 1, 00, 00, 01, 0, time.UTC).Add(5 * time.Second),
	timeIndex:   1,
	Data:        nil,
}

var badHeaders = Headers{ // Missing canonical name BaseDateTime
	Fields: strings.Split("MMSI,Timestamp,LAT,LON,SOG,COG,Heading,"+
		"VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ","),
}
var badHeaders2 = Headers{ // Missing canonical name MMSI
	Fields: strings.Split("BaseDateTime,LAT,LON,SOG,COG,Heading,"+
		"VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ","),
}

var goodHeaders = Headers{
	Fields: strings.Split("MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,"+
		"VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ","),
}

var basicHeadersString = `MMSI,BaseDateTime,LAT,LON`
var defaultHeadersString = `MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo`
var nonCanonicalHeadersString = `MMSI,Timestamp,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo`

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("errorReader used for testing")
}

type errorWriter struct{}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("errorWriter used for testing")
}

type testReader struct {
	buf  bytes.Buffer
	data string
}

func newTestReader() *testReader {
	tr := testReader{
		buf:  bytes.Buffer{},
		data: "376494000, 2017-12-01T00:00:00, 30.28963, -110.73522, 9.4, 158.2, 511.0\n",
	}

	return &tr
}

func (tr testReader) Read(p []byte) (n int, err error) {
	tr.buf.WriteString(tr.data)
	return tr.buf.Read(p)
}

type errorMatcher struct{}

func (*errorMatcher) Match(*Record) (bool, error) {
	return false, fmt.Errorf("errorMatcher used for testing")
}

type trueMatcher struct{}

func (*trueMatcher) Match(*Record) (bool, error) {
	return true, nil
}
