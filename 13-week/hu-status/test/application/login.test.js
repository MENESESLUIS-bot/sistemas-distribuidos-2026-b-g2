import { test } from 'node:test';
import assert from 'node:assert/strict';
import { Login } from '../../src/application/usecase/login.js';
import { newAdministrator } from '../../src/domain/access/administrator.js';
import { ErrAdministratorNotFound, ErrInvalidCredentials } from '../../src/domain/access/errors.js';

function fakeTokens() {
  return {
    async issue(administratorId) {
      return { token: `fake-token-for-${administratorId}`, expiresIn: 3600 };
    },
  };
}

test('login succeeds and returns a non-empty token', async () => {
  const admin = newAdministrator('alice', 'correct-horse-battery-staple');
  admin.id = 'admin-1';
  const administrators = { findByUsername: async () => admin };
  const login = new Login({ administrators, tokens: fakeTokens() });

  const result = await login.execute('alice', 'correct-horse-battery-staple');
  assert.ok(result.accessToken.length > 0);
  assert.equal(result.expiresIn, 3600);
});

test('login with an unknown username maps to the generic invalid-credentials error', async () => {
  const administrators = {
    findByUsername: async () => {
      throw new ErrAdministratorNotFound();
    },
  };
  const login = new Login({ administrators, tokens: fakeTokens() });

  await assert.rejects(() => login.execute('ghost', 'whatever'), ErrInvalidCredentials);
});

test('login with a wrong password maps to the same generic invalid-credentials error', async () => {
  const admin = newAdministrator('alice', 'correct-horse-battery-staple');
  admin.id = 'admin-1';
  const administrators = { findByUsername: async () => admin };
  const login = new Login({ administrators, tokens: fakeTokens() });

  await assert.rejects(() => login.execute('alice', 'wrong-password'), ErrInvalidCredentials);
});
