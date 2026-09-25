#!/bin/sh
# Companion poll. An unchanged workflow tree prints nothing and exits zero.
# A change prints one wake line and a quota-axi snapshot. This script does not
# start Iris, rank models, or decide reminders.
set -u

workspace=${SPYNEL_WORKSPACE:-.}
state="$workspace/.spynel/extensions-state/quota"
mkdir -p "$state" 2>/dev/null || exit 0

fingerprint() {
	find "$workspace/.spynel/tasks" "$workspace/.spynel/goals" \
		-type f -name '*.md' -exec cksum {} + 2>/dev/null | sort
}

current=$(fingerprint)
previous=$(cat "$state/fingerprint" 2>/dev/null || true)
if [ "$current" = "$previous" ]; then
	exit 0
fi
printf '%s\n' "$current" >"$state/fingerprint"
printf 'wake: workflow files changed\n'
if command -v quota-axi >/dev/null 2>&1; then
	quota-axi --no-credential-refresh 2>/dev/null || printf 'quota-axi failed\n'
else
	printf 'quota-axi is not installed\n'
fi
exit 0
