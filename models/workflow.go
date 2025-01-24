package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkflowStatus string
type StepType string
type ResultType string

const (
	StatusPending   WorkflowStatus = "pending"
	StatusApproved  WorkflowStatus = "approved"
	StatusRejected  WorkflowStatus = "rejected"
	StatusCompleted WorkflowStatus = "completed"

	StepTypeApproval StepType = "approval"
	StepTypeTask     StepType = "task"
	StepTypeDecision StepType = "decision"
	StepTypeProcess  StepType = "process"

	ResultTypeInvoice      ResultType = "invoice"
	ResultTypeDocument     ResultType = "document"
	ResultTypeReport       ResultType = "report"
	ResultTypeNotification ResultType = "notification"
)

type WorkflowStep struct {
	ID           primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Type         StepType               `json:"type" bson:"type"`
	Title        string                 `json:"title" bson:"title"`
	Description  string                 `json:"description" bson:"description"`
	AssignedTo   primitive.ObjectID     `json:"assigned_to" bson:"assigned_to"`
	Status       WorkflowStatus         `json:"status" bson:"status"`
	NextSteps    []primitive.ObjectID   `json:"next_steps" bson:"next_steps"`
	Conditions   map[string]interface{} `json:"conditions,omitempty" bson:"conditions,omitempty"`
	ProcessData  map[string]interface{} `json:"process_data,omitempty" bson:"process_data,omitempty"`
	ResultType   ResultType             `json:"result_type,omitempty" bson:"result_type,omitempty"`
	StepData     map[string]interface{} `json:"step_data,omitempty" bson:"step_data,omitempty"`
	RequiredData []string               `json:"required_data,omitempty" bson:"required_data,omitempty"`
}

type WorkflowResult struct {
	Type      ResultType             `json:"type" bson:"type"`
	Data      map[string]interface{} `json:"data" bson:"data"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
}

type ProcessContext struct {
	CurrentData map[string]interface{} `json:"current_data" bson:"current_data"`
	ProcessData map[string]interface{} `json:"process_data" bson:"process_data"`
	SharedData  map[string]interface{} `json:"shared_data" bson:"shared_data"`
	StepResults map[string]interface{} `json:"step_results" bson:"step_results"`
}

type Workflow struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Type        string                 `json:"type" bson:"type"`
	CreatedBy   primitive.ObjectID     `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
	Status      WorkflowStatus         `json:"status" bson:"status"`
	Steps       []WorkflowStep         `json:"steps" bson:"steps"`
	CurrentStep primitive.ObjectID     `json:"current_step" bson:"current_step"`
	Results     []WorkflowResult       `json:"results,omitempty" bson:"results,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Context     ProcessContext         `json:"context" bson:"context"`
}
