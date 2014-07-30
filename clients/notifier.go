package clients

type Notifier interface {
	Send(addresses []string, subject, content string) error
}
