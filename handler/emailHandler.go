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
