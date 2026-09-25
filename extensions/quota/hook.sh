#!/bin/sh
# Records quota evidence when a task is claimed. Always exits zero.
# Does not cancel the claim, rank a model, or touch the semantic heartbeat.
set -u

workspace=${SPYNEL_WORKSPACE:-}
payload=$(cat)

field() {
	printf '%s' "$payload" | sed -n 's/.*"'"$1"'":"\([^"]*\)".*/\1/p'
}

note() {
	printf 'quota: %s\n' "$1" >&2
}

if [ -z "$workspace" ] || [ ! -d "$workspace" ]; then
	note "workspace is unset; skipping"
	exit 0
fi

state="$workspace/.spynel/extensions-state/quota"
mkdir -p "$state" 2>/dev/null || {
	note "cannot write state"
	exit 0
}
snapshot="$state/latest.txt"
if command -v quota-axi >/dev/null 2>&1; then
	quota-axi --no-credential-refresh >"$snapshot" 2>/dev/null || printf 'quota-axi failed\n' >"$snapshot"
else
	printf 'quota-axi is not installed\n' >"$snapshot"
fi

file=$(field file)
case "$file" in
"$workspace"/*) ;;
*)
	note "claim has no workspace file"
	exit 0
	;;
esac
if [ ! -f "$file" ]; then
	note "claim file is missing"
	exit 0
fi
if grep -q 'Quota evidence:' "$file" 2>/dev/null; then
	exit 0
fi
if grep -q '## Progress' "$file" 2>/dev/null; then
	printf -- '- Quota evidence: %s (not a route)\n' "$snapshot" >>"$file"
else
	printf '\n## Progress\n\n- Quota evidence: %s (not a route)\n' "$snapshot" >>"$file"
fi
exit 0
