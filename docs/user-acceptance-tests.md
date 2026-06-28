# User Acceptance Tests (UAT)

This document defines the formal acceptance scenarios that the customer executed during the Sprint Review on 26 June 2026.

> **Note:** The customer tested a web-based prototype during the session. The Telegram Mini App was finalized and deployed on 29 June 2026 (Sunday) and is now the primary delivery channel.

---

## UAT-001: Complete questionnaire and receive recommendations

**Status:** Passed

**Steps:**
1. Open the prototype (web version).
2. Answer all questions.
3. Submit the questionnaire.

**Expected result:**
- Recommendations are shown with personality-based reasoning.
- The scoring logic (four personality types, weighted scoring) is explained.

**Actual result:**
Passed. The customer executed the questionnaire successfully. Confirmed it is significantly better than the previous Sprint.

**Feedback received:**
- Several questions feel too direct.
- Some answers reveal fragrance categories too obviously.
- Recommendations should feel personality-driven rather than deterministic.

**Evidence:** [Sprint Review recording 01:50–03:52](https://github.com/BratishkaDurovaTg/SWP-AromaType/blob/main/reports/week4/customer-review-transcript.md)

---

## UAT-002: View recommendation results and understand the logic

**Status:** Passed

**Steps:**
1. Complete the questionnaire.
2. Review the recommended fragrances.
3. The team explains how the algorithm works.

**Expected result:**
- The customer understands the recommendation methodology.
- The ranking logic is accepted.

**Actual result:**
Passed. The team explained the recommendation algorithm based on four personality types and weighted scoring. The customer accepted the methodology.

**Evidence:** [Sprint Review recording 06:05–08:28](https://github.com/BratishkaDurovaTg/SWP-AromaType/blob/main/reports/week4/customer-review-transcript.md)

---

## UAT-003: Identify UI issues and product data gaps

**Status:** Passed with defects

**Steps:**
1. Navigate through the prototype.
2. Check for visual inconsistencies or incorrect data.

**Expected result:**
- The interface displays correct data.
- No visual defects are present.

**Actual result:**
Passed with defects. The customer noticed incorrect percentage values (frontend bug). Recommended removing unnecessary UI elements. Also requested a real product database with images and richer descriptions.

**Follow-up actions:**
- Fix percentage display bug
- Remove redundant UI elements
- Implement admin panel for product catalog (approved by customer)
- Add real product data and images

**Evidence:** [Sprint Review recording 09:07–12:26](https://github.com/BratishkaDurovaTg/SWP-AromaType/blob/main/reports/week4/customer-review-transcript.md)

---

## Summary of Customer Feedback and Resulting PBIs

| Feedback point | Status | Response |
|----------------|--------|----------|
| Questions feel too direct; answer order should be remixed. | To Do | Will be added to Product Backlog. |
| Frontend percentage values are incorrect. | Done | Bug logged and fixed. |
| Unnecessary UI elements should be removed. | Done | UI cleanup done. |
| Real product images and richer descriptions needed. | In progress | Admin panel in progress; product data to be added. |
| Admin panel for product catalog management. | In Progress | Approved and being implemented this Sprint. |

---

## Quality Gates Continuing into Future Work

- Manual UAT with customer at each Sprint Review
- Frontend validation before release
- Recommendation algorithm verification
- Regression testing after new features are added