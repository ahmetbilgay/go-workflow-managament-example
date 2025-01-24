package services

import (
	"context"
	"errors"
	"go-workflow/models"
	"go-workflow/websocket"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowService struct {
	db  *mongo.Database
	hub *websocket.Hub
}

func NewWorkflowService(db *mongo.Database, hub *websocket.Hub) *WorkflowService {
	return &WorkflowService{
		db:  db,
		hub: hub,
	}
}

func (s *WorkflowService) CreateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	workflow.CreatedAt = time.Now()
	workflow.UpdatedAt = time.Now()
	workflow.Status = models.StatusPending

	// Context'i başlat
	workflow.Context = models.ProcessContext{
		CurrentData: make(map[string]interface{}),
		ProcessData: make(map[string]interface{}),
		SharedData:  make(map[string]interface{}),
		StepResults: make(map[string]interface{}),
	}

	// İlk adımı current step olarak ayarla
	if len(workflow.Steps) > 0 {
		workflow.CurrentStep = workflow.Steps[0].ID
	}

	result, err := s.db.Collection("workflows").InsertOne(ctx, workflow)
	if err != nil {
		return err
	}

	workflow.ID = result.InsertedID.(primitive.ObjectID)

	// İlk adımın atanmış kullanıcısına bildirim gönder
	if len(workflow.Steps) > 0 {
		firstStep := workflow.Steps[0]
		notification := &models.Notification{
			Type:       models.NotificationTypeWorkflow,
			Title:      "Yeni Workflow Görevi",
			Message:    "Size yeni bir workflow görevi atandı: " + workflow.Name,
			UserID:     firstStep.AssignedTo,
			WorkflowID: workflow.ID,
			StepID:     firstStep.ID,
			Status:     models.NotificationStatusUnread,
			CreatedAt:  time.Now(),
		}

		_, err = s.db.Collection("notifications").InsertOne(ctx, notification)
		if err != nil {
			return err
		}

		s.hub.SendToUser(firstStep.AssignedTo, websocket.Message{
			Type:    "new_workflow",
			Payload: notification,
		})
	}

	return nil
}

func (s *WorkflowService) ProcessStep(ctx context.Context, workflowID, stepID primitive.ObjectID, action string, data map[string]interface{}) error {
	var workflow models.Workflow
	err := s.db.Collection("workflows").FindOne(ctx, primitive.M{"_id": workflowID}).Decode(&workflow)
	if err != nil {
		return err
	}

	var currentStep *models.WorkflowStep
	for i, step := range workflow.Steps {
		if step.ID == stepID {
			currentStep = &workflow.Steps[i]
			break
		}
	}

	if currentStep == nil {
		return errors.New("adım bulunamadı")
	}

	// Gerekli verilerin kontrolü
	if len(currentStep.RequiredData) > 0 {
		for _, required := range currentStep.RequiredData {
			if _, exists := workflow.Context.SharedData[required]; !exists {
				return errors.New("gerekli veri eksik: " + required)
			}
		}
	}

	// Gelen veriyi current data'ya kaydet
	workflow.Context.CurrentData = data

	switch action {
	case "approve":
		currentStep.Status = models.StatusApproved

		// Step verilerini kaydet
		currentStep.StepData = data

		// Sonuç varsa kaydet
		if currentStep.ResultType != "" {
			result := models.WorkflowResult{
				Type:      currentStep.ResultType,
				Data:      data,
				CreatedAt: time.Now(),
			}
			workflow.Results = append(workflow.Results, result)

			// Sonucu step results'a da ekle
			workflow.Context.StepResults[stepID.Hex()] = data
		}

		// Paylaşılan verileri güncelle
		for key, value := range data {
			workflow.Context.SharedData[key] = value
		}

		// Sonraki adımları belirle ve bildirim gönder
		if len(currentStep.NextSteps) > 0 {
			for _, nextStepID := range currentStep.NextSteps {
				var nextStep *models.WorkflowStep
				for i, step := range workflow.Steps {
					if step.ID == nextStepID {
						nextStep = &workflow.Steps[i]
						break
					}
				}

				if nextStep != nil {
					workflow.CurrentStep = nextStep.ID

					notification := &models.Notification{
						Type:       models.NotificationTypeWorkflow,
						Title:      "Yeni Workflow Adımı",
						Message:    "Size yeni bir workflow adımı atandı: " + workflow.Name,
						UserID:     nextStep.AssignedTo,
						WorkflowID: workflow.ID,
						StepID:     nextStep.ID,
						Status:     models.NotificationStatusUnread,
						CreatedAt:  time.Now(),
					}

					_, err = s.db.Collection("notifications").InsertOne(ctx, notification)
					if err != nil {
						return err
					}

					s.hub.SendToUser(nextStep.AssignedTo, websocket.Message{
						Type:    "new_step",
						Payload: notification,
					})
				}
			}
		} else {
			workflow.Status = models.StatusCompleted
		}

	case "reject":
		currentStep.Status = models.StatusRejected
		workflow.Status = models.StatusRejected
		currentStep.StepData = data
	default:
		return errors.New("geçersiz işlem")
	}

	workflow.UpdatedAt = time.Now()
	_, err = s.db.Collection("workflows").ReplaceOne(ctx, primitive.M{"_id": workflow.ID}, workflow)
	return err
}

func (s *WorkflowService) GetWorkflow(ctx context.Context, id primitive.ObjectID) (*models.Workflow, error) {
	var workflow models.Workflow
	err := s.db.Collection("workflows").FindOne(ctx, primitive.M{"_id": id}).Decode(&workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (s *WorkflowService) GetUserWorkflows(ctx context.Context, userID primitive.ObjectID) ([]models.Workflow, error) {
	var workflows []models.Workflow
	cursor, err := s.db.Collection("workflows").Find(ctx, primitive.M{
		"steps": primitive.M{
			"$elemMatch": primitive.M{
				"assigned_to": userID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &workflows)
	return workflows, err
}
