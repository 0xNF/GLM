package notify

import (
	"net"
    "net/mail"
	  "net/smtp"
	"crypto/tls"
	"fmt"
	conf "github.com/0xNF/glm/src/internal/conf"
)

// SendMail2 Sends mail via the configured settings
func SendMail(emailConf *conf.GLMEmail, subj string, msg string) error {
	from := mail.Address{"GLM Monitor", emailConf.SenderAddress}
    to   := mail.Address{"", emailConf.RecipientAddress}
    body := msg

    // Setup headers
    headers := make(map[string]string)
    headers["From"] = from.String()
    headers["To"] = to.String()
    headers["Subject"] = subj

    // Setup message
    message := ""
    for k,v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body

    // Connect to the SMTP Server
    servername := emailConf.SmtpServer

    host, _, _ := net.SplitHostPort(servername)

    auth := smtp.PlainAuth("", emailConf.User, emailConf.Password, host)

    // TLS config
    tlsconfig := &tls.Config {
        InsecureSkipVerify: true,
        ServerName: host,
    }

    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    conn, err := tls.Dial("tcp", servername, tlsconfig)
    if err != nil {
        return err
    }

    c, err := smtp.NewClient(conn, host)
    if err != nil {
        return err
    }

    // Auth
    if err = c.Auth(auth); err != nil {
        return err
    }

    // To && From
    if err = c.Mail(from.Address); err != nil {
        return err
    }

    if err = c.Rcpt(to.Address); err != nil {
        return err
    }

    // Data
    w, err := c.Data()
    if err != nil {
        return err
    }

    _, err = w.Write([]byte(message))
    if err != nil {
        return err
    }

    err = w.Close()
    if err != nil {
        return err
    }

	c.Quit()
	
	return nil

}
