# Filesystem Publication DOX

## Purpose

- Own cross-platform atomic replacement, no-clobber publication, and one-shot directory migration primitives.

## Local Contracts

- Write complete temporary content in the destination directory, apply requested private permissions, synchronize as required, and rename only after successful preparation.
- No-clobber publication must never overwrite a concurrently created target; platform-specific replacement files preserve equivalent observable semantics.
- Callers own content validation and higher-level rollback; this package owns filesystem publication integrity.
- Cross-filesystem directory migration copies into a destination sibling, publishes the complete tree atomically, and removes the source only after publication. Preserve regular-file modes and symlinks, and reject unsupported entry types.

## Child DOX Index

No child DOX files.
