import type { ReactNode } from "react";
import { Badge } from "./Badge";
import { Card } from "./Card";

type MetricCardProps = {
  title: string;
  value: string;
  change: string;
  helper: string;
  tone: "success" | "warning" | "critical" | "info";
  icon: ReactNode;
};

export function MetricCard({ title, value, change, helper, tone, icon }: MetricCardProps) {
  return (
    <Card className="metric-card" tone="muted">
      <div className="metric-card__top">
        <div>
          <p className="metric-card__label">{title}</p>
          <p className="metric-card__value">{value}</p>
        </div>
        <div className="metric-card__icon">{icon}</div>
      </div>
      <div className="metric-card__bottom">
        <Badge tone={tone}>{change}</Badge>
        <span className="metric-card__helper">{helper}</span>
      </div>
    </Card>
  );
}
