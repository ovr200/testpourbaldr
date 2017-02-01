package main

import "gopkg.in/gomail.v2"

func resetpassword(email, newpass string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mailfrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Changement de mot de passe")
	m.SetBody("text/plain", "Ton nouveau mot de passe est: "+newpass)

	d := gomail.NewDialer(smtpserver, smtpport, "", "")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func sendpassword(email, newpass string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mailfrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Nouveau compte")
	m.SetBody("text/plain", "Bienvenue, ton email est : "+email+" et ton pass est : "+newpass)

	d := gomail.NewDialer(smtpserver, smtpport, "", "")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}