package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (*Connection, error)
	Get(ctx context.Context, id string) (*Connection, error)
	GetByName(ctx context.Context, name string) (*Connection, error)
	List(ctx context.Context, opts *ListOptions) ([]*Connection, int, error)
	ListByType(ctx context.Context, connType string) ([]*Connection, error)
	Update(ctx context.Context, id string, input UpdateInput) (*Connection, error)
	Delete(ctx context.Context, id string) error
	DeleteWithOptions(ctx context.Context, id string, opts *DeleteOptions) error
}

type service struct {
	repo      Repository
	validator *validator.Validate
}

func NewService(repo Repository) Service {
	return &service{
		repo:      repo,
		validator: GetValidator(),
	}
}

func (s *service) Create(ctx context.Context, input CreateInput) (*Connection, error) {
	if err := s.validator.Struct(input); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	if input.Config != nil {
		if err := ValidateConnectionConfig(input.Type, input.Config); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}

	now := time.Now()
	conn := &Connection{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Type:        input.Type,
		Description: input.Description,
		Config:      input.Config,
		Tags:        input.Tags,
		CreatedBy:   input.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, conn); err != nil {
		return nil, fmt.Errorf("creating connection: %w", err)
	}

	log.Info().
		Str("connection_id", conn.ID).
		Str("name", conn.Name).
		Str("type", conn.Type).
		Str("created_by", conn.CreatedBy).
		Msg("Connection created")

	return conn, nil
}

func (s *service) Get(ctx context.Context, id string) (*Connection, error) {
	conn, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *service) GetByName(ctx context.Context, name string) (*Connection, error) {
	conn, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *service) List(ctx context.Context, opts *ListOptions) ([]*Connection, int, error) {
	return s.repo.List(ctx, opts)
}

func (s *service) ListByType(ctx context.Context, connType string) ([]*Connection, error) {
	return s.repo.ListByType(ctx, connType)
}

func (s *service) Update(ctx context.Context, id string, input UpdateInput) (*Connection, error) {
	conn, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		conn.Name = *input.Name
	}
	if input.Description != nil {
		conn.Description = input.Description
	}

	if input.Config != nil {
		if err := ValidateConnectionConfig(conn.Type, input.Config); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
		conn.Config = input.Config
	}

	if input.Tags != nil {
		conn.Tags = input.Tags
	}

	conn.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, conn); err != nil {
		return nil, fmt.Errorf("updating connection: %w", err)
	}

	log.Info().
		Str("connection_id", conn.ID).
		Str("name", conn.Name).
		Msg("Connection updated")

	return conn, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.DeleteWithOptions(ctx, id, &DeleteOptions{TeardownSchedules: false})
}

func (s *service) DeleteWithOptions(ctx context.Context, id string, opts *DeleteOptions) error {
	if err := s.repo.DeleteWithOptions(ctx, id, opts); err != nil {
		return fmt.Errorf("deleting connection: %w", err)
	}

	if opts != nil && opts.TeardownSchedules {
		log.Info().
			Str("connection_id", id).
			Msg("Connection deleted with associated schedules")
	} else {
		log.Info().
			Str("connection_id", id).
			Msg("Connection deleted")
	}

	return nil
}
