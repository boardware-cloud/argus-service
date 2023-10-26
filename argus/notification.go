package argus

import (
	"github.com/boardware-cloud/model/notification"
)

type NotificationGroup struct {
}

type Notification interface {
	Notify(message string)
}

func NotificationBackward(n notification.Notification) Notification {
	switch n.Entity().Type() {
	case "EMAIL":
		entity := n.Entity().(notification.Email)
		return EmailNotification{
			To:  entity.To,
			Cc:  entity.Cc,
			Bcc: entity.Bcc,
		}
	}
	return nil
}

type EmailNotification struct {
	To  []string
	Cc  []string
	Bcc []string
}

func (e EmailNotification) Notify(message string) {
	emailSender.SendHtml(
		emailSender.Email,
		"Uptime monitor alert",
		message,
		e.To,
		e.Cc,
		e.Bcc,
	)
}
