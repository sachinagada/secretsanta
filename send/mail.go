package send

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"net/smtp"

	"embed"
)

type Config struct {
	// Username is the email address from which all the emails will be sent.
	// eg: username@gmail.com
	Username string
	// Password is the password for the email. Note that if your account uses
	// 2FA, another password has to be generated to give the application access
	// to the account and send emails. More information can be found here for a
	// gmail account: https://support.google.com/mail/?p=InvalidSecondFactor
	Password string
	SMTPHost string
	SMPTPort string
	// Subject for the email messages
	Subject string
}

// DefaultConfig returns the default Config
func DefaultConfig() *Config {
	return &Config{
		SMTPHost: "smtp.gmail.com",
		SMPTPort: "587",
		Subject:  "Secret Santa!",
	}
}

type Mail struct {
	from     string             // email address of the organizer
	smtpAddr string             // smtpAddress (host:port)
	subject  string             // email subject
	auth     smtp.Auth          // mechanism used to authenticate SMTP server
	tmpl     *template.Template // template of the email message that's sent to everyone
}

//go:embed mail_template.html
var mailTmpl embed.FS

// NewMail is the constructor for Mail
func NewMail(c *Config) (*Mail, error) {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.SMTPHost)

	tmpl, parseErr := template.ParseFS(mailTmpl, "mail_template.html")
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing message template: %w", parseErr)
	}

	return &Mail{
		from:     c.Username,
		smtpAddr: net.JoinHostPort(c.SMTPHost, c.SMPTPort),
		subject:  c.Subject,
		auth:     auth,
		tmpl:     tmpl,
	}, nil
}

type Santa struct {
	// Name of the person sending the gift
	Name string
	// Email Address of the person sending the gift
	Addr string
	// Name of the recipient that Santa will be sending gift to
	Recipient string
}

func (s Santa) execute(tmpl *template.Template) ([]byte, error) {
	buf := new(bytes.Buffer)
	execErr := tmpl.Execute(buf, s)
	if execErr != nil {
		return nil, fmt.Errorf("error executing template: %w", execErr)
	}

	return buf.Bytes(), nil
}

type ErrSend struct {
	santa   Santa
	sendErr error
}

func (e ErrSend) Error() string {
	return fmt.Sprintf("error sending mail to %q at %q address for %q recipient: %s", e.santa.Name, e.santa.Addr, e.santa.Recipient, e.sendErr)
}

// Send sends the email to the Santas with name of their recipient.
func (m *Mail) Send(santas []Santa) error {

	var sendErrs []ErrSend

	for _, s := range santas {
		b, execErr := s.execute(m.tmpl)
		if execErr != nil {
			return fmt.Errorf("error executing template for %q Santa: %w", s.Name, execErr)
		}

		to := "To: " + s.Addr
		subject := "Subject: " + m.subject
		const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

		headers := fmt.Sprintf("%s\n%s\n%s\n\n", to, subject, mime)
		msg := append([]byte(headers), b...)

		smErr := smtp.SendMail(m.smtpAddr, m.auth, m.from, []string{s.Addr}, msg)
		if smErr != nil {
			sendErrs = append(sendErrs, ErrSend{santa: s, sendErr: smErr})
		}
	}

	if len(sendErrs) == 0 {
		return nil
	}

	return fmt.Errorf("error sending email to the following people: %v", sendErrs)
}
