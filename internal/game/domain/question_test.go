package domain

import "testing"

func TestValid(t *testing.T) {
	for _, q := range Questions {
		if !q.Valid() {
			t.Errorf("%d should be valid", q)
		}
	}
	for _, q := range []Question{-1, 10} {
		if q.Valid() {
			t.Errorf("%d should not be valid", q)
		}
	}
}
