#!/bin/sh
# Enforces Firstmate delivery modes at Iris task completion.
# A required gate that is missing fails task.completed. It never starts
# `no-mistakes axi run`: that command publishes and outlives the hook timeout.
set -u

hook=${SPYNEL_HOOK:-${1:-}}
workspace=${SPYNEL_WORKSPACE:-}
payload=$(cat)
mode=
event=
outcome=

note() {
	printf 'no-mistakes: %s\n' "$1" >&2
}

field() {
	printf '%s' "$payload" | sed -n 's/.*"'"$1"'":"\([^"]*\)".*/\1/p'
}

journal() {
	note "$1: $2"
	if [ -z "$workspace" ] || [ ! -d "$workspace" ]; then
		return 0
	fi
	dir="$workspace/.spynel/extensions-state/no-mistakes"
	line="$(date -u +%Y-%m-%dT%H:%M:%SZ) hook=$hook mode=$mode action=$1 event=$event $2"
	mkdir -p "$dir" 2>/dev/null && printf '%s\n' "$line" >>"$dir/journal" 2>/dev/null || true
}

finish() {
	action=$1
	reason=$2
	journal "$action" "$reason"
	if [ "$action" = fail ] && [ "$hook" = task.completed ]; then
		exit 1
	fi
	exit 0
}

event=$(field event_id)
outcome=$(field outcome)

if [ -z "$workspace" ] || [ ! -d "$workspace" ]; then
	finish fail "workspace is unset or unreachable"
fi

mode_file="$workspace/.spynel/no-mistakes-mode"
mode=no-mistakes
if [ -f "$mode_file" ]; then
	mode=$(sed -n '1{s/^[[:space:]]*//; s/[[:space:]]*$//; p;}' "$mode_file")
	if [ -z "$mode" ]; then
		mode=no-mistakes
	fi
fi

case "$mode" in
local-only) finish skip "local-only does not run the gate" ;;
direct-PR) finish skip "direct-PR does not run the gate" ;;
no-mistakes) ;;
no-mistakes-prod-only) note "no surface classifier; requiring the gate" ;;
*) finish fail "unknown mode $mode" ;;
esac

if [ "$hook" = task.completed ]; then
	case "$outcome" in
	failed|cancelled|waiting) finish skip "outcome $outcome is not a ship" ;;
	esac
fi

if ! command -v no-mistakes >/dev/null 2>&1; then
	finish fail "no-mistakes CLI is not installed"
fi
if ! cd "$workspace" 2>/dev/null; then
	finish fail "workspace is unreachable"
fi
if ! status=$(no-mistakes status 2>/dev/null); then
	finish fail "no-mistakes status failed"
fi
if ! printf '%s' "$status" | grep -Eq 'gate:[[:space:]]*/'; then
	finish fail "workspace has no no-mistakes gate"
fi
finish pass "gate present"
