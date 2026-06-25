package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ticket-system/internal/models"
	"time"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
)

type TicketRepository interface {
	CreateTicketTx(ctx context.Context, userID string, ticket *models.Ticket) error
	ListTicketsByUserID(ctx context.Context, userID string) ([]*models.Ticket, error)
	GetTicketByIDAndUserID(ctx context.Context, ticketID int64, userID string) (*models.Ticket, error)
	UpdateTicketStatus(ctx context.Context, ticketID int64, userID string, status string) (*models.Ticket, error)
}

type postgresTicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &postgresTicketRepository{db: db}
}

func (r *postgresTicketRepository) CreateTicketTx(ctx context.Context, userID string, ticket *models.Ticket) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	ticketQuery := `
		INSERT INTO tickets (name, "desc", status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err = tx.QueryRowContext(ctx, ticketQuery, ticket.Name, ticket.Desc, ticket.Status, ticket.CreatedAt, ticket.UpdatedAt).Scan(&ticket.ID)
	if err != nil {
		return fmt.Errorf("failed to insert ticket: %w", err)
	}

	junctionQuery := `
		INSERT INTO user_tickets (user_id, ticket_id)
		VALUES ($1, $2)`

	_, err = tx.ExecContext(ctx, junctionQuery, userID, ticket.ID)
	if err != nil {
		return fmt.Errorf("failed to associate user and ticket: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *postgresTicketRepository) ListTicketsByUserID(ctx context.Context, userID string) ([]*models.Ticket, error) {
	query := `
		SELECT t.id, t.name, t.desc, t.status, t.created_at, t.updated_at
		FROM tickets t
		JOIN user_tickets ut ON t.id = ut.ticket_id
		WHERE ut.user_id = $1
		ORDER BY t.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}
	defer rows.Close()

	tickets := []*models.Ticket{}
	for rows.Next() {
		var t models.Ticket
		err := rows.Scan(&t.ID, &t.Name, &t.Desc, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket row: %w", err)
		}
		tickets = append(tickets, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading ticket rows: %w", err)
	}

	return tickets, nil
}

func (r *postgresTicketRepository) GetTicketByIDAndUserID(ctx context.Context, ticketID int64, userID string) (*models.Ticket, error) {
	query := `
		SELECT t.id, t.name, t.desc, t.status, t.created_at, t.updated_at
		FROM tickets t
		JOIN user_tickets ut ON t.id = ut.ticket_id
		WHERE t.id = $1 AND ut.user_id = $2`

	row := r.db.QueryRowContext(ctx, query, ticketID, userID)

	var t models.Ticket
	err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by id: %w", err)
	}

	return &t, nil
}

func (r *postgresTicketRepository) UpdateTicketStatus(ctx context.Context, ticketID int64, userID string, status string) (*models.Ticket, error) {
	ticket, err := r.GetTicketByIDAndUserID(ctx, ticketID, userID)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE tickets
		SET status = $1, updated_at = $2
		WHERE id = $3`

	now := time.Now()
	_, err = r.db.ExecContext(ctx, query, status, now, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket status in DB: %w", err)
	}

	ticket.Status = status
	ticket.UpdatedAt = now
	return ticket, nil
}
