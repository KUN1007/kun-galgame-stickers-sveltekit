import { resolveMultilingual } from '~/features/pack/types'

interface Envelope<T> {
  code: number
  data?: T
}

interface PackShape {
  title: Record<string, string | undefined>
  cover_url: string
  is_official: boolean
  content_rating: number
  sticker_count: number
  author?: { name: string }
  catalog_work?: { name: Record<string, string | undefined> }
  works?: { name: Record<string, string | undefined> }[]
}

/**
 * The share image for a pack, as a redirect.
 *
 * og:image points here rather than straight at nextmoe-og so the signing key
 * stays on the server, and so the address a page publishes never changes even
 * when the pack's title does -- the signature moves, this URL does not.
 *
 * Falls back to the pack's own cover when no OG key is configured: a square
 * sticker makes a worse share card than a composed one, but it is the subject.
 */
export default defineCachedEventHandler(
  async (event) => {
    const packId = getRouterParam(event, 'packId')
    const config = useRuntimeConfig()
    const fallback = `${config.public.siteUrl}/title.webp`

    const pack = await $fetch<Envelope<PackShape>>(
      `${config.apiBaseUrl}/api/v1/packs/${packId}`,
      { timeout: 8000 }
    ).catch(() => null)

    if (!pack || pack.code !== 0 || !pack.data) {
      return sendRedirect(event, fallback, 302)
    }
    const data = pack.data
    const locale = String(getQuery(event).locale ?? 'zh-cn')

    const badges: string[] = []
    if (data.is_official) badges.push('官方')
    if (data.content_rating === 1) badges.push('R18')
    badges.push(`${data.sticker_count} 张`)

    const game = data.catalog_work ?? data.works?.[0]
    const url =
      buildOgUrl('work', {
        title: ogText(resolveMultilingual(data.title, locale), 200) ?? 'Sticker pack',
        originalName: ogText(game ? resolveMultilingual(game.name, locale) : '', 200),
        cover: data.cover_url || undefined,
        label: ogText(data.author?.name, 80),
        badges: badges.slice(0, 4)
      }) ?? data.cover_url ?? fallback

    return sendRedirect(event, url, 302)
  },
  {
    name: 'og-pack',
    getKey: (event) => `${getRouterParam(event, 'packId')}:${getQuery(event).locale ?? 'zh-cn'}`,
    swr: true,
    maxAge: 60 * 60,
    staleMaxAge: 60 * 60 * 24
  }
)
