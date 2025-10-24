package repository

import (
	"context"
	"time"

	"github.com/tuusuario/restaurant-reservas/reservations-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationRepository interface {
	Create(ctx context.Context, r *domain.Reservation) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error)
	Update(ctx context.Context, r *domain.Reservation) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type mongoRepo struct{ coll *mongo.Collection }

func NewMongoRepository(coll *mongo.Collection) ReservationRepository {
	return &mongoRepo{coll: coll}
}

func (m *mongoRepo) Create(ctx context.Context, r *domain.Reservation) (primitive.ObjectID, error) {
	r.ID = primitive.NilObjectID
	r.CreatedAt = time.Now()
	r.UpdatedAt = r.CreatedAt
	res, err := m.coll.InsertOne(ctx, r)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (m *mongoRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
	var r domain.Reservation
	if err := m.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (m *mongoRepo) Update(ctx context.Context, r *domain.Reservation) error {
	r.UpdatedAt = time.Now()
	_, err := m.coll.UpdateByID(ctx, r.ID, bson.M{"$set": r})
	return err
}

func (m *mongoRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := m.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
