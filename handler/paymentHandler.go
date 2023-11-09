package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"miniproject/config"
	"miniproject/entity"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	xendit "github.com/xendit/xendit-go/v3"
	invoice "github.com/xendit/xendit-go/v3/invoice"
)

func TopupDeposit(c echo.Context) error {
	userID := c.Get("user").(int)

	// Mendapatkan jumlah deposit yang akan ditambahkan dari body request
	input := struct {
		DepositAmount float64 `json:"deposit_amount"`
	}{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Mendapatkan data pengguna dari database
	var user entity.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to get user"})
	}

	// Menambahkan dana ke saldo deposit
	user.DepositAmount += float64(input.DepositAmount) // Ubah tipe data ke float64
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to top-up deposit"})
	}

	// Kirim deposit ke Xendit
	if err := sendDepositToXendit(c, userID, float32(input.DepositAmount)); err != nil { // Ubah tipe data ke float32
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Gagal mengirim deposit ke Xendit"})
	}

	// Mengirim email notifikasi
	if err := sendTopUpEmail(user.Email, input.DepositAmount); err != nil {
		// Jika gagal mengirim email, Anda dapat menangani kesalahan di sini
		// Misalnya, Anda dapat mencatat kesalahan atau memberikan respons khusus
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Deposit topped up", "new_balance": user.DepositAmount})
}

func sendDepositToXendit(c echo.Context, userId int, depositAmount float32) error {
	judulInvoice := fmt.Sprintf("Invoice order user id = %v", userId)
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(judulInvoice, depositAmount)

	xenditClient := xendit.NewClient(os.Getenv("XENDIT_API_KEY"))

	resp, r, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return err
	}

	// response from `CreateInvoice`: Invoice
	fmt.Fprintf(os.Stdout, "Response from `InvoiceApi.CreateInvoice`: %v\n", resp)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "sukses xendit",
		"respon":  resp,
	})
}

func MakePayment(c echo.Context) error {
	userID := c.Get("user").(int)

	input := new(entity.Payment)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Cari rental history berdasarkan ID
	var rentalHistory entity.RentalHistory
	if err := config.DB.First(&rentalHistory, input.RentalHistoryID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Rental history not found"})
	}

	// Membuat pembayaran
	payment := entity.Payment{
		UserID:          userID,
		RentalHistoryID: input.RentalHistoryID,
		PaymentDate:     time.Now(),
		IsDeposit:       input.IsDeposit,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to make payment"})
	}

	// Mengurangi deposit pengguna jika pembayaran adalah deposit
	if input.IsDeposit {
		var user entity.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to get user"})
		}

		// Check if deposit is sufficient
		if user.DepositAmount < rentalHistory.TotalCost {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Insufficient deposit amount"})
		}

		user.DepositAmount -= rentalHistory.TotalCost
		if err := config.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to update user deposit"})
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "Payment made", "payment": payment})
}
