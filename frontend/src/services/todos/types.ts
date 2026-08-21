export interface TodoResponse {
  id: string;
  title: string;
  description?: string;
  due_date?: string;
  created_at: string;
  created_by: string;
  modified_at: string;
  modified_by: string;
}

export interface CreateTodoRequest {
  title: string;
  description?: string;
  due_date?: string;
}

export type UpdateTodoRequest = Partial<CreateTodoRequest>;
