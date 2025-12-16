package creature

import "testing"

func TestFeed(t *testing.T) {
	p := Pet{Name: "griffin", Hungry: 50}
	expected := p.Hungry - FeedVal
	p.Feed()
	if p.Hungry != expected {
		t.Errorf("期望饥饿为30，实际饥饿值为%d", p.Hungry)
	}

	p.Hungry = 5
	p.Feed()
	if p.Hungry != 0 {
		t.Errorf("期望饥饿为0，实际饥饿值为%d", p.Hungry)
	}
}
