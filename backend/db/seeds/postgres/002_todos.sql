INSERT INTO todos (id, title, description, due_date, created_by, modified_by) VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
   'Buy groceries', 'Milk, bread, eggs', '2026-05-20T10:00:00Z',
   '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
   'Review plugin specs', NULL, NULL,
   '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3',
   'Pay rent', 'Wire transfer to landlord', '2026-05-31T23:59:00Z',
   '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4',
   'Prepare demo', 'Run TodoApp through preview-env', '2026-05-15T14:00:00Z',
   '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa5',
   'Update documentation', NULL, NULL,
   '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222')
ON CONFLICT (id) DO NOTHING;
