package notifyer

type Notifyer interface {
	Send(text string) (err error)
}
