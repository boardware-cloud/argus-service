package argus

import (
	"fmt"
	"strings"
	"time"

	"github.com/boardware-cloud/model/notification"
)

type notificationData struct {
	Status    string
	CheckedAt time.Time
	Target    string
}

func (n notificationData) Message() string {
	return fmt.Sprintf(`
		<html>
			<body>
				<div>Url: %s</div>
				<div>Check time: %s</div>
				<div>Status: %s</div>
			</body>
		</html>`, n.Target, n.CheckedAt, n.Status)
}

type NotificationGroup struct{}

type Notification interface {
	Notify(notificationData)
}

func NotificationBackward(n notification.Notification) Notification {
	switch n.Entity().Type() {
	case "EMAIL":
		entity := n.Entity().(notification.Email)
		return EmailNotification{
			To:       entity.To,
			Cc:       entity.Cc,
			Bcc:      entity.Bcc,
			Template: entity.Template,
		}
	}
	return nil
}

type EmailNotification struct {
	To       []string
	Cc       []string
	Bcc      []string
	Template *string
}

func (e EmailNotification) Notify(notificationData notificationData) {
	var message string
	if e.Template != nil {
		message = templateStringHelper(notificationData, *e.Template)
	} else {
		message = notificationData.Message()
	}
	emailSender.SendHtml(
		emailSender.Email,
		"Uptime monitor alert",
		message,
		e.To,
		e.Cc,
		e.Bcc,
	)
}

func templateStringHelper(notificationData notificationData, template string) string {
	s := strings.ReplaceAll(template, "__STATUS__", notificationData.Status)
	s = strings.ReplaceAll(s, "__TIME__", notificationData.CheckedAt.Local().Format("2006 01-02 15:04"))
	s = strings.ReplaceAll(s, "__TARGET__", notificationData.Target)
	return s
}
