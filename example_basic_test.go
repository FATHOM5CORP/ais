// Ex1 demonstrates the basics of creating, writing to and saving a RecordSet
package ais_test

import (
	"strings"

	"github.com/FATHOM5/ais"
)

// This example shows the basic usage of creating a new RecordSet and then
// using it to write a Record and finally saving the RecordSet to a csv file.
func Example() {
	rs := ais.NewRecordSet()
	defer rs.Close()

	h := strings.Split("MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo", ",")
	data := strings.Split("477307900,2017-12-01T00:00:03,36.90512,-76.32652,0.0,131.0,352.0,FIRST,IMO9739666,VRPJ6,1004,moored,337,,,", ",")

	rs.SetHeaders(ais.NewHeaders(h, nil))

	rec1 := ais.Record(data)
	err := rs.Write(rec1)
	if err != nil {
		panic(err)
	}
	err = rs.Flush()
	if err != nil {
		panic(err)
	}

	err = rs.Save("test.csv")
	if err != nil {
		panic(err)
	}

}
