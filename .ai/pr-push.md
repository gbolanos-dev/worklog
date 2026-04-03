# PR Push

Generate consistent, high-signal PR drafts from local git context.

**Trigger:** User asks to prepare a PR, draft PR content, or push with a PR message.

## Inputs

- `git branch --show-current`
- `git diff --name-status` and `git diff --stat`
- `git log --oneline <base>..HEAD` (when commits exist on branch)

## Workflow

1. Inspect branch name and diff before writing anything.
2. Choose exactly one category using the precedence below.
3. Generate title: `Category: short imperative description`
4. Generate PR body with sections: `## What changed`, `## Why`, `## Testing done`.
5. Keep all claims grounded in observed git changes. Do not invent tests.

## Category Precedence

1. **Docs** — docs-only files (`*.md`, docs dirs)
2. **Test** — dominant change is tests (`*_test.*`, test dirs)
3. **Bug** — fix keywords (`fix`, `bug`, `regression`, `panic`) and diff aligns
4. **Refactor** — behavior unchanged, structure improved
5. **Chore** — dependency/config/infra housekeeping
6. **Feature** — default for net-new behavior

## Output Template

```
Title: Category: short description

## What changed
- ...

## Why
- ...

## Testing done
- ...
```

## Rules

- Imperative tense, lowercase after colon, no trailing period.
- Title under ~72 characters.
- Use specific file/component names, not vague wording.
- Same format for commit message title and PR title.
- Do **not** use Conventional Commits format (`feat(...)`, `fix(...)`).
- If no tests were run, state: `- Not run (reason)`.
- Never use em dashes (`—`). Use regular dashes (`-`) or rewrite the sentence.

## Post-Merge

After a PR is merged, check off the corresponding item in `PLAN.md` and note the PR number.
