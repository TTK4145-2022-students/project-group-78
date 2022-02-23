package distributor

import "testing"

func TestDistributor(t *testing.T) {
	d := New(1)
	d.RegisterClient(nil, nil)
}