# Jean Roadmap

## Phase 1: Code Cleanup & Simplification

### Refactoring
- [ ] Extract common code between TUI and menubar versions
- [ ] Create shared `internal/usage` package for data fetching
- [ ] Add proper error handling with retries
- [ ] Add logging (configurable verbosity)
- [ ] Remove debug code and unused files

### Code Quality
- [ ] Add unit tests for cookie decryption
- [ ] Add unit tests for API client
- [ ] Add integration tests (mock API)
- [ ] Run `go vet`, `staticcheck`, `golangci-lint`
- [ ] Add godoc comments to exported functions

### Architecture
- [ ] Consider using interfaces for testability
- [ ] Add context support for graceful shutdown
- [ ] Proper signal handling (SIGTERM, SIGINT)

---

## Phase 2: Configuration

### Config File Support
- [ ] Create `~/.jean/config.yaml` support
- [ ] Configurable options:
  ```yaml
  refresh_interval: 30s
  progress_bar:
    style: parallelogram  # circle, square, block
    width: 5
  notifications:
    enabled: true
    threshold: 80  # notify at 80% usage
  ```

### CLI Flags
- [ ] `--config` - custom config path
- [ ] `--refresh` - override refresh interval
- [ ] `--verbose` - enable debug logging
- [ ] `--version` - show version info

### Environment Variables
- [ ] `JEAN_CONFIG` - config file path
- [ ] `JEAN_REFRESH_INTERVAL` - refresh interval
- [ ] `JEAN_LOG_LEVEL` - logging verbosity

---

## Phase 3: GitHub Repository

### Repository Setup
- [ ] Create GitHub repo (public or private initially)
- [ ] Add .gitignore (Go template)
- [ ] Add LICENSE (MIT)
- [ ] Add CONTRIBUTING.md
- [ ] Add CODE_OF_CONDUCT.md
- [ ] Set up branch protection rules

### Documentation
- [ ] Add screenshots/GIFs to README
- [ ] Create GitHub wiki for detailed docs
- [ ] Add CHANGELOG.md
- [ ] Add SECURITY.md for vulnerability reporting

### Community
- [ ] Issue templates (bug report, feature request)
- [ ] Pull request template
- [ ] Discussion board setup

---

## Phase 4: CI/CD Pipeline

### GitHub Actions Workflows

#### Build & Test
```yaml
# .github/workflows/ci.yml
- Run tests on push/PR
- Run linters (golangci-lint)
- Build for macOS (amd64, arm64)
- Upload build artifacts
```

#### Release
```yaml
# .github/workflows/release.yml
- Trigger on tag push (v*)
- Build universal binary
- Create GitHub release
- Generate changelog
- Upload binaries
```

### Quality Gates
- [ ] Require passing tests for merge
- [ ] Require linter pass
- [ ] Code coverage reporting (codecov)
- [ ] Dependency vulnerability scanning (dependabot)

---

## Phase 5: Release & Distribution

### Binary Distribution
- [ ] Universal macOS binary (amd64 + arm64)
- [ ] Compress with UPX (optional)
- [ ] GitHub Releases with checksums

### Homebrew
- [ ] Create Homebrew tap (`homebrew-jean`)
- [ ] Formula for `brew install jean`
- [ ] Cask for GUI app (optional)

### macOS App Bundle
- [ ] Create proper .app bundle
- [ ] Add Info.plist
- [ ] Add app icon
- [ ] DMG installer with drag-to-Applications

### Code Signing & Notarization
- [ ] Apple Developer account
- [ ] Code signing certificate
- [ ] Notarization for Gatekeeper
- [ ] Hardened runtime

---

## Phase 6: Open Sourcing

### Pre-Launch Checklist
- [ ] Security audit (no hardcoded secrets)
- [ ] Remove any personal/org-specific code
- [ ] Ensure all dependencies are MIT/Apache compatible
- [ ] Add license headers to source files
- [ ] Final README review

### Launch
- [ ] Make repository public
- [ ] Write launch blog post (optional)
- [ ] Post to Hacker News / Reddit
- [ ] Tweet announcement

### Post-Launch
- [ ] Monitor issues and discussions
- [ ] Respond to PRs
- [ ] Tag first public release (v1.0.0)

---

## Phase 7: Future Features

### Notifications
- [ ] Native macOS notifications
- [ ] Configurable thresholds (warn at 80%, alert at 95%)
- [ ] Sound alerts (optional)

### Multiple Accounts
- [ ] Support multiple Claude accounts
- [ ] Account switcher in dropdown
- [ ] Per-account usage display

### Historical Tracking
- [ ] Local SQLite database for history
- [ ] Usage trends over time
- [ ] Daily/weekly/monthly reports
- [ ] Export to CSV/JSON

### Claude Code Integration
- [ ] Monitor local Claude Code sessions
- [ ] OpenTelemetry metrics collection
- [ ] Per-session token tracking
- [ ] Cost estimation

### Advanced UI
- [ ] Preferences window
- [ ] Keyboard shortcuts
- [ ] Touch Bar support (if applicable)
- [ ] Widgets for macOS Notification Center

### Cross-Platform
- [ ] Linux support (systray)
- [ ] Windows support (systray)
- [ ] Web dashboard (optional)

---

## Version Milestones

| Version | Milestone | Target |
|---------|-----------|--------|
| v0.1.0 | Alpha - Basic functionality | ✅ Done |
| v0.2.0 | Beta - Config + cleanup | TBD |
| v0.3.0 | RC - CI/CD + distribution | TBD |
| v1.0.0 | Stable - Open source release | TBD |
| v1.1.0 | Notifications | TBD |
| v1.2.0 | Historical tracking | TBD |
| v2.0.0 | Claude Code integration | TBD |
