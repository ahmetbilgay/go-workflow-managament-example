package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationType string

const (
	NotificationTypeWorkflow NotificationType = "workflow"
	NotificationTypeTask     NotificationType = "task"
	NotificationTypeSystem   NotificationType = "system"
)

type NotificationStatus string

const (
	NotificationStatusUnread   NotificationStatus = "unread"
	NotificationStatusRead     NotificationStatus = "read"
	NotificationStatusArchived NotificationStatus = "archived"
)

type Notification struct {
	ID         primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Type       NotificationType       `json:"type" bson:"type"`
	Title      string                 `json:"title" bson:"title"`
	Message    string                 `json:"message" bson:"message"`
	UserID     primitive.ObjectID     `json:"user_id" bson:"user_id"`
	WorkflowID primitive.ObjectID     `json:"workflow_id,omitempty" bson:"workflow_id,omitempty"`
	StepID     primitive.ObjectID     `json:"step_id,omitempty" bson:"step_id,omitempty"`
	Status     NotificationStatus     `json:"status" bson:"status"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
	ReadAt     *time.Time             `json:"read_at,omitempty" bson:"read_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
