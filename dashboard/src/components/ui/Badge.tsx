import type { PropsWithChildren } from "react";
import { cx } from "../../lib/cx";

type BadgeTone = "critical" | "warning" | "success" | "neutral" | "info";

type BadgeProps = PropsWithChildren<{
  tone?: BadgeTone;
  className?: string;
}>;

export function Badge({ children, tone = "neutral", className }: BadgeProps) {
  return <span className={cx("badge", `badge--${tone}`, className)}>{children}</span>;
}
