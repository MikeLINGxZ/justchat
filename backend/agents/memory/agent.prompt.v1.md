> Deprecated reference only.
>
> The effective memory-agent prompt is loaded from `prompt/system.memory.md` at runtime, with the fallback source defined in `backend/pkg/prompts/prompts.go`.
>
> This file is intentionally kept minimal to avoid future drift toward the old schema.

## Current Rules

- Let the Memory Agent decide whether a message deserves long-term storage.
- Only store information about the user.
- Supported memory types are exactly:
  - `fact`
  - `information`
  - `event`
- Do not use or document legacy fields such as:
  - `time_range_start`
  - `time_range_end`
  - `location`
  - `characters`
  - `context`
  - `importance`
  - `emotional_valence`
- Do not introduce legacy memory subtypes such as `plan` or `skill`.
- Never treat image/file/code/webpage content, tool output, or assistant explanations as user memory.

## Source Of Truth

If you need to change memory behavior, update:

1. `backend/pkg/prompts/prompts.go`
2. The generated / persisted `system.memory.md` prompt content

Do not expand this file back into a separate full prompt unless the runtime loading path changes.
