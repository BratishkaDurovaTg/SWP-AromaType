# Week 6 Reflection

## Learning points

During Week 6 (Sprint 4), the team focused on delivering a stable trial release, preparing the product for customer transition, and completing customer-facing documentation.

We learned that **customer handover is not just about code** — it requires clear documentation, access management, and transparent communication about what is being transferred and what remains under team control. Creating `docs/customer-handover.md` helped us structure the handover conversation and identify gaps early.

We also learned that **trial releases are valuable**. By letting the customer try the Week 6 trial version, we received feedback on usability issues that we would not have discovered otherwise. Issues like cart overlap, missing validation errors, and confusing UI elements were identified and prioritised for the next Sprint.

Another key lesson was that **infrastructure planning must start earlier**. The domain setup (`aroma-type.ru`), VPS access management, and Telegram bot token transfer took longer than expected. In future projects, we would prepare these items before the final handover week.

## Validated assumptions

Several assumptions were confirmed during Sprint 4:

- **The product is usable without payment integration.** The customer accepted the simplified checkout flow (Telegram contact for orders) and confirmed that it is sufficient for the current MVP scope.
- **The customer is willing to accept the product with follow-up items.** The customer confirmed that the current handover level ("Ready for independent use") is acceptable, provided that the remaining actions (domain transfer, bot token transfer, VPS access) are completed in the following week.
- **The redesigned UI meets customer expectations.** The customer approved the pastel colour palette, premium styling, and overall interface direction.

## Friction and gaps

Several challenges were identified during the Sprint:

- **Domain and infrastructure transfer** — The domain `aroma-type.ru` was purchased but not yet transferred to the customer. The VPS is still under team control, and the Telegram bot tokens have not been shared.
- **Remaining UI bugs** — Several issues identified during UAT (cart overlap, price by volume, gender filter, psychotype bar colour) could not all be fixed before the trial release.
- **Documentation gaps** — While `docs/customer-handover.md` was created, some sections (database backup, monitoring setup) are still incomplete and will be addressed in the next Sprint.
- **CI coverage** — Some critical modules still have lower coverage than the 30% target, requiring further test expansion.

## Planned response

Based on the Week 6 learnings, the following actions are planned for Sprint 5 (Week 7):

- Complete the transfer of domain, VPS access, and Telegram bot tokens to the customer.
- Fix the remaining critical UI bugs identified during UAT.
- Update `docs/customer-handover.md` with the final handover status and confirmed access details.
- Improve documentation for database backup and recovery procedures.
- Expand automated tests to improve coverage for critical modules.
- Prepare the final `MVP v3` release and the public sanitized demo video.