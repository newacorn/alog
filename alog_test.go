package alog

import "testing"

func _TestFileMailLogger(t *testing.T) {
	log := NewFileMailLogger("test.log", MailConfig{
		Host:     "smtp.126.com",
		Port:     465,
		UserName: "w1013d27@126.com",
		Password: "",
		Subject:  "simple fle mail logger test",
		From:     "w1013d27@126.com",
		To:       []string{"w1013d27@126.com"},
	})
	log.Error().Msg("test")
}
