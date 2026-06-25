package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ticket-system/internal/middleware"
	"ticket-system/internal/service"

	"github.com/go-chi/chi/v5"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

type createTicketRequest struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Description string `json:"description"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	name := req.Name
	if name == "" {
		name = req.Title
	}
	desc := req.Desc
	if desc == "" {
		desc = req.Description
	}

	if name == "" {
		RespondWithError(w, http.StatusBadRequest, "Ticket name or title is required")
		return
	}

	ticket, err := h.ticketService.CreateTicket(r.Context(), userID, name, desc)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create ticket")
		return
	}

	RespondWithJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tickets, err := h.ticketService.ListTickets(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch tickets")
		return
	}

	RespondWithJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid ticket ID format")
		return
	}

	ticket, err := h.ticketService.GetTicket(r.Context(), ticketID, userID)
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			RespondWithError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve ticket")
		return
	}

	RespondWithJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid ticket ID format")
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Status == "" {
		RespondWithError(w, http.StatusBadRequest, "Status is required")
		return
	}

	ticket, err := h.ticketService.UpdateTicketStatus(r.Context(), ticketID, userID, req.Status)
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			RespondWithError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) || errors.Is(err, service.ErrInvalidStatus) {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Failed to update ticket status")
		return
	}

	RespondWithJSON(w, http.StatusOK, ticket)
}
