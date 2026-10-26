import { test } from 'node:test';
import assert from 'node:assert/strict';
import { newAdministrator } from '../../src/domain/access/administrator.js';
import { ErrInvalidCredentials } from '../../src/domain/access/errors.js';

test('newAdministrator never stores the plaintext password', () => {
  const admin = newAdministrator('alice', 'correct-horse-battery-staple');
  assert.notEqual(admin.passwordHash, 'correct-horse-battery-staple');
  assert.ok(admin.passwordHash.startsWith('$2'));
});

test('authenticate succeeds on the correct password', () => {
  const admin = newAdministrator('alice', 'correct-horse-battery-staple');
  assert.doesNotThrow(() => admin.authenticate('correct-horse-battery-staple'));
});

test('authenticate rejects a wrong password with the generic error', () => {
  const admin = newAdministrator('alice', 'correct-horse-battery-staple');
  assert.throws(() => admin.authenticate('wrong-password'), ErrInvalidCredentials);
});
