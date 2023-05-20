package utilities

import "testing"

func TestTimeForamt(t *testing.T) {
	var unixSecond int64 = 1684577267
	formatted := EntryTimeFormatter(unixSecond)
	t.Log(formatted)
}
