package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ticket-system/internal/models"
	"ticket-system/internal/repository"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidStatus           = errors.New("invalid ticket status value")
	ErrTicketNotFound          = errors.New("ticket not found")
)

type TicketService interface {
	CreateTicket(ctx context.Context, userID, name, desc string) (*models.Ticket, error)
	ListTickets(ctx context.Context, userID string) ([]*models.Ticket, error)
	GetTicket(ctx context.Context, ticketID int64, userID string) (*models.Ticket, error)
	UpdateTicketStatus(ctx context.Context, ticketID int64, userID, status string) (*models.Ticket, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
}

func NewTicketService(ticketRepo repository.TicketRepository) TicketService {
	return &ticketService{ticketRepo: ticketRepo}
}

func (s *ticketService) CreateTicket(ctx context.Context, userID, name, desc string) (*models.Ticket, error) {
	if name == "" {
		return nil, errors.New("ticket name/title is required")
	}

	now := time.Now()
	ticket := &models.Ticket{
		Name:      name,
		Desc:      desc,
		Status:    models.StatusOpen,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.ticketRepo.CreateTicketTx(ctx, userID, ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

func (s *ticketService) ListTickets(ctx context.Context, userID string) ([]*models.Ticket, error) {
	return s.ticketRepo.ListTicketsByUserID(ctx, userID)
}

func (s *ticketService) GetTicket(ctx context.Context, ticketID int64, userID string) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetTicketByIDAndUserID(ctx, ticketID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrTicketNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to retrieve ticket: %w", err)
	}
	return ticket, nil
}

func (s *ticketService) UpdateTicketStatus(ctx context.Context, ticketID int64, userID, newStatus string) (*models.Ticket, error) {
	if newStatus != models.StatusOpen && newStatus != models.StatusInProgress && newStatus != models.StatusClosed {
		return nil, ErrInvalidStatus
	}

	ticket, err := s.GetTicket(ctx, ticketID, userID)
	if err != nil {
		return nil, err
	}

	if ticket.Status == newStatus {
		return ticket, nil
	}

	if ticket.Status == models.StatusClosed {
		return nil, fmt.Errorf("%w: closed tickets cannot be updated", ErrInvalidStatusTransition)
	}

	if ticket.Status == models.StatusInProgress && newStatus == models.StatusOpen {
		return nil, fmt.Errorf("%w: cannot move from in_progress back to open", ErrInvalidStatusTransition)
	}

	updatedTicket, err := s.ticketRepo.UpdateTicketStatus(ctx, ticketID, userID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to execute status update: %w", err)
	}

	return updatedTicket, nil
}
