package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

const (
	MIMEText = "text/plain"
	MIMEHTML = "text/html"

	boundary = "bbe6dfbb11db4ec9ab4f7cbae59b3415"
)

type Mail struct {
	From        mail.Address
	To          []string
	CC          []string
	BCC         []string
	ReplyTo     string
	Subject     string
	Content     string
	ContentType string
	Headers     []Header
	Attachments []Attachment
}

type Attachment struct {
	Name   string
	Data   []byte
	Inline bool
}

type Header struct {
	Key   string
	Value string
}

func (m *Mail) AttachFile(filePath string, inline bool) (err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	_, fileName := filepath.Split(filePath)
	m.Attachments = append(m.Attachments, Attachment{
		Name:   fileName,
		Data:   data,
		Inline: inline,
	})
	return nil
}

func (m *Mail) AttachData(name string, data []byte, inline bool) {
	m.Attachments = append(m.Attachments, Attachment{
		Name:   name,
		Data:   data,
		Inline: inline,
	})
}

func (m *Mail) Receivers() []string {
	receivers := m.To
	for _, cc := range m.CC {
		receivers = append(receivers, cc)
	}
	for _, bcc := range m.BCC {
		receivers = append(receivers, bcc)
	}
	return receivers
}

func (m *Mail) Build() (data []byte, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = v.(error)
		}
	}()

	buffer := &bytes.Buffer{}
	m.buildHeader(buffer)
	m.buildContent(buffer)
	m.buildAttachments(buffer)

	return buffer.Bytes(), nil
}

func (m *Mail) buildHeader(w io.Writer) {
	m.writeHeader(w, "From", m.From.String())
	m.writeHeader(w, "Subject", fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(m.Subject))))
	m.writeHeader(w, "Mime-Version", "1.0")
	m.writeHeader(w, "Date", time.Now().Format(time.RFC1123Z))
	if len(m.CC) > 0 {
		m.writeHeader(w, "Cc", strings.Join(m.CC, ","))
	}
	if len(m.ReplyTo) > 0 {
		m.writeHeader(w, "Reply-To", m.ReplyTo)
	}
	if len(m.Attachments) > 0 {
		m.writeHeader(w, "Content-Type", "multipart/mixed; writeBoundary="+boundary)
		m.writeString(w, "\r\n")
		m.writeBoundary(w, false)
	}
}

func (m *Mail) buildContent(w io.Writer) {
	m.writeHeader(w, "Content-Type", fmt.Sprintf("%s;charset=utf-8", m.ContentType))
	m.writeString(w, "\r\n")
	m.writeString(w, m.Content)
	m.writeString(w, "\r\n")
}

func (m *Mail) buildAttachments(w io.Writer) {
	if len(m.Attachments) == 0 {
		return
	}
	for _, attachment := range m.Attachments {
		m.writeString(w, "\r\n")
		m.writeBoundary(w, false)
		if attachment.Inline {
			m.writeHeader(w, "Content-Type", "message/rfc822")
			m.writeHeader(w, "Content-Disposition", fmt.Sprintf("inline; filename=%q", attachment.Name))
			m.writeString(w, "\r\n")
			m.write(w, attachment.Data)
		} else {
			m.writeHeader(w, "Content-Transfer-Encoding", "base64")
			m.writeHeader(w, "Content-Disposition", fmt.Sprintf(`attachment; filename="=?UTF-8?B?%s?="`,
				base64.StdEncoding.EncodeToString([]byte(attachment.Name))))
			mimeType := mime.TypeByExtension(filepath.Ext(attachment.Name))
			if mimeType != "" {
				m.writeHeader(w, "Content-Type", mimeType)
			} else {
				m.writeHeader(w, "Content-Type", "application/octet-stream")
			}
			m.writeString(w, "\r\n")
			b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
			base64.StdEncoding.Encode(b, attachment.Data)
			// write base64 content in lines of up to 76 chars
			for i, l := 0, len(b); i < l; i++ {
				m.write(w, []byte{b[i]})
				if (i+1)%76 == 0 {
					m.writeString(w, "\r\n")
				}
			}
		}
		m.writeString(w, "\r\n")
		m.writeString(w, "\r\n")
	}
	m.writeBoundary(w, true)
}

func (m *Mail) writeHeader(w io.Writer, key, value string) {
	m.writeString(w, fmt.Sprintf("%s: %s\r\n", key, value))
}

func (m *Mail) writeBoundary(w io.Writer, end bool) {
	if end {
		m.writeString(w, "--%s--", boundary)
	} else {
		m.writeString(w, "--%s\r\n", boundary)
	}
}

func (m *Mail) write(w io.Writer, data []byte) {
	if _, err := w.Write(data); err != nil {
		panic(err)
	}
}

func (m *Mail) writeString(w io.Writer, format string, args ...interface{}) {
	if _, err := io.WriteString(w, fmt.Sprintf(format, args...)); err != nil {
		panic(err)
	}
}

type Sender struct {
	addr     string
	username string
	password string
	auth     smtp.Auth
}

func NewSender(addr, username, password string) *Sender {
	return &Sender{
		addr:     addr,
		username: username,
		password: password,
		auth:     smtp.PlainAuth("", username, password, strings.Split(addr, ":")[0]),
	}
}

func (s *Sender) Send(m Mail) (err error) {
	data, err := m.Build()
	if err != nil {
		return err
	}
	return smtp.SendMail(s.addr, s.auth, m.From.Address, m.Receivers(), data)
}

func (s *Sender) SendText(from mail.Address, to []string, subject string, content string) error {
	return s.Send(Mail{
		From:        from,
		To:          to,
		Subject:     subject,
		Content:     content,
		ContentType: MIMEText,
	})
}

func (s *Sender) SendHTML(from mail.Address, to []string, subject string, content string) error {
	return s.Send(Mail{
		From:        from,
		To:          to,
		Subject:     subject,
		Content:     content,
		ContentType: MIMEHTML,
	})
}
