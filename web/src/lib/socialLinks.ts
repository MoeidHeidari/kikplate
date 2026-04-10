export function getSocialLink(
  social: Array<{ type: string; link: string }> | undefined,
  type: string,
): string | undefined {
  const t = type.toLowerCase()
  const found = social?.find((s) => s.type.toLowerCase() === t)
  const link = found?.link?.trim()
  if (!link || link === "#") return undefined
  return link
}
