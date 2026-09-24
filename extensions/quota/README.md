# Quota heartbeat companion

Optional. Copy `wake.sh` to `.spynel/extensions/quota/wake.sh`. Iris runs it
before a semantic heartbeat provider turn.

| First line | Result |
| --- | --- |
| `absorb` | No provider turn. |
| anything else | Provider turn, with the remaining output appended as a quota snapshot. |

`absorb` means `.spynel/tasks/waiting` has no Markdown file. A waiting task
prints `dispatch` and, when `quota-axi` is installed, a
`quota-axi --no-credential-refresh` snapshot. The script does not rank models
or send a reminder. A missing script leaves the heartbeat unchanged.
