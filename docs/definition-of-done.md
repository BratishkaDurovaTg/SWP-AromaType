# Definition of Done

A Product Backlog Item (PBI) may be marked as **Done** only when all of the following conditions are satisfied:

1. All issue acceptance criteria have been completed and verified.

2. The completed work has been reviewed and approved by at least one other team member through a pull request review.

3. All required Continuous Integration (CI) checks pass successfully, including build, linting, and automated workflows configured for the project.

4. Relevant automated tests have been implemented and pass successfully. This includes unit tests, integration tests, or other applicable automated tests depending on the type of work.

5. Relevant automated Quality Requirement Tests (QRT), such as dependency or security checks, have been executed successfully where applicable.

6. Code coverage for critical modules meets the team's agreed coverage expectations and does not decrease below the accepted threshold.

7. Testing and verification evidence is preserved through pull requests, CI workflow results, test reports, or linked project documentation.

8. Verification evidence is linked to the corresponding Product Backlog Item whenever appropriate.

9. Related documentation has been updated when required.

10. `CHANGELOG.md` has been updated for every user-visible change.

11. For user stories, all linked supporting PBIs required for implementation, testing, review, and verification have been completed.

12. For implementation and technical PBIs, the linked pull request has been merged into the protected default branch.

13. The Product Backlog Item Work Status has been updated to **Done**.
