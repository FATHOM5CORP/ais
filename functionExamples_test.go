// Tests and examples for package ais
package ais_test

import (
	"fmt"
	"strings"

	"github.com/FATHOM5/ais"
)

// Example demonstrates a simple use of the Distance function.
func ExampleRecord_Distance() {
	h := strings.Split("MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ",")
	headers := ais.NewHeaders(h, nil)
	latIndex, _ := headers.Contains("LAT")
	lonIndex, _ := headers.Contains("LON")

	data1 := strings.Split("477307900,2017-12-01T00:00:03,36.90512,-76.32652,0.0,131.0,352.0,FIRST,IMO9739666,VRPJ6,1004,moored,337,,,", ",")
	data2 := strings.Split("477307902,2017-12-01T00:00:03,36.91512,-76.22652,2.3,311.0,182.0,SECOND,IMO9739800,XHYSF,,underway using engines,337,,,", ",")
	rec1 := ais.Record(data1)
	rec2 := ais.Record(data2)

	nm, err := rec1.Distance(rec2, latIndex, lonIndex)
	if err != nil {
		panic(err)
	}
	fmt.Printf("The ships are %.1fnm away from one another.\n", nm)

	// Output:
	// The ships are 4.8nm away from one another.
}

func ExampleRecord_ParseTime() {
	h := strings.Split("MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ",")
	data := strings.Split("477307900,2017-12-01T00:00:03,36.90512,-76.32652,0.0,131.0,352.0,FIRST,IMO9739666,VRPJ6,1004,moored,337,,,", ",")

	headers := ais.NewHeaders(h, nil)
	rec := ais.Record(data)

	timeIndex, _ := headers.Contains("BaseDateTime")

	t, err := rec.ParseTime(timeIndex)
	if err != nil {
		panic(err)
	}
	fmt.Printf("The record timestamp is at %s\n", t.Format(ais.TimeLayout))

	// Output
	// The record timestamp is at 2017-12-01T00:00:03
}
