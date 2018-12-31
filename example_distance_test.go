// Ex2 demonstrates how to use the Record.Distance() function
package ais_test

import (
	"fmt"
	"strings"

	"github.com/FATHOM5/ais"
)

// This example demonstrates how to contruct two ais.Record types and compute
// the haversine distance between them.
func Example_distance() {

	h := ais.Headers{
		Fields: strings.Split("MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ","),
	}
	idxMap, ok := h.ContainsMulti("LAT", "LON")
	if !ok {
		panic("missing one or more required headers LAT and LON")
	}

	data1 := strings.Split("477307900,2017-12-01T00:00:03,36.90512,-76.32652,0.0,131.0,352.0,FIRST,IMO9739666,VRPJ6,1004,moored,337,,,", ",")
	data2 := strings.Split("477307902,2017-12-01T00:00:03,36.91512,-76.22652,2.3,311.0,182.0,SECOND,IMO9739800,XHYSF,,underway using engines,337,,,", ",")
	rec1 := ais.Record(data1)
	rec2 := ais.Record(data2)

	nm, err := rec1.Distance(rec2, idxMap["LAT"], idxMap["LON"])
	if err != nil {
		panic(err)
	}
	fmt.Printf("The ships are %.1fnm away from one another.\n", nm)

	// Output:
	// The ships are 4.8nm away from one another.
}
