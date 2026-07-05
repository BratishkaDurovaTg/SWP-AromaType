# Week 5 Reflection

## Learning points

During this Sprint, the team learned that documenting architecture and development processes helps clarify the product structure and communication between team members. Creating Architecture Decision Records (ADRs) and maintaining `docs/development-process.md` provided a clear reference for technical choices.

We also learned that delivering MVP v2 required balancing new features, bug fixes, and technical documentation. The domain setup and HTTPS configuration (`aroma-type.ru`) highlighted the importance of proper infrastructure planning before the deployment phase.

Customer feedback collected through UAT sessions proved valuable in identifying usability issues, especially in the cart and ordering flow, which were not obvious during development.

## Validated assumptions

The Sprint Review confirmed that the redesigned UI (pastel colors, premium styling, consistent fonts) met customer expectations. The cart functionality (US-012) and sample set ordering flow (US-004) were successfully delivered and approved by the customer.

The customer also validated the need for administrator tools and a real product catalog. The decision to use Telegram Mini App as the primary platform was confirmed, although domain and HTTPS setup required additional effort.

## Friction and gaps

Several issues were identified during the Sprint and still require attention:

- **Cart overlap with checkout button**: The frontend bug was logged and is in progress.
- **Price does not update by volume in product card**: Admin flow needs redesign.
- **Gender filter shows only "Female"**: Bug needs to be fixed.
- **Psychotype bar color changes to pink**: Design fix is pending.

Additionally, configuring the domain (`aroma-type.ru`) and obtaining a public IP for the University VM took longer than expected. This highlighted the need for earlier infrastructure planning.

## Planned response

In the next Sprint, we will focus on:
- Fixing the remaining UI bugs (cart overlap, gender filter, price by volume).
- Finalizing the product catalog and administrator tools.
- Improving automated testing and CI checks.
- Preparing the final presentation and submission.

We will also ensure that all feedback collected during UAT and Sprint Review is addressed in the upcoming Sprint.