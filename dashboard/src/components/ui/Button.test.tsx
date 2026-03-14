import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button } from "./Button";

describe("Button", () => {
  it("uses button semantics and forwards click handlers", async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();

    render(<Button onClick={handleClick}>Deploy rule</Button>);

    const button = screen.getByRole("button", { name: "Deploy rule" });
    await user.click(button);

    expect(button).toHaveAttribute("type", "button");
    expect(button).toHaveClass("button--primary");
    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});