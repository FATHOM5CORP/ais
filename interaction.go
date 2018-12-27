package ais

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
)

// InteractionFields are the default column headers used to write a csv file of two vessel
// interactions. The first field InteractionHash is an ParirHash64 return value that uniquely
// identifies this interaction and Distance(nm) is the haversine distance between the two vessels.
const InteractionFields = "InteractionHash,Distance(nm)," +
	"MMSI_1,BaseDateTime_1,LAT_1,LON_1,SOG_1,COG_1,Heading_1,VesselName_1,IMO_1,CallSign_1,VesselType_1,Status_1,Length_1,Width_1,Draft_1,Cargo_1,Geohash_1," +
	"MMSI_2,BaseDateTime_2,LAT_2,LON_2,SOG_2,COG_2,Heading_2,VesselName_2,IMO_2,CallSign_2,VesselType_2,Status_2,Length_2,Width_2,Draft_2,Cargo_2,Geohash_2"

// RecordPair holds pointers to two Records.
type RecordPair struct {
	rec1 *Record
	rec2 *Record
}

// Interactions is an abstraction for two-vessel interactions.  It requires a set of
// Headers that correspond to the Record slices being compared and it requires a set of
// Headers for the output.  The default for OutputHeaders is the const InteractionFields
// with a nil dictionary. The data held by interactions is a
// map[hash]*RecordPair.  This guarantees a non-duplicative set of interactions in the output.
type Interactions struct {
	RecordHeaders Headers                // for the Records that will be used to create interactions
	OutputHeaders Headers                // for an output RecordSet that may be written from the 2-ship interactions
	hashIndices   [4]int                 // Headers index values for MMSI, BaseDateTime, LAT, and LON
	data          map[uint64]*RecordPair // uint64 index is PairHash64 return value
}

// NewInteractions creates a new set of interactions.  It requires a set of Headers from the
// RecordSet that will be searched for Interactions.  These Headers are required to contain "MMSI",
// "BaseDateTime", "LAT", and "LON" in order to uniquely identify an interaction. The returned
// *Interactions has its output file Headers set to ais.InteractionHeaders by default.
func NewInteractions(h Headers) (*Interactions, error) {
	inter := new(Interactions)
	inter.OutputHeaders = Headers{
		fields:     strings.Split(InteractionFields, ","),
		dictionary: nil,
	}
	inter.RecordHeaders = h
	inter.data = make(map[uint64]*RecordPair)

	// Find the index values for the required headers now so that the expensive parsing
	// operation only has to be perormed once at initilization
	mmsiIndex, _ := h.Contains("MMSI")
	timeIndex, _ := h.Contains("BaseDateTime")
	latIndex, _ := h.Contains("LAT")
	lonIndex, _ := h.Contains("LON")
	inter.hashIndices = [4]int{mmsiIndex, timeIndex, latIndex, lonIndex}

	return inter, nil
}

// Len returns the number of Interactions in the set.
func (inter *Interactions) Len() int {
	return len(inter.data)
}

// AddCluster adds all of the interactions in a given cluster to the set of Interactions
func (inter *Interactions) AddCluster(c *Cluster) error {
	for i := range c.Data() {
		err := inter.writeInteractions(c.data[i:])
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteInteraction appends to the set for each pair of interaction in the slice.
// Note that calls to writeInteractions stemming from a sliding window will not hold
// their order due to the randomization of ranging over a map.  This occurs because
// the Window holds its data in a map and after a Slide() the order of these records
// will be iterated differently. Therefore, this means that the PairHash for a given
// pair of records may be recorded as the hash of {rec1, rec2} or {rec2, rec1} and
// both must be checked for existence before a new *RecordPair is inserted into the
// interactions map.
func (inter *Interactions) writeInteractions(data []*Record) error {
	if len(data) <= 1 { // only write two vessel interactions
		return nil
	}
	rec1 := data[0]
	for _, rec2 := range data[1:] {
		// Ignore pairs where it is subsequent reports of the same MMSI
		if (*rec1)[inter.hashIndices[0]] == (*rec2)[inter.hashIndices[0]] {
			continue
		}
		hash, err := PairHash64(rec1, rec2, inter.hashIndices)
		hash2, err := PairHash64(rec2, rec1, inter.hashIndices)
		if err != nil {
			return fmt.Errorf("write interactions: %v", err)
		}
		_, ok1 := inter.data[hash]
		_, ok2 := inter.data[hash2]
		if !ok1 && !ok2 { // neither Record order has been inserted
			inter.data[hash] = &RecordPair{rec1, rec2}
		}
	}
	return nil
}

// Save the interactions to a CSV file.
func (inter *Interactions) Save(filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("interactions save: %v", err)
	}

	w := csv.NewWriter(out)
	err = w.Write(inter.OutputHeaders.fields)
	if err != nil {
		return fmt.Errorf("interactions save: %v", err)
	}
	w.Flush()

	latIndex, _ := inter.RecordHeaders.Contains("LAT")
	lonIndex, _ := inter.RecordHeaders.Contains("LON")

	written := 1
	for hash, pair := range inter.data {
		d, err := pair.rec1.Distance(*(pair.rec2), latIndex, lonIndex)
		if err != nil {
			return fmt.Errorf("interactions save: %v", err)
		}
		pairData := []string{fmt.Sprintf("%0#16x", hash), fmt.Sprintf("%.1f", d)}
		pairData = append(pairData, (*pair.rec1)...)
		pairData = append(pairData, (*pair.rec2)...)
		w.Write(pairData)
		written++
		if written%flushThreshold == 0 {
			w.Flush()
			if err := w.Error(); err != nil {
				return fmt.Errorf("interactions save: flush error: %v", err)
			}
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("interactions save: flush error: %v", err)
	}

	return nil
}

// PairHash64 returns a 64 bit fnv hash from two AIS records based on the string values of
// MMSI, BaseDateTime, LAT, and LON for each vessel. Indices must
// contain the index values in rec1 and rec2 for MMSI, BaseDateTime, LAT and LON.
func PairHash64(rec1, rec2 *Record, indices [4]int) (uint64, error) {
	h64 := fnv.New64a()
	for i := range indices {
		_, err := h64.Write([]byte((*rec1)[i]))
		if err != nil {
			return 0, err
		}
		_, err = h64.Write([]byte((*rec2)[i]))
		if err != nil {
			return 0, err
		}
	}

	return h64.Sum64(), nil
}
