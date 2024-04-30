package utilities

import "testing"

func TestEntryTimeFormatter(t *testing.T) {
	s := EntryTimeFormatter(0)
	t.Logf(s)
}
