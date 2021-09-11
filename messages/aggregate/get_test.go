package aggregate

import "testing"

const APIServer = "https://api2.aleph.im"

func TestAggregateMessageRetrievalSuccess(t *testing.T) {
	sgc := AggregateGetConfiguration{
		APIServer: APIServer,
		Address:   "0x629fBDA22F485720617C8f1209692484C0359D43",
		Keys:      []string{"satoshi"},
	}
	msg, err := Get(sgc)
	if err != nil {
		t.Errorf("Failed to retrieve store msg, reason: %v", err)
	}
	if msg == nil {
		t.Errorf("Aggregate msg should not be nil")
	}
}