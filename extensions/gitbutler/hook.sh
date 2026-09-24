#!/bin/sh
# Shells the GitButler CLI on Iris task lifecycle hooks: a claimed task gets its
# own branch, and a terminal task commits the workspace onto it.
#
# Every exit is 0. A nonzero hook exit fails the orchestration operation that
# delivered the event, and version control is advisory here.
set -u

hook=${SPYNEL_HOOK:-${1:-}}
workspace=${SPYNEL_WORKSPACE:-}
extension=${SPYNEL_EXTENSION:-gitbutler}
payload=$(cat)

note() {
	printf 'gitbutler: %s\n' "$1" >&2
}

# The payload is one JSON line of orchestrator-produced identifiers and paths.
# Values containing escaped quotes are out of contract and are not parsed.
field() {
	printf '%s' "$payload" | sed -n 's/.*"'"$1"'":"\([^"]*\)".*/\1/p'
}

if [ -z "$workspace" ]; then
	note "SPYNEL_WORKSPACE is unset; skipping $hook"
	exit 0
fi
if ! command -v but >/dev/null 2>&1; then
	note "but CLI is not installed; skipping $hook"
	exit 0
fi
if ! cd "$workspace" 2>/dev/null; then
	note "workspace $workspace is unreachable; skipping $hook"
	exit 0
fi
if ! but status >/dev/null 2>&1; then
	note "$workspace is not a GitButler workspace; skipping $hook"
	exit 0
fi

task_file=$(field file)
if [ -z "$task_file" ]; then
	note "payload carries no task file; skipping $hook"
	exit 0
fi
stem=$(basename "$task_file")
stem=${stem%.*}
slug=$(printf '%s' "$stem" | tr -cs 'A-Za-z0-9._-' '-' | sed 's/^-*//; s/-*$//')
if [ -z "$slug" ]; then
	note "task file $task_file has no usable branch name; skipping $hook"
	exit 0
fi
branch="iris/$slug"

case "$hook" in
task.claimed)
	if but branch list --json 2>/dev/null | grep -Eq "\"name\": *\"$branch\""; then
		note "branch $branch already exists"
		exit 0
	fi
	if ! but branch new "$branch" >/dev/null 2>&1; then
		note "but branch new $branch failed"
	fi
	;;
task.completed)
	event=$(field event_id)
	receipts="$workspace/.spynel/extensions-state/$extension"
	ledger="$receipts/completed-events"
	# task.completed is delivered at least once, so the commit is guarded by a
	# durable ledger rather than by an in-memory check.
	if [ -n "$event" ] && [ -f "$ledger" ] && grep -Fxq "$event" "$ledger"; then
		note "event $event already committed"
		exit 0
	fi
	outcome=$(field outcome)
	route=$(field route)
	message="iris(${route:-task}): $stem ${outcome:-completed}"
	if ! but commit --branch "$branch" -m "$message" >/dev/null 2>&1; then
		note "but commit on $branch failed or had nothing to commit"
		exit 0
	fi
	if [ -n "$event" ] && ! { mkdir -p "$receipts" && printf '%s\n' "$event" >>"$ledger"; } 2>/dev/null; then
		note "could not record event $event; a retry may commit again"
	fi
	;;
*)
	note "unexpected hook $hook"
	;;
esac
exit 0
