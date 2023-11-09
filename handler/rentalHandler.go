package handler

import (
	"errors"
	"io"
	"miniproject/config"
	"miniproject/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// @Summary Rent an equipment
// @Description Rent an equipment based on rental date, return date, and equipment ID
// @ID rent-equipment
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param request body entity.RentalHistory true "Rental request body"
// @Success 201 {object} map[string]interface{} "Equipment rented successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 400 {object} map[string]interface{} "Equipment is not available"
// @Failure 500 {object} map[string]interface{} "Failed to rent equipment" "Failed to update equipment availability"
// @Router /rent [post]
func RentEquipment(c echo.Context) error {
	userID := c.Get("user").(int)

	input := new(entity.RentalHistory)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Periksa ketersediaan stok
	var availability int
	config.DB.Raw("SELECT availability FROM equipments WHERE id = ?", input.EquipmentId).Scan(&availability)

	if availability <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Equipment is not available"})
	}

	// Hitung total biaya berdasarkan biaya harian dan selisih tanggal
	startDate, _ := time.Parse("2006-01-02", input.RentalDate)
	endDate, _ := time.Parse("2006-01-02", input.ReturnDate)
	days := int(endDate.Sub(startDate).Hours() / 24)

	var dailyRentalCost float64
	config.DB.Raw("SELECT daily_rental_cost FROM equipments WHERE id = ?", input.EquipmentId).Scan(&dailyRentalCost)

	totalCost := dailyRentalCost * float64(days)

	rental := entity.RentalHistory{
		UserId:      userID,
		EquipmentId: input.EquipmentId,
		RentalDate:  input.RentalDate,
		ReturnDate:  input.ReturnDate,
		TotalCost:   totalCost, // Menyimpan total biaya di database
	}

	if err := config.DB.Create(&rental).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to rent equipment"})
	}

	// Mengurangi stok peralatan di database
	if err := config.DB.Exec("UPDATE equipments SET availability = availability - 1 WHERE id = ?", input.EquipmentId).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to update equipment availability"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "success rent", "rental": rental})
}

// @Summary Get all rental histories for a user
// @Description Get all rental histories for the authenticated user
// @ID get-all-rental-histories
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Success 200 {object} []entity.RentalHistory "List of rental histories"
// @Failure 500 {object} map[string]interface{} "Failed to fetch rental histories"
// @Router /rental-histories [get]
func GetAllRentalHistories(c echo.Context) error {
	userID := c.Get("user").(int)

	var rentalHistories []entity.RentalHistory
	if err := config.DB.Where("user_id = ?", userID).Find(&rentalHistories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to fetch rental histories"})
	}

	// Ambil fakta menarik dari API
	fact, err := getInterestingFact()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to get interesting fact"})
	}

	// Tambahkan fakta menarik ke respons
	response := map[string]interface{}{
		"rentalHistories": rentalHistories,
		"interestingFact": fact,
	}

	return c.JSON(http.StatusOK, response)
}

func getInterestingFact() (string, error) {
	url := "https://facts-by-api-ninjas.p.rapidapi.com/v1/facts"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "3a35f7776bmsh8a2c394316f6f25p1dc119jsn1bd33e7e6e55")
	req.Header.Add("X-RapidAPI-Host", "facts-by-api-ninjas.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// @Summary Delete a rental history by ID
// @Description Delete a rental history by its ID
// @ID delete-rental-history
// @Accept json
// @Produce json
// @Param authorization header string true "JWT token"
// @Param id path int true "Rental history ID"
// @Success 200 {object} map[string]interface{} "Rental history deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Rental history not found"
// @Failure 500 {object} map[string]interface{} "Failed to get rental history" "Failed to delete rental history" "Failed to update equipment availability"
// @Router /rental-histories/{id} [delete]
func DeleteRentalHistory(c echo.Context) error {
	userID := c.Get("user").(int) // Mendapatkan ID pengguna dari token

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid ID"})
	}

	rental := entity.RentalHistory{}
	if err := config.DB.First(&rental, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Rental history not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to get rental history"})
	}

	// Memeriksa apakah pengguna yang menghapus rental history adalah pemiliknya
	if rental.UserId != userID {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Unauthorized"})
	}

	// Simpan EquipmentId sebelum menghapus
	equipmentID := rental.EquipmentId

	if err := config.DB.Delete(&rental).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to delete rental history"})
	}

	// Tambahkan stok peralatan kembali
	if err := config.DB.Exec("UPDATE equipments SET availability = availability + 1 WHERE id = ?", equipmentID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to update equipment availability"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Rental history deleted"})
}
