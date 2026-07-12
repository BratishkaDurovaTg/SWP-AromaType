# User Acceptance Tests

This document tracks manual UAT scenarios. UAT evidence supports release readiness, but it does not count as automated QRT evidence.

## UAT-001: Complete Questionnaire and View Recommendations

**Linked user stories:** US-001, US-006, US-007

**Scenario:** A user opens AromaType, completes the guided questionnaire, and receives a recommendation result with a profile and matching fragrances.

**Acceptance signal:** The user can finish the questionnaire without blocking UI defects and can understand why each fragrance was recommended.

**Evidence:** To be recorded in the relevant sprint review notes or issue comments after customer/team validation.

---

## UAT-002: View Product Details

**Linked user stories:** US-003

**Scenario:** A user opens a recommended fragrance card and views notes, accords, price, volume options, and order contact action.

**Acceptance signal:** Product information is understandable and matches the catalog data.

**Evidence:** To be recorded in the relevant sprint review notes or issue comments after customer/team validation.

---

## UAT-003: Catalog Manager Adds a Fragrance

**Linked user stories:** US-002, US-010

**Scenario:** A catalog manager creates a fragrance record with product text, notes, accords, price, volume options, psychotype scores, and image URL through the separate catalog management workflow.

**Acceptance signal:** The created fragrance can appear in the recommendation/catalog flow when active.

**Evidence:** To be recorded in the relevant sprint review notes or issue comments after customer/team validation.

---

## UAT-004: Complete Checkout with Validation

**Linked user stories:** US-004, US-012

**Scenario:** A user adds fragrances to the cart, proceeds to checkout, and submits the order with valid address, phone number, and email.

**Acceptance signal:** Validation errors are shown for missing or incorrect fields. The order is submitted successfully when all fields are valid.

**Evidence:** Recorded in the Sprint Review transcript (2026-07-03, 04:27–05:32). Validation issues identified and logged as bugs.

---

## UAT-005: Remove Items from Cart with Confirmation

**Linked user stories:** US-004, US-012

**Scenario:** A user removes an item from the cart and sees a confirmation dialog before the item is deleted.

**Acceptance signal:** A confirmation dialog appears before removal. The item is deleted only after the user confirms.

**Evidence:** Recorded in the Sprint Review transcript (2026-07-03, 07:36–08:11). Confirmation dialog is missing and will be implemented in the next Sprint.

---

## UAT-006: Gender Field Displays Correctly

**Linked user stories:** US-003

**Scenario:** A user opens a fragrance product card and sees the correct gender category (Male / Female / Unisex) based on the database data.

**Acceptance signal:** The gender field is present and matches the data stored in the database.

**Evidence:** Verified in Sprint 4. Gender now correctly retrieved from the database and displayed on the UI.

---

## UAT-007: Price Updates by Volume

**Linked user stories:** US-004, US-012

**Scenario:** A user selects different volumes (3 ml, 5 ml, 10 ml) in the product card and sees the price update accordingly.

**Acceptance signal:** The displayed price changes correctly when the volume is changed.

**Evidence:** Verified in Sprint 4. Price updates correctly by volume in the product card and cart.

---

## UAT-008: Admin Can Edit Fragrance Information

**Linked user stories:** US-010

**Scenario:** An admin logs into the catalog bot, edits an existing fragrance's details (name, price, notes, gender), and the changes are reflected in the public Mini App.

**Acceptance signal:** Edited fragrance information is updated in the database and displayed correctly in the product card.

**Evidence:** Verified in Sprint 4. Admin edits are saved and visible in the catalog.

