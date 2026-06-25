# Definition of Done

A Product Backlog Item (PBI) may be marked as Done only when all of the following conditions are satisfied:

1. All issue acceptance criteria are satisfied.

2. The completed work has been reviewed by at least one other team member.

3. Required tests, checks, and verification activities have been completed successfully.

4. Verification evidence is preserved through normal workflow artifacts such as issues, pull requests, reviews, workflows, or comments.

5. Relevant CI checks pass, including formatting, static analysis, unit tests, integration tests, coverage reporting, Docker build, link checking, and the additional QA check documented in [testing.md](testing.md).

6. Critical product modules affected by the change keep at least 30% automated line coverage unless a documented replacement or TA-approved exception exists.

7. Related documentation has been updated when necessary.

8. CHANGELOG.md has been updated for every user-visible change.

9. For user stories, linked supporting PBIs provide the required implementation, review, and verification evidence.

10. For implementation and technical PBIs, the linked pull request has been merged into the protected default branch.

11. The Product Backlog Item Work Status has been updated to Done.
