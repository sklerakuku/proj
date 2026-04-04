package service

import (
	"context"
	"errors"

	"github.com/sklerakuku/tracker-web/internal/model"
	"github.com/sklerakuku/tracker-web/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// AUTH
func (s *Service) Register(ctx context.Context, username, password, role string) (*model.User, error) {
	if role == "" {
		role = "worker"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return s.repo.CreateUser(ctx, username, string(hash), role)
}

func (s *Service) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// Template
func (s *Service) CreateTemplate(ctx context.Context, name, description string, tasks []model.TaskTemplate) (*model.Template, error) {
	template, err := s.repo.CreateTemplate(ctx, name, description)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		taskID, err := s.repo.CreateTemplateTask(ctx, template.ID, task.Title, task.ForRole, task.IsFileRequired, task.PlanHours)
		if err != nil {
			return nil, err
		}

		for _, depID := range task.DependsOn {
			err = s.repo.AddTemplateDependency(ctx, taskID, depID)
			if err != nil {
				return nil, err
			}
		}
	}

	return s.repo.GetTemplateByID(ctx, template.ID)
}

func (s *Service) GetAllTemplates(ctx context.Context) (model.Templates, error) {
	return s.repo.GetAllTemplates(ctx)
}

func (s *Service) GetTemplate(ctx context.Context, id int) (*model.Template, error) {
	return s.repo.GetTemplateByID(ctx, id)
}

// Process
func (s *Service) CreateProcessFromTemplate(ctx context.Context, templateID int, title string) (*model.Process, error) {
	template, err := s.repo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	process, err := s.repo.CreateProcess(ctx, templateID, title, "draft")
	if err != nil {
		return nil, err
	}

	taskIDMap := make(map[int]int) // template_task_id -> task_id

	for _, tt := range template.Tasks {
		taskID, err := s.repo.CreateTask(ctx, process.ID, tt.Title, tt.ForRole, tt.IsFileRequired, tt.PlanHours)
		if err != nil {
			return nil, err
		}
		taskIDMap[tt.ID] = taskID
	}

	for _, tt := range template.Tasks {
		for _, depTemplateID := range tt.DependsOn {
			if taskID, ok := taskIDMap[tt.ID]; ok {
				if depTaskID, ok := taskIDMap[depTemplateID]; ok {
					err = s.repo.AddTaskDependency(ctx, taskID, depTaskID)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return s.repo.GetProcessByID(ctx, process.ID)
}

func (s *Service) GetProcess(ctx context.Context, id int) (*model.Process, error) {
	return s.repo.GetProcessByID(ctx, id)
}

func (s *Service) UpdateTaskStatus(ctx context.Context, taskID int, status string) error {
	if status != "pending" && status != "in_progress" && status != "done" {
		return errors.New("invalid status")
	}

	if status == "done" {
		deps, err := s.repo.GetTaskDependencies(ctx, taskID)
		if err != nil {
			return err
		}

		allDone, err := s.repo.CheckTasksStatus(ctx, deps, "done")
		if err != nil {
			return err
		}

		if !allDone {
			return errors.New("cannot complete task: dependencies not finished")
		}
	}

	return s.repo.UpdateTaskStatus(ctx, taskID, status)
}
