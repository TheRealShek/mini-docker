# AGENTS.md

## Default Workflow

Default to guidance before code.

For any task:

1. Read the relevant files and understand the current behavior.
2. State the important findings first.
3. Explain the recommended change, including the affected files and meaningful trade-offs.
4. Wait for confirmation before editing unless the user has clearly asked for implementation or provided a complete direct spec.

## When To Code

Write code only when one of these is true:

- The user explicitly asks for implementation.
- The user gives a direct, complete spec with no unresolved decisions.

If the request is ambiguous and the answer would change the implementation, ask one focused clarifying question before proceeding.

## Guidance Standards

- Prefer the correct solution over the easiest one.
- Point to exact files, functions, or sections when discussing changes.
- Trace the relevant behavior far enough to avoid guessing; go deeper when the change affects shared logic or user-visible behavior.
- Check what already exists before proposing a new abstraction or pattern.
- When there are multiple valid approaches with meaningful trade-offs, explain them and ask which direction to take.

## Code Standards

- Keep changes minimal, complete, and readable.
- Prefer real implementation over placeholder examples when implementation is requested.
- Add comments only when they explain why a non-obvious choice exists.
- Use validation and error handling where the design requires them; avoid defensive code that hides problems without a clear reason.

## Questions

- Ask only questions that materially affect the answer.
- Keep questions focused and skip them when the answer is already clear from context.

## Documentation Rules

- `README.md` and `Docs/BriefContext.md` must describe the codebase as it exists now.
- Other docs may describe intended architecture or roadmap and do not need to match partial implementation state unless the user explicitly asks.
- Keep `Docs/BriefContext.md` at or below 200 lines. If an update would exceed that limit, condense it before proceeding.
- Update `Docs/BriefContext.md` after significant implementation changes or architectural decisions so future agents have an accurate short-form context.
