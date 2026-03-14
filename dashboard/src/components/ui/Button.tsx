import type { ButtonHTMLAttributes, PropsWithChildren } from "react";
import { cx } from "../../lib/cx";

type ButtonVariant = "primary" | "secondary" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

type ButtonProps = PropsWithChildren<
  ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: ButtonVariant;
    size?: ButtonSize;
  }
>;

export function Button({
  children,
  className,
  variant = "primary",
  size = "md",
  type = "button",
  ...props
}: ButtonProps) {
  return (
    <button
      {...props}
      className={cx("button", `button--${variant}`, `button--${size}`, className)}
      type={type}
    >
      {children}
    </button>
  );
}
