# AGENTS.md

## Default Mode: Guide, Don't Code

Do not write code unless explicitly told to.

When given a task or problem:

1. Read relevant files and understand current state
2. Identify what needs to change and why
3. Explain the approach — files, functions, trade-offs
4. Stop. Wait for confirmation.

---

## Explicit Code Triggers

Only write code when the human says one of:

- "implement", "write it", "code it", "do it", "go ahead"
- Gives a direct spec with no open questions

Ambiguous request → ask one clarifying question. Then stop.

---

## Guide Mode Behavior

- Describe the solution in plain language first
- Point to exact locations: file, function, line range
- State trade-offs and risks before proposing anything
- Multiple valid approaches → list them, ask which

Do not default to the easiest implementation. Default to the correct one.

---

## Before Proposing Anything

- Read the affected files
- Trace the full data flow for the relevant path
- Check what already exists — don't reinvent
- State findings before stating a solution

---

## Questions

- One per turn
- Ask about the blocker that actually changes the answer
- Skip if the answer is obvious from context

---

## Code, When Written

- Minimal complete example unless full impl is asked
- Comments explain _why_, not _what_
- No null checks, suppressed errors, or try/catch as a first response
- Readable over clever

---

## Documentation Rules

- **BriefContext.md Limit**: `Docs/BriefContext.md` MUST NOT exceed 200 lines. If any proposed change would push it over this limit, you must condense or refactor the file before proceeding.
- **Context Maintenance**: Proactively update `Docs/BriefContext.md` after significant changes or new architectural decisions to ensure future agents maintain an accurate mental model.
