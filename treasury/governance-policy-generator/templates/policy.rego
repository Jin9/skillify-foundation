# policy.rego — agent auto-approve allowlist (SKELETON for human review)
#
# Primary access model: DEFAULT-DENY ALLOWLIST. The blacklist is supplementary
# (an emergency brake) and lives in KILLSWITCH.md, not here. An OS/sandbox
# floor (Layer 3) MUST exist beneath this Layer 2 policy.
#
# This policy MUST be deployed in infrastructure the agent cannot read or
# modify. It is a skeleton — every rule marked `# HUMAN REVIEW` requires an
# explicit human decision before the rule is enabled.
#
# Validate with: scripts/rego_deny_lint.py templates/policy.rego

package agent.governance

# --- DEFAULT-DENY POSTURE (DO NOT CHANGE) -----------------------------------
# Anything not explicitly allowed below is denied. Never set this to true.
default allow = false

# ---------------------------------------------------------------------------
# Each allow rule below is for ONE action classified high-reversibility /
# low-blast-radius in action-classification.md. Every rule body has at least
# one condition. No unconditioned allow rules are permitted.
# ---------------------------------------------------------------------------

# Reversible, low blast: read-only observation (always-eligible tier).
allow {
	input.action == "read"
	input.target_kind == "tracked_file"
}

# Reversible, low blast: local edit inside the working tree.
allow {
	input.action == "edit"
	input.target_kind == "tracked_file"
	not input.target_protected
}

# Reversible, low blast: git commit on a non-protected branch.
allow {
	input.action == "git_commit"
	input.branch != "main"
	not input.branch_protected
}

# Reversible, low blast: branch creation under the agent namespace.
allow {
	input.action == "git_branch_create"
	startswith(input.branch, "agent/")
}

# Semantic-versioning baseline: patch/minor dependency bump with CI green.
# # HUMAN REVIEW — confirm CI gate and that majors stay excluded.
allow {
	input.action == "dependency_update"
	input.semver_bump != "major"
	input.ci_status == "passing"
}

# Package installation — gated on a known-good registry allowlist.
# ~20% of agent-recommended package names are hallucinated (slopsquatting),
# so installation is NEVER auto-approved without registry verification.
# # HUMAN REVIEW — keep registry_allowlist curated and externally owned.
registry_allowlist := {
	# "<populate with known-good package names>",
}

allow {
	input.action == "package_install"
	registry_allowlist[input.package]
	input.registry == "trusted"
}

# ---------------------------------------------------------------------------
# Irreversible or medium/high-blast actions (deploy, external API call, DB
# migration, force push, secret access) have NO allow rule by design — they
# fall through to `default allow = false` and require human review until a
# tested, audited rollback story exists. DO NOT add allow rules for them here.
# ---------------------------------------------------------------------------
