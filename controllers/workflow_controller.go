package controllers

import (
	"encoding/json"
	"go-workflow/models"
	"go-workflow/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkflowController struct {
	workflowService *services.WorkflowService
}

func NewWorkflowController(workflowService *services.WorkflowService) *WorkflowController {
	return &WorkflowController{
		workflowService: workflowService,
	}
}

func (c *WorkflowController) RegisterRoutes(r chi.Router) {
	r.Post("/workflows", c.CreateWorkflow)
	r.Get("/workflows/{id}", c.GetWorkflow)
	r.Post("/workflows/{id}/steps/{stepId}/process", c.ProcessStep)
	r.Get("/workflows/user/{userId}", c.GetUserWorkflows)
}

func (c *WorkflowController) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.workflowService.CreateWorkflow(r.Context(), &workflow); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workflow)
}

func (c *WorkflowController) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Geçersiz ID formatı", http.StatusBadRequest)
		return
	}

	workflow, err := c.workflowService.GetWorkflow(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(workflow)
}

func (c *WorkflowController) ProcessStep(w http.ResponseWriter, r *http.Request) {
	workflowID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Geçersiz workflow ID formatı", http.StatusBadRequest)
		return
	}

	stepID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "stepId"))
	if err != nil {
		http.Error(w, "Geçersiz step ID formatı", http.StatusBadRequest)
		return
	}

	var data struct {
		Action string                 `json:"action"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.workflowService.ProcessStep(r.Context(), workflowID, stepID, data.Action, data.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *WorkflowController) GetUserWorkflows(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Geçersiz user ID formatı", http.StatusBadRequest)
		return
	}

	workflows, err := c.workflowService.GetUserWorkflows(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(workflows)
}
