package services

import (
	"fmt"
	"restaurant-system/internal/domain"
	"time"

	"github.com/google/uuid"
)

// ReviewService handles review business logic
type ReviewService struct {
	// Add repository when implemented
}

// NewReviewService creates a new review service
func NewReviewService() *ReviewService {
	return &ReviewService{}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(req CreateReviewRequest) (*domain.Review, error) {
	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Validate individual ratings if provided
	if req.FoodRating != nil && (*req.FoodRating < 1 || *req.FoodRating > 5) {
		return nil, fmt.Errorf("food rating must be between 1 and 5")
	}
	if req.ServiceRating != nil && (*req.ServiceRating < 1 || *req.ServiceRating > 5) {
		return nil, fmt.Errorf("service rating must be between 1 and 5")
	}
	if req.AmbienceRating != nil && (*req.AmbienceRating < 1 || *req.AmbienceRating > 5) {
		return nil, fmt.Errorf("ambience rating must be between 1 and 5")
	}

	review := &domain.Review{
		ID:            uuid.New().String(),
		ReservationID: req.ReservationID,
		CustomerName:  req.CustomerName,
		Rating:        req.Rating,
		Title:         req.Title,
		Comment:       req.Comment,
		IsAnonymous:   req.IsAnonymous,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set individual ratings if provided
	if req.FoodRating != nil {
		review.FoodRating = *req.FoodRating
	}
	if req.ServiceRating != nil {
		review.ServiceRating = *req.ServiceRating
	}
	if req.AmbienceRating != nil {
		review.AmbienceRating = *req.AmbienceRating
	}

	// TODO: Save to database
	// err := s.reviewRepo.Create(review)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to create review: %w", err)
	// }

	return review, nil
}

// CreateReviewRequest represents the request to create a review
type CreateReviewRequest struct {
	ReservationID  string `json:"reservation_id" binding:"required"`
	CustomerName   string `json:"customer_name" binding:"required"`
	Rating         int    `json:"rating" binding:"required,min=1,max=5"`
	Title          string `json:"title"`
	Comment        string `json:"comment"`
	FoodRating     *int   `json:"food_rating"`
	ServiceRating  *int   `json:"service_rating"`
	AmbienceRating *int   `json:"ambience_rating"`
	IsAnonymous    bool   `json:"is_anonymous"`
}
