package model

import (
	"time"
)

type (
	Users []*User

	User struct {
		ID           int       `json:"id"`
		Username     string    `json:"username"`
		PasswordHash string    `json:"-"`
		Role         string    `json:"role"`
		CreatedAt    time.Time `json:"created_at"`
	}
)

type (
	Processes []*Process

	Process struct {
		ID         int        `json:"id"`
		Title      string     `json:"title"`
		Status     string     `json:"status"`
		StartedAt  *time.Time `json:"started_at"`
		FinishedAt *time.Time `json:"finished_at"`
		Tasks      Tasks      `json:"tasks"`
	}
)

type (
	Tasks []*Task

	Task struct {
		ID             int         `json:"id"`
		ProcessID      int         `json:"process_id"`
		Title          string      `json:"title"`
		Comment        string      `json:"comment"`
		Status         string      `json:"status"`
		ForRole        string      `json:"for_role"`
		StartedAt      *time.Time  `json:"started_at"`
		FinishedAt     *time.Time  `json:"finished_at"`
		PlanHours      int         `json:"plan_hours"`
		IsFileRequired bool        `json:"is_file_required"`
		Attachments    Attachments `json:"attachments"`
		Users          Users       `json:"users"`
		DependsOn      []int       `json:"depends_on"`
	}
)

type (
	Attachments []*Attachment

	Attachment struct {
		ID         int       `json:"id"`
		TaskID     int       `json:"task_id"`
		FilePath   string    `json:"file_path"`
		FileSizeKb int       `json:"file_size_kb"`
		UploadedAt time.Time `json:"uploaded_at"`
	}
)

type (
	Templates []*Template

	Template struct {
		ID          int           `json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Tasks       TaskTemplates `json:"tasks"`
	}
)

type (
	TaskTemplates []*TaskTemplate

	TaskTemplate struct {
		ID             int    `json:"id"`
		TemplateID     int    `json:"template_id"`
		Title          string `json:"title"`
		IsFileRequired bool   `json:"is_file_required"`
		PlanHours      int    `json:"plan_hours"`
		ForRole        string `json:"for_role"`
		DependsOn      []int  `json:"depends_on"`
	}
)

func NewProcessFromTemplate(temp *Template, title string) *Process {
	proc := &Process{
		Title:  title,
		Status: "draft",
		Tasks:  make(Tasks, len(temp.Tasks)),
	}
	return proc
}
