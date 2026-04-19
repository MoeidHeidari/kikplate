/**
 * Format a slug into a human-readable label.
 * @example formatExplorerCategoryLabel("backend") => "Backend"
 * @example formatExplorerCategoryLabel("backend-api") => "Backend API"
 * @example formatExplorerCategoryLabel("backend_api") => "Backend API"
 * @param slug 
 * @returns 
 */
export function formatExplorerCategoryLabel(slug: string): string {
  const s = slug.trim()
  if (!s) return ""
  return s
    .split(/[-_]/)
    .map((w) => (w ? w.charAt(0).toUpperCase() + w.slice(1).toLowerCase() : ""))
    .join(" ")
}
