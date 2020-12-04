package mailing

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

//Mailer object
type Mailer struct {
	message *gomail.Message
	dialer  *gomail.Dialer
}

//NewMailer mailer constructor
func NewMailer() *Mailer {
	mailer := &Mailer{}
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mailer.dialer = gomail.NewPlainDialer(os.Getenv("MAIL_HOST"), port, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASS"))
	return mailer
}

//InitMessage method
func (m *Mailer) InitMessage(to, subject string) *Mailer {
	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("MAIL_USER"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	m.message = message
	return m
}

//RecordingNotificationMailTemplate mail template for sharing recordings
var RecordingNotificationMailTemplate string = `
<section style="margin: 0 5%;">
  <header style="background: #0000FF;width: 100%;height: 78px;text-align: center;">
    <span style="display: inline-block;height: 100%;vertical-align: middle;"></span>
    <a style="text-decoration: none;" href="https://www.hexcord.com"><img alt="Hexcord Logo"
        style="vertical-align: middle;"
        src="https://res.cloudinary.com/hyper-debugger/image/upload/v1606919638/hexcord_logo.png" /></a>
  </header>
  <section style="background: #F8F8FF;">
    <div style="width: 70%;margin: auto;">
      <p
        style="font-weight:bold;font-family: Lato;font-style: normal;font-size: 16px;line-height: 26px;color: #080708;margin-bottom: 20px;">
        Hi {{recipientName}}</p>
      <p
        style="font-family: Lato;font-style: normal;font-size: 16px;line-height: 26px;color: #080708;margin-bottom: 20px;">
        {{senderName}} has shared a recording with you.
      </p>
      <div>
        <img src="{{thumbNail}}" alt="recording thumbnail" width="500" height="300" />
      </div>
      <span>{{notificationMessage}}</span>
      <p
        style="font-family: Lato;font-style: normal;font-size: 16px;line-height: 26px;color: #080708;margin-bottom: 20px;">
        You can find the recording here <a
          href="{{recordingLink}}">{{recordingLink}}</a>.
      </p>
      <p
        style="font-family: Lato;font-style: normal;font-size: 16px;line-height: 26px;color: #080708;margin-bottom: 20px;">
        Cheers!</p>
      <p
        style="font-family: Lato;font-style: normal;font-size: 16px;line-height: 26px;color: #080708;margin-bottom: 20px;">
        Hexcord.</p>
      <footer
        style="border-top: 1px solid #E6E6FF;border-bottom: 1px solid #E6E6FF;padding-top: 20px;margin-bottom: 20px;">
        <a style="text-decoration: none;" href="https://twitter.com/thehexcord">
          <span
            style="display: inline-block;width: 29px;height: 22.75px;background: url(https://res.cloudinary.com/hyper-debugger/image/upload/v1606919620/twitter_icon.png);"></span>
        </a>
      </footer>
    </div>
  </section>
</section>
`

//SendMail utility method
func (m *Mailer) SendMail(message string) error {
	m.message.SetBody("text/html", message)
	return m.dialer.DialAndSend(m.message)

	// return smtp.SendMail(os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"), auth, from, []string{to}, []byte(message))
}
