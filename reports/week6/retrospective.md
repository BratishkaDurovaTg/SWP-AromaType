# Sprint Retrospective — Week 6 (Sprint 4)

## What went well

* Successfully delivered the Week 6 trial release and made it accessible to the customer via the Telegram Mini App.
* Created `docs/customer-handover.md` — a complete handover document covering product status, access, deployment, environment variables, and remaining actions.
* Customer feedback was collected during the trial and UAT sessions, helping us identify critical issues before the final release.
* The customer approved the redesigned UI, pastel colour palette, and overall product direction.
* Critical bugs (gender filter, psychotype points validation, price by volume) were identified and fixed during the Sprint.
* Infrastructure documentation was improved, and the domain `aroma-type.ru` was purchased and configured.

## What did not go well

* Several UI bugs identified during UAT (cart overlap, price display issues, gender filter) could not all be fixed before the trial release.
* The domain and VPS access were not yet transferred to the customer by the end of Week 6.
* Telegram bot tokens and admin passwords were not shared with the customer, blocking full independent use.
* Some documentation sections (database backup, monitoring setup) remained incomplete.
* Infrastructure planning (domain, VPS, bot tokens) started later than it should have, causing delays in the handover process.

## What changed compared to the previous Sprint

Compared to Sprint 3, the team shifted focus from feature development to transition readiness and documentation. We spent more time preparing for customer handover, creating handover documentation, and fixing critical bugs identified during UAT. We also added `CONTRIBUTING.md` and `AGENTS.md` to improve repository guidance for future contributors and AI agents.

## Action points for the next Sprint

1. **Complete the domain and VPS transfer** — finalise the handover of `aroma-type.ru` and server access to the customer.
2. **Share Telegram bot tokens and admin credentials** — transfer all necessary secrets to the customer.
3. **Fix remaining UI bugs** — address cart overlap, price by volume, gender filter, and psychotype bar colour issues.
4. **Finalise documentation** — complete missing sections in `docs/customer-handover.md` and add database backup and recovery procedures.
5. **Prepare final release** — create the `MVP v3` release, update `CHANGELOG.md`, and record the public sanitized demo video.