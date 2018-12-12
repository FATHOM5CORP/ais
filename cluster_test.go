package ais

import (
	"reflect"
	"strings"
	"testing"
)

func TestCluster_Append(t *testing.T) {
	type fields struct {
		data []*Record
	}
	type args struct {
		rec *Record
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int //len(cluster.data)
	}{
		{
			name: "append to empty",
			fields: fields{
				data: []*Record{},
			},
			args: args{
				rec: &Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
			},
			want: 1,
		},
		{
			name: "append empty to empty",
			fields: fields{
				data: []*Record{},
			},
			args: args{
				rec: &Record{},
			},
			want: 1,
		},
		{
			name: "append one to len(c.data)=1",
			fields: fields{
				data: []*Record{
					&Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
				},
			},
			args: args{
				rec: &Record{"376494001", "2017-12-01T00:00:04", "31.28963", "-115.73522", "9.4", "158.2", "511.0"},
			},
			want: 2,
		},
		{
			name: "append one to len(c.data)=2",
			fields: fields{
				data: []*Record{
					&Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
					&Record{"376494002", "2017-12-01T00:00:05", "32.28963", "-114.73522", "9.4", "158.2", "511.0"}},
			},
			args: args{
				rec: &Record{"376494001", "2017-12-01T00:00:04", "31.28963", "-115.73522", "9.4", "158.2", "511.0"},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				data: tt.fields.data,
			}
			c.Append(tt.args.rec)
			got := len(c.data)
			if got != tt.want {
				t.Errorf("Cluster.Append() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestCluster_Size(t *testing.T) {
	type fields struct {
		data []*Record
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "size of empty",
			fields: fields{
				data: []*Record{},
			},
			want: 0,
		},
		{
			name: "size of one",
			fields: fields{
				data: []*Record{
					&Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
				},
			},
			want: 1,
		},
		{
			name: "size of two",
			fields: fields{
				data: []*Record{
					&Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
					&Record{"376494002", "2017-12-01T00:00:05", "32.28963", "-114.73522", "9.4", "158.2", "511.0"},
				},
			},
			want: 2,
		},
		{
			name: "size of three",
			fields: fields{
				data: []*Record{
					&Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
					&Record{"376494000", "2017-12-01T00:00:03", "30.28963", "-116.73522", "9.4", "158.2", "511.0"},
					&Record{"376494002", "2017-12-01T00:00:05", "32.28963", "-114.73522", "9.4", "158.2", "511.0"},
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				data: tt.fields.data,
			}
			if got := c.Size(); got != tt.want {
				t.Errorf("Cluster.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_String(t *testing.T) {
	type fields struct {
		data []*Record
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty cluster",
			fields: fields{
				data: []*Record{},
			},
			want: "",
		},
		{
			name: "one record cluster",
			fields: fields{
				data: []*Record{
					&testRec0,
				},
			},
			want: testRec0String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				data: tt.fields.data,
			}
			got := c.String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Cluster.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_Data(t *testing.T) {
	type fields struct {
		data []*Record
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Record
	}{
		{
			name: "empty cluster",
			fields: fields{
				data: []*Record{},
			},
			want: []*Record{},
		},
		{
			name: "one record cluster",
			fields: fields{
				data: []*Record{
					&testRec0,
				},
			},
			want: []*Record{
				&testRec0,
			},
		},
		{
			name: "three record cluster",
			fields: fields{
				data: []*Record{
					&testRec0,
					&testRec1,
					&testRec2,
				},
			},
			want: []*Record{
				&testRec0,
				&testRec1,
				&testRec2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				data: tt.fields.data,
			}
			if got := c.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cluster.Data() = %v, want %v", got, tt.want)
			}
		})
	}
}
