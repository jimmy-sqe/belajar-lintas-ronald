// Stable internal names wrapping kubb-generated hooks. Page code consumes
// only these — kubb output rename doesn't propagate.

import {
  useListTodos as kubbUseListTodos,
  useGetTodo as kubbUseGetTodo,
  useCreateTodo as kubbUseCreateTodo,
  useUpdateTodo as kubbUseUpdateTodo,
  useDeleteTodo as kubbUseDeleteTodo
} from '@/openapi';

export const useTodos = kubbUseListTodos;
export const useTodo = kubbUseGetTodo;
export const useCreateTodo = kubbUseCreateTodo;
export const useUpdateTodo = kubbUseUpdateTodo;
export const useDeleteTodo = kubbUseDeleteTodo;

export type { TodoResponse, CreateTodoRequest, UpdateTodoRequest } from './types';
