-- Pre-seeded test accounts for TodoApp preview-env (per PRD §4).
-- INTERNAL — do not use in production.
-- Plaintext passwords: "Demo1234!" (both accounts). Documented in
-- docs/todo-app/README.md.

INSERT INTO users (id, email, password_hash, name, created_by, modified_by) VALUES
  ('11111111-1111-1111-1111-111111111111',
   'demo1@example.com',
   '$2a$12$L5naLzmWDnoV8HhC0e4wmupejixf1uVQE/ptf8zbywYl.KQaibMha',
   'Demo One',
   NULL, NULL),
  ('22222222-2222-2222-2222-222222222222',
   'demo2@example.com',
   '$2a$12$eyQG9XAUzLidlEqKNWVVregaJ3T7nWGJF7tMBOHmLPV.UrbrYJx.e',
   'Demo Two',
   NULL, NULL)
ON CONFLICT (id) DO NOTHING;
