// boilerplate:axis=sample-app option=todo-app START
import { todosHandlers } from './todos';
// boilerplate:axis=sample-app option=todo-app END
// boilerplate:axis=auth option=jwt-refresh START
import * as jwtMocks from './auth/jwt-refresh';
// boilerplate:axis=auth option=jwt-refresh END
// boilerplate:axis=auth option=opaque-session START
import * as opaqueMocks from './auth/opaque-session';
// boilerplate:axis=auth option=opaque-session END

export const handlers = [
  // boilerplate:axis=sample-app option=todo-app START
  ...todosHandlers,
  // boilerplate:axis=sample-app option=todo-app END
  // boilerplate:axis=auth option=jwt-refresh START
  ...jwtMocks.authHandlers,
  // boilerplate:axis=auth option=jwt-refresh END
  // boilerplate:axis=auth option=opaque-session START
  ...opaqueMocks.authHandlers
  // boilerplate:axis=auth option=opaque-session END
];
