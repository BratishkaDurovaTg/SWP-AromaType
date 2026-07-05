# Sprint Review Summary

**Project:** MVP v2 – Telegram Bot and Mini App  
**Event:** Sprint Review (including customer-executed UAT)  
**Language:** English (translated from original meeting)  
**Duration:** 13:09

> **Moodle-only timecodes**
>
> - **Customer-executed UAT:** 00:03–11:54
> - **Sprint Review discussion:** 11:54–13:09

---

## Sprint Goal Reviewed

The Sprint Goal reviewed during the meeting was to deliver and demonstrate the MVP v2 increment, consisting of the Telegram user bot, the updated Telegram Mini App, and the administrative bot. The customer was invited to execute acceptance testing directly on the delivered system and provide additional feedback for future development.

## Delivered MVP v2 Increment

The demonstrated increment included:

- Telegram user bot with the integrated Telegram Mini App.
- Customer questionnaire and product selection workflow.
- Shopping cart functionality.
- Administrative bot (partially demonstrated due to authentication issues).
- Product management using real database data where available.

Although the complete admin demonstration could not be finished because access credentials were unavailable and one demonstration product had been removed from the database, the implemented functionality was successfully presented and reviewed.

## Addressed Customer Feedback

During the review, previously collected customer's feedback was revisited and acknowledged, including:

- Randomizing questionnaire questions.
- Correcting UI text displayed on one page.
- Removing the current sorting option after customer's approval.
- Continuing improvements based on additional feedback identified during testing.

The customer confirmed that additional comments would be provided after further testing.

## Customer-Executed UAT Results

The customer executed acceptance testing throughout the session. The review identified several issues and improvement opportunities, including:

- Missing validation for required address, phone number, and email fields.
- Payment workflow not yet implemented.
- Address selection and synchronization behaviour requiring verification.
- Limited city availability.
- Shopping cart layout issues when multiple products are added.
- Missing confirmation before removing the final product from the cart.
- Administrative login failure caused by unavailable credentials.
- Missing demonstration product due to database changes.

Overall, the customer was able to test the available functionality and confirmed that further feedback would be provided after additional testing.

## Architecture Evidence Discussed

The review demonstrated the integrated MVP v2 architecture by showing interaction between the Telegram user bot, Telegram Mini App, backend services, database, and administrative functionality. During testing, a potential frontend/backend database synchronization issue was identified for further investigation. No additional architecture changes or ADR updates were proposed during the meeting. 

## Quality Requirements and CI Evidence

The review focused on functional acceptance testing of the deployed increment. Quality observations from customer testing identified validation, usability, and data consistency issues that will remain quality priorities for subsequent development. The team will continue applying the project's existing quality assurance and continuous integration practices while implementing these improvements.

## Remaining Gaps, Risks, and Follow-up Product Backlog Items

The Sprint Review identified the following follow-up work for future Product Backlog refinement:

- Implement complete input validation.
- Complete the payment workflow.
- Fix address synchronization behaviour.
- Resolve shopping cart layout issues.
- Add confirmation before deleting the final cart item.
- Restore demonstration product data.
- Restore administrative access and complete admin feature validation.
- Implement questionnaire randomization.
- Apply remaining UI improvements identified during review.

These items will be prioritised during Product Backlog refinement for future sprints.

## Sprint Review Outcome

The customer reviewed the delivered MVP v2 increment, successfully executed acceptance testing on the available functionality, and confirmed that additional feedback would be submitted after further testing. The Sprint Review concluded with agreement to continue addressing the identified improvements and, if necessary, schedule an additional demonstration after restoring the remaining administrative functionality. 
