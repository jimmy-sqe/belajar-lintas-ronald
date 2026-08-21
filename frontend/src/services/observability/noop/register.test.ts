import { register } from './register';

describe('noop register — observability noop', () => {
  it('resolves without throwing', async () => {
    await expect(register()).resolves.toBeUndefined();
  });
});
