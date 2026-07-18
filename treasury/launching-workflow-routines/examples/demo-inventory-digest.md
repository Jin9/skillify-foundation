---
routine: demo-inventory-digest
summary: Two-node demo - inventory a directory, then digest the inventory into one paragraph. Dependency-free smoke test for the launcher.
version: 0.1.0
default_on_fail: stop
tags: demo, smoke-test
---

Smallest possible routine: proves the format, the gates, and the verification loop without dispatching any other skill. Launch it to smoke-test the launcher on any host.

## Inputs

- target_dir: absolute path of the directory to inventory (required)

## Nodes

### Node: inventory

- purpose: List the files in the target directory with sizes, newest first
- executor: agent-inline
- tier: small
- inputs: user.target_dir
- outputs: 01-inventory.md
- gate: after

### Node: digest

- purpose: Compress the inventory into a one-paragraph digest of what the directory contains and what stands out
- executor: agent-inline
- tier: small
- inputs: 01-inventory.md
- outputs: 02-digest.md
- gate: after
