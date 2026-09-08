import { createHmac } from 'node:crypto'

/**
 * Builds a signed nextmoe-og URL.
 *
 * The contract is frozen upstream: the signature covers `template + "\n" + d`,
 * not `d` alone, and `d` is base64url of the field JSON. Identical fields in an
 * identical key order produce an identical URL, which is also the render cache
 * key -- so changing a pack's title changes the URL and therefore the image.
 *
 * Returns null when no site key is configured, which is the caller's cue to
 * fall back to the artwork the page already has.
 */
export const buildOgUrl = (
  template: string,
  fields: Record<string, unknown>
): string | null => {
  const { ogBaseUrl, ogSiteKey } = useRuntimeConfig()
  if (!ogSiteKey) return null

  // Upstream rejects an over-long field rather than truncating it, and an
  // undefined one would still change the URL for no visible gain.
  const clean = Object.fromEntries(
    Object.entries(fields).filter(
      ([, value]) =>
        value !== undefined &&
        value !== null &&
        value !== '' &&
        !(Array.isArray(value) && value.length === 0)
    )
  )
  const d = Buffer.from(JSON.stringify(clean), 'utf8').toString('base64url')
  const sig = createHmac('sha256', ogSiteKey).update(`${template}\n${d}`).digest('base64url')
  return `${ogBaseUrl.replace(/\/$/, '')}/v1/og/${template}?d=${d}&sig=${sig}`
}

/** Upstream caps every string field; going over is a 400, not a trim. */
export const ogText = (value: string | undefined, max: number): string | undefined => {
  const trimmed = (value ?? '').trim()
  if (!trimmed) return undefined
  const runes = [...trimmed]
  return runes.length > max ? runes.slice(0, max - 1).join('') + '…' : trimmed
}
