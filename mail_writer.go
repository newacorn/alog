package alog

import (
	"github.com/phuslu/log"
	"github.com/wneessen/go-mail"
	"sync"
)

var _ = log.Writer(&MailWriter{})

type MailWriter struct {
	Host     string
	Port     int
	UserName string
	Password string
	Subject  string
	From     string
	To       []string
	TLS      bool
	msg      *mail.Msg
	cli      *mail.Client
	once     sync.Once
}

func (m *MailWriter) initMail() {
	msg := mail.NewMsg()
	err := msg.From(m.From)
	if err != nil {
		panic(err)
	}
	err = msg.To(m.To...)
	if err != nil {
		panic(err)
	}
	msg.Subject(m.Subject)
	client, err := mail.NewClient(m.Host, mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(m.UserName), mail.WithPassword(m.Password))
	if err != nil {
		panic(err)
	}
	m.cli = client
	m.msg = msg
}

func (m *MailWriter) WriteEntry(entry *log.Entry) (n int, err error) {
	m.once.Do(func() {
		m.initMail()
	})
	n = len(entry.Value())
	m.msg.SetBodyString(mail.TypeTextPlain, string(entry.Value()))
	err = m.cli.DialAndSend(m.msg)
	return
}

type mailConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	Subject  string
	From     string
	To       []string
}

func NewFileMailLogger(logFile string, mailCfg mailConfig) log.Logger {
	return log.Logger{
		Caller: 1,
		Writer: &log.MultiEntryWriter{&log.FileWriter{
			Filename:  logFile,
			LocalTime: true,
			MaxSize:   100 * 1024 * 1024,
		}, &MailWriter{
			UserName: mailCfg.UserName,
			Password: mailCfg.Password,
			Host:     mailCfg.Host,
			Port:     mailCfg.Port,
			TLS:      true,
			From:     mailCfg.From,
			To:       mailCfg.To,
			Subject:  mailCfg.Subject,
		}},
	}
}
