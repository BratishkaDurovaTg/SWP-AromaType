# Week 7 Reflection

## Learning points

During Week 7 (Sprint 5), the team focused on follow-up maintenance, bug fixes,
customer feedback response, and final delivery of MVP v3.

We learned that **customer handover expectations can shift**. In Week 6, we
assumed the customer wanted full infrastructure transfer (VPS, domain, bot
tokens). In the Week 7 Sprint Review, the customer explicitly stated that
deployment was not a priority — they preferred receiving the complete source
code, administrator components, and documentation through a private channel.
This reinforced the importance of asking the customer directly rather than
making assumptions.

We also confirmed that **focused bug-fix sprints are effective**. Sprint 5 was
the smallest sprint in terms of new features, but the team resolved all
outstanding customer-reported issues: duplicate product images, cart scrolling,
checkout validation errors, and the confusing psychotype score input in the
admin bot. The customer accepted all fixes and confirmed the product is ready
for handover.

Another key lesson was that **administrator experience matters for long-term
product usefulness**. Simplifying the psychotype score entry in the catalog bot
(from a confusing multi-step input to comma-separated values) was a small
change that significantly improved the admin workflow. The customer explicitly
appreciated this improvement and requested additional administrator
documentation.

## Validated assumptions

Several assumptions were confirmed during Sprint 5:

- **The customer values source code and documentation over deployment.**
  The customer confirmed that the most important deliverable is the complete
  frontend, backend, and administrator source code together with supporting
  documentation. Deployment and infrastructure configuration were explicitly
  deprioritised.

- **The remaining UI bugs were the critical blockers.** Once the cart overlap,
  duplicate image, missing error messages, and psychotype input issues were
  fixed, the customer accepted the product without requesting additional
  functional changes.

- **A pre-recorded demo is sufficient for handover.** The customer did not
  request a live demo during the final review — the trial release from Week 6
  and the fix verification in Week 7 were enough to confirm readiness.

## Friction and gaps

Several challenges were encountered during the Sprint:

- **Administrator documentation was not completed.** While the psychotype
  input was simplified, the written administrator guide for psychotype
  assignment was still pending by the end of the Sprint. This was noted as a
  follow-up item.

- **Private credential transfer remains unresolved.** Repository access and
  admin credentials are still scheduled for private-channel transfer after
  submission. The process was discussed but not executed during the Sprint.

- **Database backup automation was not implemented.** The planned backup
  mechanism was deferred due to the customer's deprioritisation of
  infrastructure.

- **Git tags were not created for v1.3.0 and v1.4.0 releases.** The changelog
  was updated but the corresponding SemVer tags were not pushed.

## Planned response

Based on the Week 7 learnings:

- Complete the private transfer of repository access, source code, admin
  credentials, and administrator documentation after the Sprint Review.
- Create SemVer tags for the v1.3.0 and v1.4.0 releases.
- Record and link the public sanitized demo video for MVP v3.
- Prepare the Week 7 rehearsal presentation and the Demo Day slide deck.
