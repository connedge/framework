package store

import (
	"context"
	"gorm.io/gorm"
)

type Model[E any] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}

type GormStore[M Model[E], E any] struct {
	db *gorm.DB
}

func NewStore[M Model[E], E any](db *gorm.DB) *GormStore[M, E] {
	return &GormStore[M, E]{
		db: db,
	}
}

func (r *GormStore[M, E]) FindByID(ctx context.Context, id any) (E, error) {
	var model M
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return *new(E), err
	}

	return model.ToEntity(), nil
}
