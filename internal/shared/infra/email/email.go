// TODO: mover para pkg
package email

type Email struct {
	From    string
	To      []string
	Html    string
	Subject string
	Cc      []string
	ReplyTo string
}
