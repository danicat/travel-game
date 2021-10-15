package main

import "testing"

func TestKBHandler(t *testing.T) {
	hdl := NewKBHandler()
	key := hdl.Read()
	hdl.Cancel()
	if key != KeyNone {
		t.Fatal("expected KeyNone")
	}
}
