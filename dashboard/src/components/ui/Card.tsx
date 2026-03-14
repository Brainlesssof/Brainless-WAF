import type { PropsWithChildren, ReactNode } from "react";
import { cx } from "../../lib/cx";

type CardTone = "default" | "highlight" | "muted";

type CardProps = PropsWithChildren<{
  title?: string;
  eyebrow?: string;
  actions?: ReactNode;
  tone?: CardTone;
  className?: string;
}>;

export function Card({ children, title, eyebrow, actions, tone = "default", className }: CardProps) {
  return (
    <section className={cx("card", `card--${tone}`, className)}>
      {(title ?? eyebrow ?? actions) && (
        <header className="card__header">
          <div>
            {eyebrow ? <p className="card__eyebrow">{eyebrow}</p> : null}
            {title ? <h2 className="card__title">{title}</h2> : null}
          </div>
          {actions ? <div className="card__actions">{actions}</div> : null}
        </header>
      )}
      {children}
    </section>
  );
}
