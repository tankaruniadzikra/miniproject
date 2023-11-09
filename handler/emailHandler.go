package handler

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func sendRegistrationEmail(recipientEmail string) error {
	// Konfigurasi email
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", "tankaruniadzikra1@gmail.com") // Ganti dengan alamat email pengirim
	m.SetHeader("To", recipientEmail)                  // Menggunakan alamat email penerima dari parameter fungsi
	m.SetHeader("Subject", "Selamat bergabung!")
	m.SetBody("text/plain", "Selamat, pendaftaran Anda berhasil!")

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	// Email terkirim dengan sukses
	return nil
}

func sendTopUpEmail(userEmail string, depositAmount float64) error {
	// Konfigurasi email
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", "tankaruniadzikra1@gmail.com") // Ganti dengan alamat email pengirim
	m.SetHeader("To", userEmail)                       // Menggunakan alamat email penerima dari parameter fungsi
	m.SetHeader("Subject", "Deposit Top-Up")
	m.SetBody("text/plain", fmt.Sprintf("Deposit Anda telah ditambahkan sebesar %.2f", depositAmount))

	// Mengirim email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	// Email terkirim dengan sukses
	return nil
}

// key := "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"
// from := "tankaruniadzikra1@gmail.com"
// host := "smtp-relay.sendinblue.com"
// port := 587
// to := "tankaruniadzikra@gmail.com"

// m := gomail.NewMessage()

// m.SetHeader("From", from)
// m.SetHeader("To", to)

// m.SetHeader("Subject", "Registration Successful")
// m.SetBody("text/plain", "Thank you for registering with our service!")

// d := gomail.NewDialer(host, port, from, key)

// if err := d.DialAndSend(m); err != nil {
// 	return err
// }

// return nil

// m := gomail.NewMessage()

// m.SetHeader("From", "tankaruniadzikra1@gmail.com")
// m.SetHeader("To", userEmail)

// m.SetHeader("Subject", "Registration Successful")
// m.SetBody("text/plain", "Thank you for registering with our service!")

// d := gomail.NewDialer("smtp-relay.brevo.com", 587, "tankaruniadzikra@gmail.com", "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")

// if err := d.DialAndSend(m); err != nil {
// 	return err
// }

// return nil
