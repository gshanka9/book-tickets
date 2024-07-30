package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"runtime"
)

var logger = logrus.New()

func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
}

func main() {
	http.HandleFunc("/book", bookHotelHandler)
	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func bookHotelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var bookingRequest struct {
		HotelID string `json:"hotel_id"`
		UserID  string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bookingRequest); err != nil {
		logWithStack(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if bookingRequest.HotelID == "" || bookingRequest.UserID == "" {
		err := logWithStack("HotelID and UserID must be provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Infof("Booking successful for hotel %s by user %s", bookingRequest.HotelID, bookingRequest.UserID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Booking successful"))
}

func logWithStack(msg interface{}) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Error(msg)
	if err, ok := msg.(error); ok {
		return err
	}
	return fmt.Errorf("%v", msg)
}
