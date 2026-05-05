-- ============================================================
-- Наполнение тестовыми данными
-- ============================================================

-- 1. Пользователи
-- Пароль для всех: password123 (bcrypt hash)
INSERT INTO users (username, password_hash, role) VALUES
    ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'),
    ('ivanov', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'manager'),
    ('petrova', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
    ('sidorov', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker');

-- 2. Шаблоны процессов
INSERT INTO templates (name, description) VALUES
    ('Процесс найма', 'Стандартный процесс найма нового сотрудника'),
    ('Онбординг', 'Процесс ввода сотрудника в должность'),
    ('Заказ оборудования', 'Процесс закупки и выдачи рабочего оборудования');

-- 3. Задачи для шаблона "Процесс найма" (id=1)
INSERT INTO template_tasks (template_id, title, role, is_file_required, plan_done_hours) VALUES
    (1, 'Собеседование', 'manager', false, 2),
    (1, 'Проверка документов', 'worker', true, 1),
    (1, 'Согласование оффера', 'manager', false, 1),
    (1, 'Оформление', 'admin', true, 3);

-- Зависимости для шаблона "Процесс найма"
INSERT INTO template_dependencies (template_task_id, depends_on_id) VALUES
    (2, 1),  -- Проверка документов после Собеседования
    (3, 2),  -- Согласование оффера после Проверки документов
    (4, 3);  -- Оформление после Согласования оффера

-- 4. Задачи для шаблона "Онбординг" (id=2)
INSERT INTO template_tasks (template_id, title, role, is_file_required, plan_done_hours) VALUES
    (2, 'Подготовка рабочего места', 'admin', false, 4),
    (2, 'Знакомство с командой', 'manager', false, 1),
    (2, 'Получение доступов', 'admin', false, 2),
    (2, 'Обучение', 'worker', false, 8);

-- Зависимости для шаблона "Онбординг"
INSERT INTO template_dependencies (template_task_id, depends_on_id) VALUES
    (6, 5),  -- Знакомство с командой после Подготовки рабочего места
    (7, 5),  -- Получение доступов после Подготовки рабочего места
    (8, 6),  -- Обучение после Знакомства с командой
    (8, 7);  -- Обучение после Получения доступов

-- 5. Задачи для шаблона "Заказ оборудования" (id=3)
INSERT INTO template_tasks (template_id, title, role, is_file_required, plan_done_hours) VALUES
    (3, 'Сбор заявок', 'worker', false, 3),
    (3, 'Согласование бюджета', 'manager', false, 2),
    (3, 'Закупка', 'admin', true, 5),
    (3, 'Выдача сотрудникам', 'worker', false, 1);

-- Зависимости для шаблона "Заказ оборудования"
INSERT INTO template_dependencies (template_task_id, depends_on_id) VALUES
    (10, 9),   -- Согласование бюджета после Сбора заявок
    (11, 10),  -- Закупка после Согласования бюджета
    (12, 11);  -- Выдача после Закупки

-- 6. Тестовые процессы
INSERT INTO processes (template_id, title, status, start_date) VALUES
    (1, 'Найм Иванова И.И.', 'in_progress', NOW() - INTERVAL '3 days'),
    (1, 'Найм Петровой А.С.', 'draft', NOW() - INTERVAL '1 day'),
    (2, 'Онбординг Сидорова', 'done', NOW() - INTERVAL '14 days'),
    (3, 'Закупка ноутбуков', 'in_progress', NOW() - INTERVAL '10 days');

-- 7. Задачи для процесса "Найм Иванова" (process_id=1)
INSERT INTO tasks (process_id, title, status, role, is_file_required, plan_done_hours, started_at, finished_at) VALUES
    (1, 'Собеседование', 'done', 'manager', false, 2, NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days'),
    (1, 'Проверка документов', 'in_progress', 'worker', true, 1, NOW() - INTERVAL '2 days', NULL),
    (1, 'Согласование оффера', 'pending', 'manager', false, 1, NULL, NULL),
    (1, 'Оформление', 'pending', 'admin', true, 3, NULL, NULL);

-- Зависимости задач для процесса 1
INSERT INTO task_dependencies (task_id, depends_on_id) VALUES
    (2, 1),
    (3, 2),
    (4, 3);

-- 8. Задачи для процесса "Найм Петровой" (process_id=2)
INSERT INTO tasks (process_id, title, status, role, is_file_required, plan_done_hours) VALUES
    (2, 'Собеседование', 'pending', 'manager', false, 2),
    (2, 'Проверка документов', 'pending', 'worker', true, 1),
    (2, 'Согласование оффера', 'pending', 'manager', false, 1),
    (2, 'Оформление', 'pending', 'admin', true, 3);

-- Зависимости задач для процесса 2
INSERT INTO task_dependencies (task_id, depends_on_id) VALUES
    (6, 5),
    (7, 6),
    (8, 7);

-- 9. Задачи для процесса "Онбординг Сидорова" (process_id=3)
INSERT INTO tasks (process_id, title, status, role, is_file_required, plan_done_hours, started_at, finished_at) VALUES
    (3, 'Подготовка рабочего места', 'done', 'admin', false, 4, NOW() - INTERVAL '14 days', NOW() - INTERVAL '10 days'),
    (3, 'Знакомство с командой', 'done', 'manager', false, 1, NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),
    (3, 'Получение доступов', 'done', 'admin', false, 2, NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days'),
    (3, 'Обучение', 'done', 'worker', false, 8, NOW() - INTERVAL '8 days', NOW() - INTERVAL '1 day');

-- Зависимости задач для процесса 3
INSERT INTO task_dependencies (task_id, depends_on_id) VALUES
    (10, 9),
    (11, 9),
    (12, 10),
    (12, 11);

-- 10. Задачи для процесса "Закупка ноутбуков" (process_id=4)
INSERT INTO tasks (process_id, title, status, role, is_file_required, plan_done_hours, started_at, finished_at) VALUES
    (4, 'Сбор заявок', 'done', 'worker', false, 3, NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days'),
    (4, 'Согласование бюджета', 'done', 'manager', false, 2, NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days'),
    (4, 'Закупка', 'in_progress', 'admin', true, 5, NOW() - INTERVAL '5 days', NULL),
    (4, 'Выдача сотрудникам', 'pending', 'worker', false, 1, NULL, NULL);

-- Зависимости задач для процесса 4
INSERT INTO task_dependencies (task_id, depends_on_id) VALUES
    (14, 13),
    (15, 14),
    (16, 15);

-- 11. Назначения пользователей на задачи
INSERT INTO task_assignments (user_id, task_id) VALUES
    (2, 1),  -- ivanov на Собеседование (процесс 1)
    (3, 2),  -- petrova на Проверку документов
    (4, 13), -- sidorov на Сбор заявок (процесс 4)
    (3, 14); -- petrova на Согласование бюджета