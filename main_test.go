package main

import (
	"testing"
	"time"
)

func TestCheckOut(t *testing.T) {
	// Test case 1: Check In Mall Close
	idCard := checkIn(time.Date(2024, time.July, 10, 0, 0, 0, 0, time.UTC))
	expectedFee := 0
	actualFee := checkOut(idCard, time.Date(2024, time.July, 11, 6, 0, 0, 0, time.UTC))
	if actualFee != expectedFee {
		t.Errorf("Expected fee: %d, got: %d", expectedFee, actualFee)
	}

	// Test case 2: Check In Just Before Mall Close
	idCard = checkIn(time.Date(2024, time.July, 10, 9, 59, 0, 0, time.UTC))
	expectedFee = 0
	actualFee = checkOut(idCard, time.Date(2024, time.July, 11, 6, 0, 0, 0, time.UTC))
	if actualFee != expectedFee {
		t.Errorf("Expected fee: %d, got: %d", expectedFee, actualFee)
	}

	// Test case 3: Parking Free In 2 hour
	idCard = checkIn(time.Date(2024, time.July, 7, 10, 0, 0, 0, time.UTC))
	expectedFee = 0
	actualFee = checkOut(idCard, time.Date(2024, time.July, 7, 12, 0, 0, 0, time.UTC))
	if actualFee != expectedFee {
		t.Errorf("Expected fee: %d, got: %d", expectedFee, actualFee)
	}

	// Test case 4: Normal Case but before 12:00 AM
	idCard = checkIn(time.Date(2024, time.July, 8, 10, 0, 0, 0, time.UTC))
	expectedFee = 1200 // Replace with expected fee for this case
	actualFee = checkOut(idCard, time.Date(2024, time.July, 8, 23, 30, 0, 0, time.UTC))
	if actualFee != expectedFee {
		t.Errorf("Expected fee: %d, got: %d", expectedFee, actualFee)
	}

	// Test case 5: Parking More Than One Days
	idCard = checkIn(time.Date(2024, time.July, 8, 10, 0, 0, 0, time.UTC))
	expectedFee = 4600 // Replace with expected fee for this case
	actualFee = checkOut(idCard, time.Date(2024, time.July, 10, 10, 0, 0, 0, time.UTC))
	if actualFee != expectedFee {
		t.Errorf("Expected fee: %d, got: %d", expectedFee, actualFee)
	}

	// // Test case 6: Parking Night Time
	// idCard = checkIn(time.Date(2024, time.July, 10, 10, 0, 0, 0, time.UTC))
	// expectedFee = 2000 // Replace with expected fee for this case
	// actualFee = checkOut(idCard, time.Date(2024, time.July, 11, 7, 0, 0, 0, time.UTC))
	// if actualFee != expectedFee {
	// 	t.Errorf("Expected fee: %d, got: %d", expectedFee, actualFee)
	// }
}
