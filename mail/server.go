package mail

import (
	"net"
    "net/mail"
	"net/smtp"
	"crypto/tls"
	"fmt"
	"io"
	"strconv"
)

type Server struct {
	address string //remote smtpserveraddress
	port int //remote smtpserver port
	username string //account to connect to remote smtpserver + e-mail sender
	pass string //account password
}

func NewServer(a string, po int, u, pa string) (*Server) {
	return &Server{a, po, u, pa	}
}

func (this Server) SendTLS(to, subj, body string) (error) {
	var c *smtp.Client
	var w io.WriteCloser
	defer func () {
		if w != nil {
			w.Close()
		}
		if c != nil {
			c.Quit()
		}
	}()
	
	from := mail.Address{"", this.username}
	too := mail.Address{"", to}

    // Setup headers
    headers := make(map[string]string)
    headers["From"] = this.username
    headers["To"] = to
    headers["Subject"] = subj
	
	// Setup message
    message := ""
    for k,v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body

    // Connect to the SMTP Server
    servername := this.address + ":" + strconv.Itoa(this.port)

    host, _, _ := net.SplitHostPort(servername)

    auth := smtp.PlainAuth("",
        this.username,
	    this.pass,
	    this.address)

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
	
    c, err = smtp.NewClient(conn, host)
    
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

    if err = c.Rcpt(too.Address); err != nil {
        return err
    }

    // Data
    w, err = c.Data()
    if err != nil {
        return err
    }
	
	_, err = w.Write([]byte(message))
    if err != nil {
        return err
    }
    return err
}

func (this Server) Send(to, subj, body string) (error) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", this.username, this.pass, this.address)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to2 := []string{to}
	msg := []byte("To: " + to +  "\r\n" +
		"Subject: " + subj + "\r\n" +
		"\r\n" + body + "\r\n")
	err := smtp.SendMail(this.address + ":" + strconv.Itoa(this.port), auth, this.username, to2, msg)
	return err
} 

