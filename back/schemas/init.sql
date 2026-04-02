-- 1. Users Table
-- Stores credentials and roles for system access.
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- Never store plain text passwords
    role VARCHAR(20) NOT NULL DEFAULT 'worker',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Templates
-- Blueprints for processes.
CREATE TABLE templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- 3. Template Tasks
-- Specific task definitions within a template.
CREATE TABLE template_tasks (
    id SERIAL PRIMARY KEY,
    template_id INTEGER REFERENCES templates(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    role VARCHAR(50), 
    is_file_required BOOLEAN DEFAULT FALSE,
    plan_done_hours INTEGER
);

-- 4. Processes
-- An instance of a template being executed.
CREATE TABLE processes (
    id SERIAL PRIMARY KEY,
    template_id INTEGER REFERENCES templates(id),
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    start_date TIMESTAMP,
    finished_at TIMESTAMP
);

-- 5. Tasks
-- Actual tasks generated from template_tasks for a specific process.
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    process_id INTEGER REFERENCES processes(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    is_file_required BOOLEAN DEFAULT FALSE,
    role VARCHAR(50), 
    plan_done_hours INTEGER
);

-- 6. Task Assignments (User-Task Relationship)
-- Tracks which user is working on which task.
CREATE TABLE task_assignments (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, task_id)
);

-- 7. Dependencies (Unified Table)
CREATE TABLE template_dependencies (
    template_task_id INTEGER REFERENCES template_tasks(id) ON DELETE CASCADE,
    depends_on_id INTEGER REFERENCES template_tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (template_task_id, depends_on_id),
    CHECK (template_task_id != depends_on_id)
);

CREATE TABLE task_dependencies (
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    depends_on_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, depends_on_id),
    CHECK (task_id != depends_on_id)
);

-- 8. Attachments
CREATE TABLE attachments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_size_kb INTEGER,
    upload_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
