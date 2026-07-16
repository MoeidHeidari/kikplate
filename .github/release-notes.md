# Release Notes: v0.22.0

**Previous:** v0.21.0
**Summary:** CLI improvements, SEO enhancements, sync resilience, and dark mode refinements

## Minor Features & Improvements

### CLI Enhancements (feat/cli-seo-sync-improvements)

**Search Command**
* Added positional argument support: `kik search <query>` now equivalent to `kik search --name <query>`
* Improved table output: removed duplicate NAME column, added VISIBILITY and SCHEMA columns
* SCHEMA column detects whether a plate has variables defined in synced metadata

**Describe Command**
* Now fetches and displays plate schema, modules, and files in formatted tables
* Removed redundant Rating field (Stars is sufficient)
* Schema table dynamically shows VALUES column only when at least one field has enum restrictions

**Generate Error Messages**
* Fixed error matching: switched from `==` to `errors.Is()` to properly surface wrapped errors
* `BuildZip` errors now wrapped with `ErrTemplateFailed` for accurate HTTP status codes and error messages
* Default error case now returns actual error text instead of generic "generation failed"

**Config Management**
* Added `kik config set-server <address>` subcommand for updating server address via CLI
* Works even if config file doesn't exist yet

Files Changed: 10 files
Impact: CLI UX, error reporting

### API Sync Improvements

**Metadata Capture**
* Added `Schema` field to `syncPlateManifestYAML` with proper `yaml` and `json` tags
* Schema variables now captured during plate sync for CLI and frontend use
* Explicit `json` tags on all `syncPlateManifestYAML` fields for consistent serialization

Files Changed: 1 file
Impact: Sync worker, plate metadata storage

### Web Frontend

**SEO & Discoverability**
* Added `generateMetadata` to plate detail pages with unique title, description, keywords, Open Graph tags, and canonical URLs
* Added `generateMetadata` to explore page with dynamic titles reflecting active filters
* Root layout now uses `title.template` for proper page title inheritance
* Added `/sitemap.xml` generation from all plate slugs with pagination support
* Added `/robots.txt` with crawler directives and sitemap reference
* Replaced default Next.js favicon with kikplate logo (app/icon.svg)

**Form Improvements**
* Fixed organization select in submit form: now displays org name instead of UUID when non-personal organization selected
* Replaced async `SelectValue` with explicit display lookup to eliminate race conditions

**Design System**
* Aligned all dark mode CSS variables to zinc oklch values for consistency
* Updated dark mode tokens: `card`, `popover`, `secondary`, `muted`, `accent`, `ring`, `sidebar`
* Changed overlay color in dialog and sheet components from `slate-950` to `zinc-950`
* Updated community badge tier color from `gray-*` to `zinc-*`

Files Changed: 13 files
Impact: Frontend UX, SEO, dark mode consistency

---

# Release Notes: v0.18.0

**Previous:** v0.17.0
**Summary:** 132 files changed, 5232 insertions(+), 284 deletions(-)

## Major Features

### Plate Generator System (feat/plate-generator-poc)

**Status:** Core foundation for template generation

* Added typed coercion for dynamic template rendering
* Expression conditions for conditional generation
* Template helpers for common operations
* Path safety for secure file operations
* New `/api/handler/handlers/generator_handler.go` and `/api/handler/routes/generator_routes.go`
* New model: `/api/model/generation.go`
* Database persistence layer: `/api/repository/postgres/generation.go`
* Dependency updates: `api/go.mod` and `api/go.sum`

Files Changed: 14 files
Impact: Backend API, database schema, repository layer

### Public Generate Endpoints (feat/public-generate-endpoints)

**Status:** Generator endpoints now accessible without authentication

* Made generate endpoints public (removed auth middleware requirement)
* Enhanced template helper utilities in `pkg/generator/template.go`
* Serves generation functionality for unauthenticated users

Files Changed: 4 files
Impact: API accessibility, middleware configuration

### Unified Manifest System (feat/unify-manifest-plate-yaml)

**Status:** plate.yaml is now primary manifest, kikplate.yaml as fallback

* Made `plate.yaml` the primary manifest file for submit/verify/sync operations
* Maintains backward compatibility with `kikplate.yaml` as fallback
* Updated CLI commands: `scaf.go`, `submit.go`, `verify.go`
* Updated service layer: `github_yaml.go`, `service.go`
* Updated error handling: `api/service/plate/errors.go`
* Documentation updates across README and docs files

Files Changed: 12 files
Impact: CLI, API service layer, documentation

### Password Recovery Flow (feat/password-recovery-throttling)

**Status:** Complete end-to-end password reset with email delivery

Backend:
* New model: `/api/model/password_reset.go` with token management
* Repository: `/api/repository/postgres/password_reset.go`
* Auth service methods: `RequestPasswordReset`, `ResetPassword`
* HTTP handlers and routes for both endpoints
* Server-side throttling: 1 request/min cooldown + 5 requests/hour max
* Email service integration with HTML template
* Event-driven email delivery pattern

Frontend:
* New pages: `/web/app/(auth)/forgot-password/page.tsx`, `/web/app/(auth)/reset-password/page.tsx`
* Form components: `ForgotPasswordForm.tsx`, `ResetPasswordForm.tsx`
* Hooks: `useRequestPasswordReset()`, `useResetPassword()`
* "Forgot password?" link on login page
* Auto-redirect to login on successful reset

Configuration:
* Updated logo URL to public absolute path for email rendering

Files Changed: 26 files
Impact: Backend auth layer, email service, frontend UI, database schema

### UI Redesign & Generation Support (feat/ui-plate-schema)

**Status:** Enhanced plate display and user interface

* Redesigned UI components for plate viewing
* Added generation support in UI
* Updated plate content and header tabs
* Enhanced submit form UI
* Updated styles in `web/app/globals.css`
* Improved use modal and button components
* GitHub client integration updates
* Hero search component enhancements

Files Changed: 14 files
Impact: Frontend UI/UX, user experience

### CLI Rebranding (feat/rename-cli-to-kik)

**Status:** CLI renamed from `kikplate` to `kik`

* Binary name changed from `kikplate` to `kik`
* Updated GoReleaser configuration with new binary and archive naming
* Updated Homebrew formula: `homebrew-kik` repository
* Updated Scoop bucket: `scoop-kik` repository
* Installation: `brew tap kikplate/kik && brew install kik`
* Updated all documentation and CLI examples (40+ references)
* Config directory: `~/.kik/` (instead of `~/.kikplate/`)
* GitHub workflows updated for binary naming

Files Changed: 5 files
Impact: Package manager configurations, documentation, distribution

## Bug Fixes

### Lifecycle Operations Fix (fix/lifeCyscleOps)

* Fixed lifecycle operation handling in `api/service/plate/lifecycle_ops.go`
* Fixed review badge operations in `api/service/plate/review_badge_ops.go`

Files Changed: 2 files

## Documentation Updates

* Updated CLI reference documentation
* Architecture documentation enhanced
* Configuration documentation synchronized
* Database schema documentation updated
* How-it-works guide updated with new features
* Generator tutorial added
* Updated README with feature highlights

Files Changed: 7+ documentation files

## Merge Commits

* #196: Merged CLI rename branch (feat/rename-cli-to-kik)
* #195: Merged password recovery feature (feat/password-recovery-throttling)
* #194: Merged UI redesign (feat/ui-plate-schema)
* #193: Merged manifest unification (feat/unify-manifest-plate-yaml)
* #191: Merged public generate endpoints (feat/public-generate-endpoints)
* #190: Merged generator POC (feat/plate-generator-poc)

## Testing Status

* All GitHub Actions CI checks passing
* API compilation successful
* Frontend linting passed
* No pre-existing test regressions

## Database Schema Changes

* New table: `password_resets` (token-based password recovery)
* New table: `generations` (template generation tracking)

## Breaking Changes

* CLI command changed from `kikplate` to `kik`
* Config directory moved from `~/.kikplate/` to `~/.kik/`
* `plate.yaml` now primary manifest (kikplate.yaml still supported)

## Migration Notes

For users upgrading to this release:

1. Update installation: `brew install kik` (or equivalent for your platform)
2. Config will auto-migrate from `~/.kikplate/` to `~/.kik/`
3. Existing `kikplate.yaml` files will still work but `plate.yaml` is preferred
4. Password recovery is now available without re-authentication

## Files Modified Summary

Backend (API): 48 files
Frontend (Web): 26 files
CLI: 5 files
Workflows: 2 files
Documentation: 15 files
Configuration: 2 files
Package Managers: 2 files (formula/bucket files external repos)

## Recommendation

Ready for Release: All features complete, tested, and documented. No known blockers.
