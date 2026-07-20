# Sprint Retrospective — Week 7 (Sprint 5)

## What went well

* All outstanding customer-reported UI bugs were fixed within the Sprint:
  duplicate images, cart overlap, missing checkout validation, and psychotype
  score input confusion.
* The customer accepted all fixes and confirmed the product is ready for
  handover during the Sprint Review.
* The psychotype score entry in the catalog bot was simplified from a confusing
  multi-step input to a comma-separated format, which the customer explicitly
  appreciated.
* The final MVP v3 release (v1.4.0) was prepared with updated changelog,
  documentation, and tested fixes.
* The Sprint Review transcript shows clear customer acceptance with no
  additional functional changes requested.
* Architecture documentation and development process documentation were added
  to support long-term maintainability.

## What did not go well

* Administrator documentation for psychotype assignment was started but not
  completed by the end of the Sprint.
* Private credential transfer (repository access, bot tokens, passwords) was
  discussed but not executed — it remains a post-Sprint action.
* Git tags were not pushed for the v1.3.0 and v1.4.0 releases, which would
  improve release traceability.
* Database backup automation was deferred and not implemented.
* The public sanitized demo video for MVP v3 was not yet recorded by the end
  of the Sprint.

## What changed compared to the previous Sprint

Compared to Sprint 4 (Week 6), the team made the following changes based on
the previous retrospective action points:

| Previous action point | Status | Result |
|---|---|---|
| Complete domain and VPS transfer | Not done — customer deprioritised deployment | Customer confirmed deployment is not a priority; source code transfer is sufficient |
| Share Telegram bot tokens and admin credentials | Scheduled for post-Sprint private transfer | Process was discussed but not completed |
| Fix remaining UI bugs | **Done** — all identified bugs fixed | Duplicate images, cart overlap, checkout validation, psychotype input — all resolved |
| Finalise documentation | Partially done — architecture and process docs added; admin doc still pending | Customer accepted current state with follow-up for admin documentation |
| Prepare final release | **Done** — v1.4.0 prepared with changelog | Changelog updated; tags still need to be pushed |

The team adapted to the customer's revised handover expectations: instead of
pushing infrastructure transfer, we prepared the source code package and
administrator components as the primary deliverables.

## Action points

1. **Complete the private transfer** — share repository access, source code
   archive, admin bot credentials, and administrator documentation through a
   secure private channel.
2. **Create SemVer tags** — push git tags for v1.3.0 and v1.4.0 to the remote
   repository.
3. **Record the public sanitized demo video** for MVP v3 and link it from the
   Week 7 report and the final release.
4. **Complete the administrator documentation** for psychotype assignment and
   catalog bot usage.
