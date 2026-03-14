# Dashboard

This package contains the Brainless WAF dashboard UI. Phase 3.1 establishes the design system and reusable component library that later dashboard pages will build on.

## Stack

- React 18
- TypeScript in strict mode
- Vite for local development and builds
- Recharts for time-series visualizations
- Vitest + Testing Library for component tests

## Scripts

```bash
npm install
npm run dev
npm run build
npm run lint
npm run type-check
npm run test
```

## Structure

```text
dashboard/
  src/
    app/              # App shell and showcase screen
    components/
      charts/         # Dashboard visualizations
      ui/             # Reusable design-system primitives
    lib/              # Small utilities
    styles/           # Tokens and global styles
    test/             # Test setup
```

## Design System Foundation

- CSS custom properties define brand, surface, and status tokens.
- UI primitives live under `src/components/ui` and avoid page-specific coupling.
- The initial app is a showcase screen for tokens, cards, metrics, chart styling, buttons, and status badges.
