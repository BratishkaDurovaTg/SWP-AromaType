# AromaType - Assignment 4 Report (Week 4)

## Testing Scope

Assignment 4 introduced maintained automated testing and QA gates for the AromaType product repository.

Primary testing documentation:

- [Testing status](../../docs/testing.md)
- [Quality requirements](../../docs/quality-requirements.md)
- [Quality requirement tests](../../docs/quality-requirement-tests.md)
- [User acceptance tests](../../docs/user-acceptance-tests.md)

## Critical Modules

| Critical module | Reason | Required coverage | Current coverage |
|---|---|---:|---:|
| `backend/internal/auth` | Registration, login, JWT, roles. | 30% | 52.6% |
| `backend/internal/questionnaire` | Recommendation and fragrance business logic. | 30% | 45.7% |
| `backend/internal/http` | API routing, admin access control, request/response behavior. | 30% | 37.3% |

Coverage was measured with:

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Unit and Integration Tests

Unit tests were added for:

- JWT registration, login, admin bootstrap, token verification.
- Rule-based recommendation ranking and profile generation.
- Fragrance creation validation and input cleanup.
- Photo content type allow-listing.

Integration tests were added for:

- Health endpoint response.
- Public registration, questionnaire loading, and recommendation flow through the HTTP router.
- Admin fragrance creation protection against regular user tokens.

Tests are stored in normal Go test locations under `backend/internal/**`.

## Additional QA Check Options Considered

The team considered these additional QA checks:

- Dependency vulnerability scanning.
- API contract validation against OpenAPI examples.
- Basic performance smoke testing for key endpoints.
- Accessibility checks for the frontend.
- Dependency freshness or health checking.

## Additional QA Check Selected

The selected additional QA check is Go dependency vulnerability scanning with `govulncheck`.

During local validation, this check found a reachable vulnerability in `github.com/golang-jwt/jwt/v5@v5.2.1`. The dependency was upgraded to `v5.2.2`, and the scan now reports that the application code is affected by 0 vulnerabilities.

## QA Objective or Risk Addressed

The check addresses the risk that a vulnerable dependency or reachable Go standard-library vulnerability could affect authentication, admin catalog mutation, or backend deployment integrity.

## Why This Risk Matters

AromaType stores user accounts, issues JWTs, and allows admins to create catalog items shown to users. Known reachable vulnerabilities in backend dependencies can create avoidable security and reliability risk even when the product logic itself is tested.

## Where the Check Runs in CI

The check runs in the GitHub Actions workflow:

- [Additional QA workflow](../../.github/workflows/qa.yml)

The command is:

```bash
cd backend
$(go env GOPATH)/bin/govulncheck ./...
```

## CI Evidence

CI workflows used for Assignment 4:

- [Backend workflow](../../.github/workflows/backend.yml)
- [Docker Build workflow](../../.github/workflows/docker-build.yml)
- [Lychee workflow](../../.github/workflows/lychee.yml)
- [Additional QA workflow](../../.github/workflows/qa.yml)

Branch protection evidence is maintained in the repository settings and protected-branch PR history.

## Important Limitations and Deferred QA Work

- Repository/database methods are not yet covered with a live PostgreSQL integration test because current tests focus on service and API behavior with fake repositories.
- Full browser end-to-end tests are deferred until the frontend routes stabilize.
- Accessibility and performance checks were considered but deferred behind backend testing and dependency vulnerability scanning.
- Manual UAT evidence supports release readiness but does not count as QRT evidence.
