package entity

import "time"

type User struct {
	Id            int     `json:"id"`
	Email         string  `json:"email"`
	Password      string  `json:"-"`
	DepositAmount float64 `json:"deposit_amount"`
}

type Equipment struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	Availability    int     `json:"availability"`
	DailyRentalCost float64 `json:"daily_rental_cost"`
	Category        string  `json:"category"`
}

type RentalHistory struct {
	Id          int     `json:"id"`
	UserId      int     `json:"user_id"`
	EquipmentId int     `json:"equipment_id"`
	RentalDate  string  `json:"rental_date"`
	ReturnDate  string  `json:"return_date"`
	TotalCost   float64 `json:"total_cost"`
}

type Payment struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	RentalHistoryID int       `json:"rental_history_id"`
	PaymentDate     time.Time `json:"payment_date"`
	IsDeposit       bool      `json:"is_deposit"`
}
