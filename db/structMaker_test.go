package db

import (
	mytest "github.com/a2800276/goutils"
	"testing"
)

func TestMakeStructForTable(t *testing.T) {

	cols := []col{col{name: "id", dataType: "bigint"}, col{name: "received_at", dataType: "timestamp without time zone"}, col{name: "topic", dataType: "character varying"}, col{name: "msg", dataType: "character varying"}, col{name: "bgw", dataType: "character varying"}, col{name: "beacon", dataType: "character varying"}, col{name: "timestamp", dataType: "timestamp without time zone"}, col{name: "raw", dataType: "bytea"}, col{name: "rssi", dataType: "integer"}, col{name: "event_type", dataType: "integer"}}

	should := `type Event struct {
	Id int64
	ReceivedAt time.Time
	Topic string
	Msg string
	Bgw string
	Beacon string
	Timestamp time.Time
	Raw []byte
	Rssi int32
	EventType int32
}`
	is := makeStructForTable("event", cols)
	mytest.AssertEqual(t, is, should)
}
