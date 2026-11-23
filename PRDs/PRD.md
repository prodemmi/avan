# Product Requirements Document: Messaging Service Enhancement (Phased)

## Phase 1: UI/UX Standardization & Admin Experience
**Objective:** Improve the administrative interface consistency, usability, and layout efficiency.

### 1.1 Design System & Component Consistency
* **Typography:** Implement `IRANSans` font for all administrative panels and UI components.
* **Button Standardization:** Apply a strict Material Design system. All action buttons must have standardized dimensions and consistent padding/scaling.
* **Localization & Date Handling:** 
    * Implement dynamic date formatting: Use `Jalali` (Persian) calendar when the locale is `fa` and `Gregorian` when the locale is `en`.
* **Swagger Documentation:** Standardize the Swagger UI. Disable theme-switching and animations to maintain the default Swagger color palette for stability.

### 1.2 Admin Interface Overhaul
* **User Management:** 
    * Deprecate the "inline edit" pattern in user tables.
    * Implement a **Side Drawer (Slide-out)** pattern from the left/right for editing user details to maintain context.
* **Workspace Context:** 
    * Add a **Workspace Selector** component in the Global Admin Header to allow switching between different business contexts/workspaces.
* **Role Management:** 
    * Allow the removal of custom-defined roles. 
    * *Constraint:* System-default/Primary roles must remain immutable.

---

## Phase 2: Provider Configuration & Management
**Objective:** Re-engineer the Provider lifecycle and introduce advanced routing capabilities.

### 2.1 Provider Entity Refactoring
* **Domain Re-modeling:** Refactor the hierarchy from `Workspace -> Provider` to a more robust `Business -> Workspace -> App -> Provider` model to support multi-tenant complexity.
* **Dynamic Configuration UI:** 
    * Replace static forms with dynamic input dependency logic. 
    * Input fields must reactively appear/disappear based on the selected Provider type or specific settings.

### 2.2 Provider Onboarding Wizard
Implement a multi-step Wizard for provider creation with two modes: **Blank Configuration** or **Template-based Configuration**.
* **Template Library:** 
    * **Local Providers:** Pre-configured templates for Iranian services (e.g., Kavenegar, Melipayamak, etc.).
    * **Global Providers:** Pre-configured templates for international services (e.g., Twilio, SendGrid, Slack, Webhooks).
* **Requirement:** Templates must include default required fields (API Keys, Sender IDs, Endpoint URLs) extracted from official developer documentation.

### 2.3 Advanced Provider Routing (Load Balancer)
Implement a **Provider Load Balancer** system to manage multi-provider redundancy and cost optimization.
* **Strategy Support:**
    * **Priority/Weight-based:** Route traffic based on user-defined cost-efficiency (e.g., use the cheapest provider first).
    * **Failover:** Automatically switch to the next provider in the priority list if the primary provider returns a failure/timeout.
    * *Note:* The architecture must be extensible to support additional strategies (Round Robin, Latency-based) in the future.

---

## Phase 3: Core Reliability & API Expansion
**Objective:** Enhance system robustness and provide deeper visibility for business clients.

### 3.1 Reliability & Security
* **Idempotency Implementation:** 
    * Implement an idempotency key mechanism for every "Send Message" request.
    * Ensure the system detects and ignores duplicate requests with the same key to prevent accidental double-billing/sending.
* **Granular Security:** 
    * Introduce a security configuration layer that allows defining unique security policies (IP whitelisting, API Key restrictions) at the **Provider** or **Workspace** level.

### 3.2 API & Observability
* **Message Inquiry API:** 
    * Develop a high-performance Inquiry API specifically for business clients to track message status (Sent, Delivered, Failed, Read) in real-time.
