CREATE TYPE task_status_enum AS ENUM ('pending', 'completed');

CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  
  title TEXT NOT NULL,
  status task_status_enum NOT NULL,
  description TEXT,
  
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  
  due_date TIMESTAMP,
  
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NULL
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);