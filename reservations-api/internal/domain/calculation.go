package domain

import (
	"sync"
	"time"
)

// CalculationResult holds the result of concurrent calculations
type CalculationResult struct {
	Available   bool
	BasePrice   float64
	Discount    float64
	FinalPrice  float64
	Restrictions []string
}

// PartialResult represents a single calculation result
type PartialResult struct {
	Type string
	Data interface{}
}

// AvailabilityResult represents table availability check result
type AvailabilityResult struct {
	Available bool
	Reason    string
}

// PriceResult represents price calculation result
type PriceResult struct {
	BasePrice float64
}

// DiscountResult represents discount calculation result
type DiscountResult struct {
	DiscountPercent float64
	DiscountAmount  float64
}

// CalculateReservationConcurrent performs concurrent calculations for a reservation
// Returns: availability, base price, discount, and final price
func CalculateReservationConcurrent(tableNumber int, guests int, dateTime time.Time, mealType string, ownerID string) (*CalculationResult, error) {
	results := make(chan PartialResult, 3)
	var wg sync.WaitGroup

	// Goroutine 1: Check table availability
	wg.Add(1)
	go func() {
		defer wg.Done()
		availability := checkTableAvailability(tableNumber, dateTime)
		results <- PartialResult{Type: "availability", Data: availability}
	}()

	// Goroutine 2: Calculate base price
	wg.Add(1)
	go func() {
		defer wg.Done()
		price := calculateBasePrice(guests, mealType)
		results <- PartialResult{Type: "price", Data: price}
	}()

	// Goroutine 3: Apply discounts
	wg.Add(1)
	go func() {
		defer wg.Done()
		discount := calculateDiscount(dateTime, ownerID, mealType)
		results <- PartialResult{Type: "discount", Data: discount}
	}()

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	finalResult := &CalculationResult{
		Available:    true,
		Restrictions: []string{},
	}

	for partial := range results {
		switch partial.Type {
		case "availability":
			avail := partial.Data.(AvailabilityResult)
			finalResult.Available = avail.Available
			if !avail.Available {
				finalResult.Restrictions = append(finalResult.Restrictions, avail.Reason)
			}

		case "price":
			price := partial.Data.(PriceResult)
			finalResult.BasePrice = price.BasePrice

		case "discount":
			disc := partial.Data.(DiscountResult)
			finalResult.Discount = disc.DiscountAmount
		}
	}

	// Calculate final price
	finalResult.FinalPrice = finalResult.BasePrice - finalResult.Discount
	if finalResult.FinalPrice < 0 {
		finalResult.FinalPrice = 0
	}

	return finalResult, nil
}

// checkTableAvailability simulates checking if a table is available
func checkTableAvailability(tableNumber int, dateTime time.Time) AvailabilityResult {
	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Simple logic: tables 1-10 always available, others check time
	if tableNumber <= 10 {
		return AvailabilityResult{Available: true}
	}

	// Check if it's a weekend (higher demand)
	weekday := dateTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		// Weekend: only available before 6 PM
		if dateTime.Hour() >= 18 {
			return AvailabilityResult{Available: false, Reason: "Table not available on weekend evenings"}
		}
	}

	return AvailabilityResult{Available: true}
}

// calculateBasePrice calculates the base price for a reservation
func calculateBasePrice(guests int, mealType string) PriceResult {
	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	basePerPerson := 0.0

	switch mealType {
	case MealTypeBreakfast:
		basePerPerson = 15.0
	case MealTypeLunch:
		basePerPerson = 25.0
	case MealTypeDinner:
		basePerPerson = 40.0
	case MealTypeEvent:
		basePerPerson = 75.0
	default:
		basePerPerson = 30.0
	}

	total := float64(guests) * basePerPerson

	return PriceResult{BasePrice: total}
}

// calculateDiscount applies discounts based on time and user
func calculateDiscount(dateTime time.Time, ownerID string, mealType string) DiscountResult {
	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	discountPercent := 0.0

	// Early bird discount (before 6 PM)
	if dateTime.Hour() < 18 && mealType == MealTypeDinner {
		discountPercent += 10.0
	}

	// Weekday discount (Monday-Thursday)
	weekday := dateTime.Weekday()
	if weekday >= time.Monday && weekday <= time.Thursday {
		discountPercent += 5.0
	}

	// Loyal customer discount (simulated - ID ending in even number)
	if len(ownerID) > 0 && (ownerID[len(ownerID)-1]%2 == 0) {
		discountPercent += 5.0
	}

	// Calculate discount amount (will be applied to base price later)
	// For now, just return the percent
	return DiscountResult{
		DiscountPercent: discountPercent,
		DiscountAmount:  0, // will be calculated in service layer
	}
}
