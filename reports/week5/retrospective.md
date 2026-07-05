# Sprint Retrospective

## What went well

* We successfully delivered MVP v2 with cart functionality (US-012) and sample set ordering (US-004), which were approved by the customer during the Sprint Review.
* The UI redesign based on customer feedback (pastel colors, premium styling, consistent fonts) was completed and positively received.
* Architecture documentation and ADRs were created, providing a clear reference for technical decisions.
* The domain `aroma-type.ru` was purchased and configured, and HTTPS was set up through Caddy.
* Customer feedback was collected through structured UAT sessions and used to improve the product.

## What did not go well

* Several UI bugs were identified during UAT (cart overlap, price by volume, gender filter, psychotype bar color) and could not all be fixed before the Sprint end.
* The University VM did not have a public IP, which caused delays in deployment and required alternative solutions.
* Domain and HTTPS setup took longer than expected due to DNS propagation and infrastructure limitations.
* Some planned functionality had to be postponed because higher-priority improvements were identified during the customer review.

## What changed compared to the previous Sprint

Compared to the previous Sprint, we focused more on architecture documentation and ADRs, which improved the team's understanding of technical decisions. We also delivered MVP v2 with customer-approved UI and cart functionality. The domain and HTTPS setup was completed, although it required additional effort due to infrastructure constraints. Customer feedback was collected through UAT sessions and used to prioritize bug fixes for the next Sprint.

## Process improvements for the next Sprint

* Validate new features with the customer earlier so feedback can be incorporated before the end of the Sprint.
* Continue refining the Product Backlog throughout the Sprint to keep priorities clear and make Sprint Planning more efficient.
* Ensure infrastructure planning (domain, VM, public IP) is completed earlier to avoid last-minute delays.
* Allocate more time for UI bug fixes and testing before the Sprint Review.