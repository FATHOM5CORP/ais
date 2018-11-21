package ais

import (
	"fmt"
	"sort"
	"strconv"
)

// ByGeohash implements the sort.Interface for creating a RecordSet
// sorted by Geohash.
type ByGeohash struct {
	h    Headers
	data *[]Record
}

// NewByGeohash returns a data structure suitable for sorting using
// the sort.Interface tools.
func NewByGeohash(rs *RecordSet) (*ByGeohash, error) {
	bg := new(ByGeohash)
	bg.h = rs.Headers()

	// Read the data from the underlying Recordset into a slice
	var err error
	bg.data, err = rs.loadRecords()
	if err != nil {
		return nil, fmt.Errorf("new bygeohash: unable to load data: %v", err)
	}

	return bg, nil
}

// Len function to implement the sort.Interface.
func (bg ByGeohash) Len() int { return len(*bg.data) }

// Swap function to implement the sort.Interface.
func (bg ByGeohash) Swap(i, j int) {
	(*bg.data)[i], (*bg.data)[j] = (*bg.data)[j], (*bg.data)[i]
}

//Less function to implement the sort.Interface.
func (bg ByGeohash) Less(i, j int) bool {
	geoIndex, ok := bg.h.Contains("Geohash")
	if !ok {
		panic("bybeohash less: headers does not contain Geohash")
	}
	g1, err := strconv.ParseUint((*bg.data)[i][geoIndex], 0, 64)
	if err != nil {
		panic(err)
	}
	g2, err := strconv.ParseUint((*bg.data)[j][geoIndex], 0, 64)
	if err != nil {
		panic(err)
	}
	return g1 < g2
}

// SortByGeohash returns a pointer to a new RecordSet sorted in ascending order
// by Geohash.
func (rs *RecordSet) SortByGeohash() (*RecordSet, error) {
	rs2 := NewRecordSet()
	rs2.SetHeaders(rs.Headers())

	bg, err := NewByGeohash(rs)
	if err != nil {
		return nil, fmt.Errorf("sortbygeohash: %v", err)
	}

	sort.Sort(bg)

	// Write the reports to the new RecordSet
	// NOTE: Headers are written only when the RecordSet is saved to disk
	written := 0
	for _, rec := range *bg.data {
		rs2.Write(rec)
		written++
		if written%flushThreshold == 0 {
			err := rs2.Flush()
			if err != nil {
				return nil, fmt.Errorf("sortbygeohash: flush error writing to new recordset: %v", err)
			}
		}
	}
	err = rs2.Flush()
	if err != nil {
		return nil, fmt.Errorf("sortbygeohash: flush error writing to new recordset: %v", err)
	}

	return rs2, nil
}
