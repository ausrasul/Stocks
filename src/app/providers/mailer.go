package providers

import(
	"github.com/astaxie/beego"
	"github.com/go-gomail/gomail"
	"time"
)

func SendNewsLetter(text string) error{
	users, err := GetUsers()
	if err != nil {
		return err
	}


	d := gomail.NewPlainDialer("send.one.com", 465, "aus@emailaddress.com", "password")
	s, err := d.Dial()
	if err != nil {
    	return err
	}
	m := gomail.NewMessage()
	for _, r := range users {
		expT, err := time.Parse("2006-01-02", r.Expires)
		if err != nil || expT.Before(time.Now()) {
			beego.Debug("User expired ", err)
			continue
		}
		m.SetAddressHeader("From", "aus@emailaddress.com", "Stock Picker")
		m.SetAddressHeader("To", r.Email, r.Name)
		m.SetHeader("Subject", "Stock Picks")
		m.SetBody("text/html", text) //fmt.Sprintf("Hello %s!", r.Name))
		if err := gomail.Send(s, m); err != nil {
			beego.Debug("Could not send email to ", r.Email, err)
		}
		m.Reset()
	}
	return nil
}
