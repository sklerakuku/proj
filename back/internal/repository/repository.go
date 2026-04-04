package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sklerakuku/tracker-web/internal/model"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// USER
func (r *Repository) CreateUser(ctx context.Context, username, passwordHash, role string) (*model.User, error) {
	var user model.User
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)
	RETURNING id, username, role, created_at`

	err := r.db.QueryRow(ctx, query, username, passwordHash, role).Scan(
		&user.ID, &user.Username, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE username = $1`
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// TEMPLATE
func (r *Repository) CreateTemplate(ctx context.Context, name, description string) (*model.Template, error) {
	var template model.Template
	query := `INSERT INTO templates (name, description) VALUES ($1, $2) 
	RETURNING id, name, description`
	err := r.db.QueryRow(ctx, query, name, description).Scan(&template.ID, &template.Name, &template.Description)
	if err != nil {
		return nil, err
	}
	template.Tasks = model.TaskTemplates{}
	return &template, nil
}

func (r *Repository) GetAllTemplates(ctx context.Context) (model.Templates, error) {
	rows, err := r.db.Query(ctx, `SELECT 8 FROM templates ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates model.Templates
	for rows.Next() {
		var t model.Template
		err := rows.Scan(&t.ID, &t.Name, &t.Description)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, nil
}

func (r *Repository) GetTemplateByID(ctx context.Context, id int) (*model.Template, error) {
	var template model.Template
	query := `SELECT id, name, description FROM templates WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&template.ID, &template.Name, &template.Description)
	if err != nil {
		return nil, err
	}

	// get tasks
	tasksQuery := `SELECT id, template_id, title, is_file_required, plan_done_hours, role 
	               FROM template_tasks WHERE template_id = $1`
	taskRows, err := r.db.Query(ctx, tasksQuery, id)
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()

	for taskRows.Next() {
		var task model.TaskTemplate
		err := taskRows.Scan(&task.ID, &task.TemplateID, &task.Title, &task.IsFileRequired, &task.PlanHours, &task.ForRole)
		if err != nil {
			return nil, err
		}
		template.Tasks = append(template.Tasks, &task)
	}

	// get dependencies
	for _, task := range template.Tasks {
		depQuery := `SELECT depends_on_id FROM template_dependencies WHERE template_task_id = $1`
		depRows, err := r.db.Query(ctx, depQuery, task.ID)
		if err != nil {
			return nil, err
		}
		defer depRows.Close()

		for depRows.Next() {
			var depID int
			depRows.Scan(&depID)
			task.DependsOn = append(task.DependsOn, depID)
		}
	}

	return &template, nil
}

func (r *Repository) CreateTemplateTask(ctx context.Context, templateID int, title, role string, isFileRequired bool, planHours int) (int, error) {
	var id int
	query := `INSERT INTO template_tasks (template_id, title, role, is_file_required, plan_done_hours) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(ctx, query, templateID, title, role, isFileRequired, planHours).Scan(&id)
	return id, err
}

func (r *Repository) AddTemplateDependency(ctx context.Context, taskID, dependsOnID int) error {
	query := `INSERT INTO template_dependencies (template_task_id, depends_on_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, taskID, dependsOnID)
	return err
}

// Process methods
func (r *Repository) CreateProcess(ctx context.Context, templateID int, title, status string) (*model.Process, error) {
	var process model.Process
	query := `INSERT INTO processes (template_id, title, status) VALUES ($1, $2, $3) 
	          RETURNING id, title, status`
	err := r.db.QueryRow(ctx, query, templateID, title, status).Scan(&process.ID, &process.Title, &process.Status)
	if err != nil {
		return nil, err
	}
	return &process, nil
}

func (r *Repository) CreateTask(ctx context.Context, processID int, title, role string, isFileRequired bool, planHours int) (int, error) {
	var id int
	query := `INSERT INTO tasks (process_id, title, role, is_file_required, plan_done_hours, status) 
	          VALUES ($1, $2, $3, $4, $5, 'pending') RETURNING id`
	err := r.db.QueryRow(ctx, query, processID, title, role, isFileRequired, planHours).Scan(&id)
	return id, err
}

func (r *Repository) AddTaskDependency(ctx context.Context, taskID, dependsOnID int) error {
	query := `INSERT INTO task_dependencies (task_id, depends_on_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, taskID, dependsOnID)
	return err
}

func (r *Repository) GetProcessByID(ctx context.Context, id int) (*model.Process, error) {
	var process model.Process
	query := `SELECT id, title, status, start_date, finished_at FROM processes WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&process.ID, &process.Title, &process.Status, &process.StartedAt, &process.FinishedAt)
	if err != nil {
		return nil, err
	}

	// get tasks
	tasksQuery := `SELECT id, title, comment, status, role, started_at, finished_at, plan_done_hours, is_file_required 
	               FROM tasks WHERE process_id = $1`
	taskRows, err := r.db.Query(ctx, tasksQuery, id)
	if err != nil {
		return nil, err
	}
	defer taskRows.Close()

	for taskRows.Next() {
		var task model.Task
		err := taskRows.Scan(&task.ID, &task.Title, &task.Comment, &task.Status, &task.ForRole, &task.StartedAt, &task.FinishedAt, &task.PlanHours, &task.IsFileRequired)
		if err != nil {
			return nil, err
		}
		task.ProcessID = process.ID
		process.Tasks = append(process.Tasks, &task)
	}

	// get dependencies for each task
	for _, task := range process.Tasks {
		depQuery := `SELECT depends_on_id FROM task_dependencies WHERE task_id = $1`
		depRows, err := r.db.Query(ctx, depQuery, task.ID)
		if err != nil {
			return nil, err
		}
		defer depRows.Close()

		for depRows.Next() {
			var depID int
			depRows.Scan(&depID)
			task.DependsOn = append(task.DependsOn, depID)
		}
	}

	return &process, nil
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, taskID int, status string) error {
	var startedAt, finishedAt interface{}
	now := time.Now()

	if status == "in_progress" {
		startedAt = now
		finishedAt = nil
	} else if status == "done" {
		finishedAt = now
	} else {
		startedAt = nil
		finishedAt = nil
	}

	query := `UPDATE tasks SET status = $1, started_at = COALESCE($2, started_at), finished_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, status, startedAt, finishedAt, taskID)
	return err
}

func (r *Repository) GetTaskDependencies(ctx context.Context, taskID int) ([]int, error) {
	rows, err := r.db.Query(ctx, `SELECT depends_on_id FROM task_dependencies WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []int
	for rows.Next() {
		var depID int
		rows.Scan(&depID)
		deps = append(deps, depID)
	}
	return deps, nil
}

func (r *Repository) CheckTasksStatus(ctx context.Context, taskIDs []int, status string) (bool, error) {
	if len(taskIDs) == 0 {
		return true, nil
	}

	query := `SELECT COUNT(*) FROM tasks WHERE id = ANY($1) AND status != $2`
	var count int
	err := r.db.QueryRow(ctx, query, taskIDs, status).Scan(&count)
	return count == 0, err
}
