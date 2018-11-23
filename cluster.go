package ais

import (
	"bytes"
	"fmt"
	"strconv"
)

// Cluster is an abstraction for a []*Record. The intent is that a Cluster of
// Records are vessels that share the same geohash
type Cluster struct {
	data []*Record
}

// Append adds a *Record to the underlying slice managed by the Cluster
func (c *Cluster) Append(rec *Record) {
	c.data = append(c.data, rec)
}

// Size returns the length of the underlying slice managed by the Cluster.
func (c *Cluster) Size() int { return len(c.data) }

func (c *Cluster) String() string {
	buf := bytes.Buffer{}
	for _, rec := range c.data {
		buf.WriteString(fmt.Sprint(*rec))
		buf.WriteString("\n")
	}
	return buf.String()
}

// Data returns the encapsulated data in the Cluster
func (c *Cluster) Data() []*Record { return c.data }

// ClusterMap is an abstraction for a map[geohash]*Cluster.
type ClusterMap map[uint64]*Cluster

// FindClusters returns a ClusterMap that groups Records in the window
// into common Clusters that share the same geohash.  It requires that
// the RecordSet Window is operating on has a 'Geohash' field stored as
// a Uint64 with the proper prefix for the hash (i.e. 0x for hex representation).
func (win *Window) FindClusters(geohashIndex int) ClusterMap {
	cm := make(ClusterMap)
	for _, rec := range win.Data {
		geoString := (*rec)[geohashIndex]
		geohash, err := strconv.ParseUint(geoString, 0, 64)
		if err != nil {
			panic(err)
		}
		if cluster, ok := cm[geohash]; ok {
			cluster.Append(rec)
		} else {
			cm[geohash] = func(r *Record) *Cluster {
				cl := new(Cluster)
				cl.Append(r)
				return cl
			}(rec)
		}
	}
	return cm
}
