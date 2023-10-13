package argus

type NotificationGroup struct {
}

type Notification interface {
	Notify(message string)
}

type EmailNotification struct {
}
