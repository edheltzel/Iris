#!/bin/sh
# Decides whether an Iris heartbeat tick needs a provider turn.
# absorb: no waiting task, so the tick would only spend tokens.
# dispatch: at least one waiting task. The agent still decides any reminder.
# Quota text is evidence for that agent, not a route picked here.
set -u

workspace=${SPYNEL_WORKSPACE:-.}
waiting="$workspace/.spynel/tasks/waiting"
found=
if [ -d "$waiting" ]; then
	for file in "$waiting"/*.md; do
		if [ -f "$file" ]; then
			found=1
			break
		fi
	done
fi

journal() {
	dir="$workspace/.spynel/extensions-state/quota"
	mkdir -p "$dir" 2>/dev/null && printf '%s %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$1" >>"$dir/journal" 2>/dev/null || true
}

if [ -z "$found" ]; then
	journal absorb
	printf 'absorb\n'
	exit 0
fi

journal dispatch
printf 'dispatch\n'
if command -v quota-axi >/dev/null 2>&1; then
	quota-axi --no-credential-refresh 2>/dev/null || printf 'quota-axi failed\n'
else
	printf 'quota-axi is not installed\n'
fi
exit 0
