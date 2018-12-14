package ais

import (
	"strings"
	"time"
)

var testRec0 = Record{"376494000", "2017-12-01T00:00:00", "30.28963", "-110.73522", "9.4", "158.2", "511.0"}
var testRec1 = Record{"376494001", "2017-12-01T00:00:01", "31.28963", "-111.73522", "9.4", "158.2", "511.0"}
var testRec2 = Record{"376494002", "2017-12-01T00:00:02", "32.28963", "-112.73522", "9.4", "158.2", "511.0"}

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
	fields: strings.Split("MMSI,Timestamp,LAT,LON,SOG,COG,Heading,"+
		"VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ","),
}

var goodHeaders = Headers{
	fields: strings.Split("MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,"+
		"VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ","),
}
