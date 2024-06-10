package mail

import (
	"bitmoi/utilities"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config := utilities.GetConfig("../..")
	require.NotNil(t, config)

	gSender := NewGmailSender(config)

	subject := "test"
	content := "test message"
	to := []string{"yourheehee@gmail.com"}

	err := gSender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
