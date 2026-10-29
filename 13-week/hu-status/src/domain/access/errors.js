// Port-level domain errors — use cases depend on these, never on an
// infrastructure adapter's own error type
// (library-docs/05-architecture/hexagonal-architecture.md, Dependency Rule).

// ErrInvalidCredentials is the single, generic error for any login failure — INV-002:
// login must never reveal whether the username or the password was wrong.
export class ErrInvalidCredentials extends Error {
  constructor() {
    super('invalid credentials');
    this.name = 'ErrInvalidCredentials';
  }
}

// ErrAdministratorNotFound is the port-level "not found" error.
export class ErrAdministratorNotFound extends Error {
  constructor() {
    super('administrator not found');
    this.name = 'ErrAdministratorNotFound';
  }
}
