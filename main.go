package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type ParkingRecord struct {
	ParkingID string
	CheckIn   time.Time
}

const csvFileName = "parking_records.csv"

func checkIn(checkin time.Time) string {
	nightParkingStart := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 0, 0, 0, 0, checkin.Location())
	nightParkingEnd := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 10, 0, 0, 0, checkin.Location())

	if checkin.Before(nightParkingEnd) && (checkin.Equal(nightParkingStart) || checkin.After(nightParkingStart)) {
		fmt.Println("Mall Closed. Please come back after 10 AM and before 12 AM.")
		return ""
	}

	parkingID := uuid.New().String()
	record := ParkingRecord{ParkingID: parkingID, CheckIn: checkin}

	file, err := os.OpenFile(csvFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return ""
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{record.ParkingID, record.CheckIn.Format(time.RFC3339)})
	if err != nil {
		fmt.Println("Error writing to CSV file:", err)
		return ""
	}

	return parkingID
}

func checkOut(parkingID string, checkout time.Time) int {
	file, err := os.Open(csvFileName)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return 0
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return 0
	}

	var checkin time.Time
	for _, record := range records {
		if record[0] == parkingID {
			checkin, err = time.Parse(time.RFC3339, record[1])
			if err != nil {
				fmt.Println("Error parsing check-in time:", err)
				return 0
			}
			break
		}
	}

	if checkin.IsZero() {
		fmt.Println("Parking ID not found.")
		return 0
	}

	// // Calculate the duration during night parking (0:00:00 AM - 10:00:00 AM)
	// nightParkingStart := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 0, 0, 0, 0, checkin.Location())
	// nightParkingEnd := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 10, 0, 0, 0, checkin.Location())

	// Calculate total duration in hours
	duration := checkout.Sub(checkin).Seconds()
	totalHours := 0
	if int(duration)%3600 == 0 {
		totalHours = int(duration / 3600)
	} else {
		totalHours = int((duration / 3600) + 1)
	}

	fmt.Println("Total Hour: ", totalHours)

	// Calculate fees based on rules
	totalFee := 0
	if totalHours <= 2 {
		totalFee = 0
	} else {
		// afterTwoHours := totalHours - 2

		// // Calculate fee from 12:00 PM to Midnight
		// feeFrom12PMtoMidnight := 0
		// if checkin.Hour() < 12 {
		// 	feeFrom12PMtoMidnight = 12 * 100
		// } else {
		// 	feeFrom12PMtoMidnight = (24 - checkin.Hour()) * 100
		// }

		// // Calculate fee from Midnight to checkout (if within night parking period)
		// feeFromMidnightTo7AM := 1000

		// // Total fee
		// totalFee = feeFrom12PMtoMidnight + feeFromMidnightTo7AM

		afterTwoHours := totalHours - 2
		nightParkingStart := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 0, 0, 0, 0, checkin.Location())
		nightParkingEnd := time.Date(checkin.Year(), checkin.Month(), checkin.Day(), 10, 0, 0, 0, checkin.Location())

		// Calculate the night parking fee if applicable

		fmt.Println(checkout.Before(nightParkingEnd))
		fmt.Println((checkout.Equal(nightParkingStart) || checkout.After(nightParkingStart)))

		if checkout.Before(nightParkingEnd) && (checkout.Equal(nightParkingStart) || checkout.After(nightParkingStart)) {
			// if checkout.Before(nightParkingEnd) { // Check Out Before 10 AM
			// 	totalFee += 1000
			// 	afterTwoHours -= int(duration / 3600) // Subtract the hours spent in night parking
			// } else {
			// 	nightDuration := nightParkingEnd.Sub(checkin)
			// 	totalFee += 1000
			// 	// Subtract the night duration from afterTwoHours
			// 	afterTwoHours -= int(nightDuration.Hours())
			// }
			// fmt.Println("Fuck You")
			totalFee += 1000

		}

		// Calculate the fee for hours after the first 2 hours
		if afterTwoHours > 0 {
			totalFee += afterTwoHours * 100
		}
	}

	return totalFee
}

func main() {
	idCard := ""
	totalFee := 0

	// Case 1: Check In Mall Close -> return 0
	idCard = checkIn(time.Date(2024, time.July, 10, 0, 0, 0, 0, time.UTC))
	totalFee = checkOut(idCard, time.Date(2024, time.July, 11, 6, 0, 0, 0, time.UTC))
	fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

	idCard = checkIn(time.Date(2024, time.July, 10, 9, 59, 0, 0, time.UTC))
	totalFee = checkOut(idCard, time.Date(2024, time.July, 11, 6, 0, 0, 0, time.UTC))
	fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

	// Case 2: Parking Free In 2 hour
	idCard = checkIn(time.Date(2024, time.July, 7, 10, 0, 0, 0, time.UTC))
	totalFee = checkOut(idCard, time.Date(2024, time.July, 7, 12, 0, 0, 0, time.UTC))
	fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

	// Case 3: Normal Case but before 12:00 AM
	idCard = checkIn(time.Date(2024, time.July, 8, 10, 0, 0, 0, time.UTC))
	totalFee = checkOut(idCard, time.Date(2024, time.July, 8, 23, 30, 0, 0, time.UTC))
	fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

	// Case 4: Parking More Than One Days
	idCard = checkIn(time.Date(2024, time.July, 8, 10, 0, 0, 0, time.UTC))
	totalFee = checkOut(idCard, time.Date(2024, time.July, 10, 10, 0, 0, 0, time.UTC))
	fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

	// Case 5: Parking New Month
	idCard = checkIn(time.Date(2024, time.July, 31, 10, 0, 0, 0, time.UTC))
	totalFee = checkOut(idCard, time.Date(2024, time.August, 1, 10, 0, 0, 0, time.UTC))
	fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

	// // Case 6: Parking Night Time
	// idCard = checkIn(time.Date(2024, time.July, 10, 10, 0, 0, 0, time.UTC))
	// totalFee = checkOut(idCard, time.Date(2024, time.July, 11, 7, 0, 0, 0, time.UTC))
	// fmt.Printf("Total Fee : %d Baht\n\n", totalFee)

}
