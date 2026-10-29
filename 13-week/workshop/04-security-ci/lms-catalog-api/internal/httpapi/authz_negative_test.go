package httpapi

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Negative authorization tests for the [security] story of STRIDE row E
// (boundary lms-circulation-api -> lms-catalog-api, PR #36): a valid token
// used on a route it was not issued for must answer 403.

const (
	internalSecret = "w13-test-internal-secret-0123456789abcdef"
	bookID         = "c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82"
)

func setup(t *testing.T) (http.Handler, *rsa.PrivateKey, *int) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	ok := func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(http.StatusOK) }
	h := Handlers{Create: ok, Update: ok, Delete: ok, LoanCopy: ok, ReturnCopy: ok}
	return NewRouter(&key.PublicKey, internalSecret, h), key, &calls
}

// circulationToken is what catalogclient in lms-circulation-api mints.
func circulationToken(t *testing.T) string {
	t.Helper()
	now := time.Now().UTC()
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "circulation-service", "scope": "catalog:copies:write",
		"iat": now.Unix(), "exp": now.Add(time.Minute).Unix(),
	}).SignedString([]byte(internalSecret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// administratorToken is what lms-access-api issues after POST /api/v1/auth/login.
func administratorToken(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	now := time.Now().UTC()
	s, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "admin-1", "scope": "catalog:books:write",
		"iat": now.Unix(), "exp": now.Add(time.Hour).Unix(),
	}).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func call(router http.Handler, method, path, token string) int {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code
}

func TestCirculationTokenCannotDeleteBook(t *testing.T) {
	router, _, calls := setup(t)
	if got := call(router, http.MethodDelete, "/api/v1/books/"+bookID, circulationToken(t)); got != http.StatusForbidden {
		t.Fatalf("DELETE /api/v1/books/{id} with the internal token = %d, want 403", got)
	}
	if *calls != 0 {
		t.Fatal("the delete handler ran: catalog.book row would be gone")
	}
}

func TestAdministratorTokenCannotTakeCopyWithoutLoan(t *testing.T) {
	router, key, calls := setup(t)
	if got := call(router, http.MethodPost, "/api/v1/books/"+bookID+"/loan-copy", administratorToken(t, key)); got != http.StatusForbidden {
		t.Fatalf("POST /api/v1/books/{id}/loan-copy with the Administrator token = %d, want 403", got)
	}
	if *calls != 0 {
		t.Fatal("the loan-copy handler ran: available_copies would drop without a Loan")
	}
}

func TestMissingTokenIsStill401(t *testing.T) {
	router, _, _ := setup(t)
	if got := call(router, http.MethodPost, "/api/v1/books/"+bookID+"/loan-copy", ""); got != http.StatusUnauthorized {
		t.Fatalf("loan-copy without token = %d, want 401", got)
	}
}

func TestRightTokenOnRightRoute(t *testing.T) {
	router, key, _ := setup(t)
	if got := call(router, http.MethodPost, "/api/v1/books/"+bookID+"/loan-copy", circulationToken(t)); got != http.StatusOK {
		t.Fatalf("loan-copy with the internal token = %d, want 200", got)
	}
	if got := call(router, http.MethodPatch, "/api/v1/books/"+bookID, administratorToken(t, key)); got != http.StatusOK {
		t.Fatalf("PATCH /books/{id} with the Administrator token = %d, want 200", got)
	}
}
