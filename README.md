# SWP-AromaType

AromaType is a Telegram Mini App for personalized fragrance discovery through
style, situations, and feelings rather than perfume terminology.

## Repository Structure

```text
backend/   Go API, PostgreSQL migrations, recommendation logic
frontend/  Web client and admin UI
docs/      API contract, database schema notes, product docs
```

## Local Development

Start the backend stub:

```bash
cd backend
go run ./cmd/api
```

Open:

- API health check: `/health`
- Swagger UI: `/docs`
- OpenAPI spec: `/openapi.yaml`

Start with Docker Compose:

```bash
docker compose up --build
```

PostgreSQL connection for local backend development:

```text
host: localhost
port: 5432
database: aromatype
user: aromatype
password: aromatype
```

## Team Workflow

- `main` stores stable versions.
- Work in feature branches, for example `feature/backend-auth`.
- Keep API changes documented in `docs/api/openapi.yaml`.
- Frontend can use mock data until backend endpoints are ready.
=======
# Week 2 Assignment Report - Aroma Type

## Project Information

### Project Name

Aroma Type

### Short Description

Aroma Type is a Telegram Mini App focused on fragrance discovery and perfume recommendation through guided questionnaires and personalized matching logic.

### License

[MIT License](https://github.com/BratishkaDurovaTg/SWP-AromaType/blob/dev/LICENSE)

# Repository

[GitHub Repository](https://github.com/BratishkaDurovaTg/SWP-AromaType)

---

# User Stories

[user-stories.md](user-stories.md)

---

# Prototype and Interface Artifacts

## Interactive Prototype

[Figma Prototype](https://www.figma.com/design/XjzUCfIDDInU8ZnSrpE1dC/v1?node-id=0-1&t=GuL1SLLDE3jpRViS-1)

Covered User Stories:

* US-01
* US-02
* US-03
* US-04

---

# API Interface

## Swagger UI

[Swagger UI](https://outlet-lonely-jacksonville-surfing.trycloudflare.com/docs)

Note: the Swagger deployment is temporary and depends on the local development environment remaining active.

## OpenAPI Specification

The OpenAPI specification is available through the Swagger UI deployment.


## Implemented API Endpoints

* `POST /api/auth/register`
* `POST /api/auth/login`

## Postman Collection

No Postman collection was created during Week 2.
Swagger UI was used for API testing and demonstration instead.

---

# MVP v0

## Report

[mvp-v0-report.md](mvp-v0-report.md)

---

## Deployment

[Frontend MVP Deployment](https://t.me/aroma_type_test_bot)

---

## Accessible Implementation

The current MVP v0 frontend deployment and temporary backend Swagger deployment serve as the accessible implementation for Week 2.

## Run Instructions

1. Open the deployment URL.
2. Launch the Telegram Mini App.
3. Verify successful initialization.
4. Run smoke-check scenario.

---

## Public Video Demonstration

[Video Demonstration](https://anonfilesnew.com/s/RuCHERmNy_q)

---

# Pull Requests and Reviews

## PR/MR Template

No dedicated PR template was used during Week 2.

## Reviewed PRs


* [PR #2 - docs: add analysis report](https://github.com/BratishkaDurovaTg/SWP-AromaType/pull/2)
* [PR #3 - docs: add customer meeting summary](https://github.com/BratishkaDurovaTg/SWP-AromaType/pull/3)
* [PR #4 - docs: add customer meeting transcript](https://github.com/BratishkaDurovaTg/SWP-AromaType/pull/4)
* [PR #5 - docs: add llm usage report](https://github.com/BratishkaDurovaTg/SWP-AromaType/pull/5)
* [PR #6 - docs: add mvp v0 report](https://github.com/BratishkaDurovaTg/SWP-AromaType/pull/6)
* [PR #7 - docs: add week 2 user stories](https://github.com/BratishkaDurovaTg/SWP-AromaType/pull/7)


---

# Lychee Link Checking

## Lychee Configuration

[Lychee Configuration](https://github.com/BratishkaDurovaTg/SWP-AromaType/blob/dev/.github/workflows/lychee.yml)

## Latest Successful Protected Branch Run

[Latest Successful Protected Branch Run](https://github.com/BratishkaDurovaTg/SWP-AromaType/actions/runs/27499290355)

---

## Excluded Links

The following temporary links were excluded from automatic validation:

* Temporary Cloudflare Swagger deployment

### Manual Verification

All excluded links were manually verified in the browser before submission.

---

# Screenshots

## Protected Default Branch

![Protected Branch](images/protected-branch.png)

---

## Reviewed Pull Request

![Reviewed PR](images/reviewed-pr.png)

---

## Interactive Prototype

![Prototype](images/prototype.png)

---

## MVP v0 Deployment

![MVP v0](images/mvp-v0.png)

---

# Coverage

## Prototype Coverage

The interactive prototype covers the following stable user story IDs:

* US-01
* US-02
* US-03
* US-04

The prototype demonstrates:

* onboarding flow;
* questionnaire navigation;
* loading and error states;
* fragrance recommendation screens.

---

## MVP v0 Coverage

MVP v0 currently provides the technical foundation for:

* Telegram Mini App integration;
* frontend deployment;
* backend authentication setup;
* environment validation.

Related User Stories:

* US-01
* US-02

Detailed smoke-check documentation is available in:
[mvp-v0-report.md](mvp-v0-report.md)

---

# Customer Review Artifacts

## Customer Transcript

The customer approved transcript publication for Assignment 2 documentation purposes.

[customer-meeting-transcript.md](customer-meeting-transcript.md)

---

## Customer Meeting Summary

[customer-meeting-summary.md](customer-meeting-summary.md)

---

# Weekly Analysis

[analysis.md](analysis.md)

---

# LLM Usage Report

[llm-report.md](llm-report.md)
