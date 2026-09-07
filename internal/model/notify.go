package model

type Notify struct {
	UserID     int64
	Message    string
	LinkButton *LinkButton
}

type LinkButton struct {
	Text string
	Link string
}
