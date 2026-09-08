interface Author {
  id: number
}

interface Pack {
  id: string
  updated_at: string
  author: Author
}

interface Tag {
  slug: string
}

interface Envelope<T> {
  code: number
  data?: T
}

interface SitemapEntry {
  loc: string
  changefreq: 'weekly'
  priority: number
  lastmod?: string
  _i18nTransform: true
}

const PAGE_SIZE = 50
const MAX_PAGES = 40

const entry = (loc: string, priority: number, lastmod?: string): SitemapEntry => ({
  loc,
  changefreq: 'weekly',
  priority,
  ...(lastmod ? { lastmod } : {}),
  _i18nTransform: true
})

export default defineCachedEventHandler(
  async () => {
    const urls: SitemapEntry[] = [entry('/', 1), entry('/about', 0.6)]
    const apiBase = useRuntimeConfig().apiBaseUrl

    // Individual sticker pages are deliberately absent: they are linked from
    // every pack page, so crawlers still reach them, and listing them here grew
    // the file by one URL per sticker per locale with no ceiling.
    try {
      const authors = new Set<number>()

      for (let page = 1; page <= MAX_PAGES; page++) {
        const resp = await $fetch<Envelope<{ packs: Pack[]; total: number }>>(
          `${apiBase}/api/v1/packs`,
          { query: { page, limit: PAGE_SIZE }, timeout: 15000 }
        )
        const packs = resp?.code === 0 ? (resp.data?.packs ?? []) : []
        for (const pack of packs) {
          urls.push(entry(`/pack/${pack.id}`, 0.8, pack.updated_at))
          authors.add(pack.author.id)
        }
        if (packs.length < PAGE_SIZE) break
      }

      for (const uid of authors) urls.push(entry(`/u/${uid}`, 0.4))

      const tagResp = await $fetch<Envelope<Tag[]>>(`${apiBase}/api/v1/tags`, { timeout: 15000 })
      for (const tag of tagResp?.code === 0 ? (tagResp.data ?? []) : []) {
        urls.push(entry(`/tag/${tag.slug}`, 0.5))
      }
    } catch {
      // Static pages still ship from this handler; the rest retries on the next
      // cache miss.
    }

    return urls
  },
  {
    name: 'sitemap-urls',
    getKey: () => 'all',
    swr: true,
    maxAge: 60 * 60 * 6,
    staleMaxAge: 60 * 60 * 24 * 7
  }
)
