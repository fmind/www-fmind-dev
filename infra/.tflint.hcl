// Only the terraform ruleset, which ships inside the tflint binary. The
// google ruleset is a plugin tflint downloads from the GitHub API, which would
// make `check:tofu` need network access and a token — the check stays offline.
plugin "terraform" {
  enabled = true
  preset  = "recommended"
}
