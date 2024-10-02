package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MeherKandukuri/studioClasses_API/models"
)

type CreateClassRequest struct {
	ClassName string `json:"class_name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Capacity  int    `json:"capacity"`
}

type BookingRequest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

var bookings = make(map[string][]string)
var classStorage = make(map[time.Time]models.Class)

func PostCreateClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Used wrong request Method type to access this function", http.StatusMethodNotAllowed)
		return
	}

	var req CreateClassRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.ClassName == "" || req.StartDate == "" || req.EndDate == "" || req.Capacity <= 0 {
		http.Error(w, "Entered, Missing or invalid fields.", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "Invalid endDate format", http.StatusBadRequest)
		return
	}
	// normalizing dates:
	startDate, endDate = normalizeDate(startDate), normalizeDate(endDate)

	if startDate.After(endDate) {
		http.Error(w, "start date cannot be after end date", http.StatusBadRequest)
		return
	}

	class := models.Class{
		ClassName: req.ClassName,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
	}

	currentDate := startDate
	for !currentDate.After(endDate) {
		if _, exists := classStorage[currentDate]; exists {
			http.Error(w, fmt.Sprintf("Class already exists on %v", currentDate.Format("2006-01-02")), http.StatusConflict)
			return
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Reset currentDate to startDate
	currentDate = startDate
	for !currentDate.After(endDate) {
		// Add class to storage for each date
		classStorage[currentDate] = class
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Class created successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func PostCreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Used wrong request method type to access this function", http.StatusMethodNotAllowed)
		return
	}

	var reqBooking BookingRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBooking); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if reqBooking.Name == "" || reqBooking.Date == "" {
		http.Error(w, "Missing required fields: name or date or both", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", reqBooking.Date)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}
	date = normalizeDate(date)
	datestr := date.Format("2006-01-02")

	if _, found := classStorage[date]; !found {
		http.Error(w, "We don't have a class on this day", http.StatusBadRequest)
		return
	}

	booking := models.Booking{
		Name: reqBooking.Name,
		Date: date,
	}
	// This check is done considering there is only one person with one name. We can achieve this functionality using unique id.
	namesInClass := bookings[datestr]
	username := strings.ToLower(booking.Name)
	for _, name := range namesInClass {
		if strings.ToLower(name) == username {
			http.Error(w, "You have already enrolled into class", http.StatusConflict)
			return
		}
	}

	bookings[datestr] = append(bookings[datestr], booking.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message := fmt.Sprintf("%s has been enrolled for class on %s", booking.Name, datestr)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": message}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
