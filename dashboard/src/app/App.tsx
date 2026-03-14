import { lazy, Suspense } from "react";
import { Activity, ArrowUpRight, Blocks, ShieldAlert, Siren, Waves } from "lucide-react";
import { Badge } from "../components/ui/Badge";
import { Button } from "../components/ui/Button";
import { Card } from "../components/ui/Card";
import { MetricCard } from "../components/ui/MetricCard";
import { SectionHeader } from "../components/ui/SectionHeader";

const TrafficChart = lazy(() => import("../components/charts/TrafficChart"));

const trafficData = [
  { time: "08:00", allowed: 460, blocked: 24, throttled: 8 },
  { time: "09:00", allowed: 510, blocked: 40, throttled: 11 },
  { time: "10:00", allowed: 620, blocked: 58, throttled: 14 },
  { time: "11:00", allowed: 700, blocked: 82, throttled: 20 },
  { time: "12:00", allowed: 760, blocked: 94, throttled: 23 },
  { time: "13:00", allowed: 730, blocked: 66, throttled: 18 },
  { time: "14:00", allowed: 810, blocked: 102, throttled: 27 },
];

const tokens = [
  { label: "Canvas", value: "var(--color-canvas)", tone: "token-swatch--canvas" },
  { label: "Panel", value: "var(--color-panel)", tone: "token-swatch--panel" },
  { label: "Primary", value: "var(--color-accent)", tone: "token-swatch--accent" },
  { label: "Success", value: "var(--color-success)", tone: "token-swatch--success" },
  { label: "Warning", value: "var(--color-warning)", tone: "token-swatch--warning" },
  { label: "Critical", value: "var(--color-critical)", tone: "token-swatch--critical" },
];

const typographyScale = [
  { label: "Display", sample: "Operational visibility without visual noise" },
  { label: "Heading", sample: "Cards, sections, and dense dashboard labels remain legible" },
  { label: "Body", sample: "Readable by default, tuned for incident response sessions" },
];

export function App() {
  return (
    <div className="app-shell">
      <div className="app-shell__backdrop app-shell__backdrop--north" />
      <div className="app-shell__backdrop app-shell__backdrop--south" />
      <main className="app-shell__content">
        <section className="hero">
          <div>
            <Badge tone="info">Phase 3.1 foundation</Badge>
            <h1 className="hero__title">Component library and design system for the Brainless WAF dashboard.</h1>
            <p className="hero__description">
              This initial dashboard build defines the visual language, reusable primitives, chart styling, and density rules for
              later Overview, Events, Rules, and Settings screens.
            </p>
          </div>
          <div className="hero__actions">
            <Button size="lg">View tokens</Button>
            <Button size="lg" variant="secondary">
              Open component gallery
            </Button>
          </div>
        </section>

        <div className="metrics-grid">
          <MetricCard
            title="Requests / sec"
            value="12.8k"
            change="+8.2%"
            helper="Rolling 5 minute window"
            tone="success"
            icon={<Activity size={20} />}
          />
          <MetricCard
            title="Threat pressure"
            value="184"
            change="High"
            helper="Critical + warning events"
            tone="critical"
            icon={<ShieldAlert size={20} />}
          />
          <MetricCard
            title="Rule coverage"
            value="1,274"
            change="92 live"
            helper="Custom and managed rules"
            tone="info"
            icon={<Blocks size={20} />}
          />
          <MetricCard
            title="Rate-limited"
            value="2.1%"
            change="Stable"
            helper="Per-endpoint throttling"
            tone="warning"
            icon={<Siren size={20} />}
          />
        </div>

        <div className="content-grid">
          <Card
            eyebrow="Visualization"
            title="Traffic surface"
            actions={
              <Badge tone="success">
                Live-ready <ArrowUpRight size={14} />
              </Badge>
            }
            className="content-grid__wide"
          >
            <p className="card__body-copy">
              Charts use the same spacing, status tones, and contrast rules as the rest of the UI so future pages inherit a
              consistent monitoring language.
            </p>
            <Suspense fallback={<div className="chart-fallback">Loading traffic surface...</div>}>
              <TrafficChart data={trafficData} />
            </Suspense>
          </Card>

          <Card eyebrow="Palette" title="Operational tokens" tone="highlight">
            <div className="token-grid">
              {tokens.map((token) => (
                <div key={token.label} className="token-swatch">
                  <div className={token.tone} />
                  <div>
                    <p className="token-swatch__label">{token.label}</p>
                    <p className="token-swatch__value">{token.value}</p>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </div>

        <section className="section-block">
          <SectionHeader
            eyebrow="Reusable primitives"
            title="Buttons, badges, and cards"
            description="These components establish shared interaction states and density defaults before feature pages arrive."
            actions={<Button variant="ghost">Inspect patterns</Button>}
          />

          <div className="component-grid">
            <Card title="Buttons" eyebrow="Actions">
              <div className="button-row">
                <Button>Primary action</Button>
                <Button variant="secondary">Secondary</Button>
                <Button variant="ghost">Ghost</Button>
              </div>
            </Card>

            <Card title="Status badges" eyebrow="Semantics">
              <div className="badge-row">
                <Badge tone="critical">Critical</Badge>
                <Badge tone="warning">Warning</Badge>
                <Badge tone="success">Healthy</Badge>
                <Badge tone="info">Syncing</Badge>
                <Badge tone="neutral">Draft</Badge>
              </div>
            </Card>

            <Card title="Typography" eyebrow="Hierarchy" className="component-grid__wide">
              <div className="type-scale">
                {typographyScale.map((item) => (
                  <div key={item.label} className="type-scale__row">
                    <span className="type-scale__label">{item.label}</span>
                    <p className="type-scale__sample">{item.sample}</p>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        </section>

        <section className="section-block">
          <SectionHeader
            eyebrow="Layout rules"
            title="Built for dense monitoring surfaces"
            description="The shell balances command-center density with readable incident workflows, using flexible two-column and metric-grid primitives."
          />

          <div className="principles-grid">
            <Card title="Signal-first density" eyebrow="Principle">
              <p className="card__body-copy">
                Surfaces prioritize status, metrics, and event context before decorative elements. Accent color is reserved for actionability,
                not chrome.
              </p>
            </Card>
            <Card title="Consistent status language" eyebrow="Principle">
              <p className="card__body-copy">
                Success, warning, and critical tones are shared across charts, badges, and cards to reduce re-learning during incidents.
              </p>
            </Card>
            <Card title="Responsive behavior" eyebrow="Principle">
              <p className="card__body-copy">
                The shell collapses from asymmetric dashboard columns to a single stack without losing hierarchy on smaller screens.
              </p>
            </Card>
            <Card title="Animation restraint" eyebrow="Principle">
              <p className="card__body-copy">
                Motion is limited to subtle reveal and hover states so the UI feels alive without compromising clarity during investigations.
              </p>
            </Card>
          </div>
        </section>

        <footer className="footer-note">
          <Waves size={16} />
          <span>Phase 3.2 can now build Overview and Events screens on top of the shared tokens, cards, and chart conventions defined here.</span>
        </footer>
      </main>
    </div>
  );
}
