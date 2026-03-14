# Governance

This document describes how the Brainless WAF project is governed, how decisions are made, and how you can gain more responsibility in the project.

---

## Principles

Brainless WAF is governed by these principles:

- **Transparency** — decisions are made publicly in GitHub Discussions or issues, not in private channels
- **Meritocracy** — influence is earned through sustained, quality contributions
- **Consensus-first** — we seek consensus before voting; voting is a last resort
- **Inclusivity** — we actively lower barriers for new contributors
- **Security-first** — when in doubt, the more secure option wins

---

## Roles

### Contributor

**Anyone** who has opened an issue, submitted a PR, written a rule, helped someone in Discord, or contributed in any way.

- No formal process — just show up and contribute
- Listed in [CONTRIBUTORS.md](CONTRIBUTORS.md) after first merged contribution

### Reviewer

A trusted contributor with knowledge of a specific area who reviews PRs in that area.

**How to become one:** A Core Maintainer invites you after noticing consistent, high-quality reviews or contributions over 2+ months.

**Responsibilities:**
- Review PRs in your area of expertise
- Leave constructive, actionable feedback
- Approve PRs you're confident are ready to merge

**Permissions:** Can approve PRs (cannot merge without a second approval from a Maintainer for significant changes)

### Component Owner

Deep expertise in one major component (core engine, management API, dashboard, rules, deployment).

**How to become one:** 3+ months of consistent contributions to that component + nomination by a Core Maintainer.

**Responsibilities:**
- Primary reviewer for PRs touching your component
- Weigh in on architectural decisions for your component
- Maintain and update component documentation

**Permissions:** Can merge PRs in their component (single approval)

### Core Maintainer

Sets project direction, cuts releases, and makes final decisions.

**How to become one:**
1. 6+ months of consistent, high-quality contributions across multiple areas
2. Demonstrated judgment in reviews and discussions
3. Vote by existing Core Maintainers (simple majority)

**Responsibilities:**
- Review and merge PRs across all components
- Cut releases and maintain the changelog
- Moderate community spaces
- Respond to security vulnerability reports
- Set the project roadmap (with community input)
- Enforce the Code of Conduct

**Permissions:** Full repository access, release publishing, Docker Hub and Helm registry publishing

### Security Team

A subset of Core Maintainers plus invited security researchers who handle vulnerability reports.

- Membership by invitation only
- Membership is not publicly disclosed to protect against targeted attacks
- Operates under the process defined in [SECURITY.md](SECURITY.md)

---

## Decision Making

### Day-to-day decisions

Most decisions — bug fixes, routine features, documentation improvements — are made through the PR process:

- Open a PR
- Get at least 1 approval (2 for security-sensitive changes)
- CI passes
- Merge

No formal discussion needed. If a reviewer has concerns, they raise them in the PR.

### Significant decisions

Changes that affect architecture, public API, release cadence, or project direction require broader discussion:

1. Open a **GitHub Discussion** tagged `RFC` (Request for Comments)
2. Include: motivation, proposed change, alternatives considered, migration impact
3. Community discussion period: **14 days minimum**
4. Core Maintainers reach consensus (or vote if consensus fails)
5. Decision documented and linked from the relevant PR or issue

### Breaking changes

Breaking changes require:
- RFC process above
- Entry in CHANGELOG under a new major or minor version
- Migration guide in `docs/migration/`
- Minimum 30-day deprecation notice before removal

### Voting

Voting is used only when consensus cannot be reached after good-faith discussion.

- Each Core Maintainer gets one vote
- Simple majority wins (50%+1)
- Votes are cast publicly in the GitHub Discussion
- Voting period: 7 days
- Minimum participation: 3 Core Maintainers must vote

---

## Roadmap

The project roadmap ([ROADMAP.md](ROADMAP.md)) is set by Core Maintainers with community input:

- Major roadmap decisions are discussed via RFC
- Quarterly review of roadmap priorities in a public GitHub Discussion
- Any community member may propose additions via a Discussion tagged `roadmap`
- Community upvotes (👍 reactions) on issues influence prioritization

---

## Conflict Resolution

If contributors disagree:

1. **Direct resolution** — discuss in the PR or issue thread, seek to understand the other perspective
2. **Maintainer mediation** — ask a Core Maintainer to weigh in
3. **Core Maintainer vote** — if mediation fails, Core Maintainers vote

Personal disputes that cannot be resolved through the above process are escalated to the Code of Conduct process.

---

## Stepping Down

Maintainers and Reviewers may step down at any time by notifying the Core Maintainer group. There is no minimum commitment. Life happens — we understand.

Inactive Maintainers (no activity for 6+ months) may be moved to Emeritus status after a private check-in. Emeritus Maintainers retain their GitHub access but are not counted for voting quorum.

---

## Changes to This Document

Changes to this Governance document require:
- RFC process (14-day discussion period)
- Approval by 2/3 of Core Maintainers

---

## Current Core Maintainers

| Name | GitHub | Area |
|------|--------|------|
| *(Founding team — to be listed at first commit)* | — | — |

---

*This governance model is inspired by the [CNCF governance template](https://github.com/cncf/project-template) and the governance models of projects like Prometheus, Envoy, and ModSecurity.*
