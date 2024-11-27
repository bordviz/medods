package smtp

type SMTPServer struct{}

func NewSMTPServer() *SMTPServer {
	return &SMTPServer{}
}

func (s *SMTPServer) SendEmail(email string, message string) {}
