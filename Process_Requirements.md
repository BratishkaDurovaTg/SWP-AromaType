# Process Requirements

## Team Roles

- **Team Lead:** Sergey Berezhnoy — coordination, customer communication, final review
- **Frontend Developer:** Dilya Akhmerova — Telegram Mini App implementation, UI integration
- **Backend Developer:** Nikita Matveev — API, database, deployment
- **UX/UI Designer:** Liza Sotnikova — design, prototypes, UI-kit
- **QA/Documentation:** Viktoria Zorkaltceva — testing, bug tracking, reports

## Workflow

### Issues
- Every user story has a stable ID (US-001, US-002, ...)
- Every PBI has a clear title, description, type, MoSCoW priority, and Story Points
- Acceptance criteria are written for MVP v1 tasks and Sprint-selected tasks
- Work Status: `To Do` | `Ready` | `In Progress` | `Review` | `Done`

### Branches and PRs
- Branch naming: `<issue-number>-short-description` (e.g., `12-api-recommend`)
- Each PR must be linked to an issue (`Related to #X` or `Closes #X`)
- PRs require at least one approval from another team member
- Merge commits only (no squash or rebase)

### Definition of Done


### Estimation
- Story Points are used with the Modified Fibonacci scale: 1, 2, 3, 5, 8, 13, 20, 40, 100
- Estimation is done collaboratively (Planning Poker)

### Sprint Cadence
- Sprint runs from Monday to Sunday
- Sprint Goal is defined in the Sprint milestone description
- Sprint Review with the customer is conducted at the end of each Sprint
- Sprint Retrospective is conducted after the Sprint Review

### MVP Versioning
- MVP versions are tracked using the `MVP version` field in GitHub Projects or labels: `mvp-v1`, `mvp-v2`, `mvp-v3`

### Customer Communication
- All customer meetings are recorded with permission
- Meeting summaries and transcripts (sanitized) are stored in `reports/weekX/`
- Written consent for MIT-licensed public development was obtained before repository creation
