import type { PackQuery } from '~/features/pack/api'

export type PackSort = 'new' | 'hot'
export type PackScope = 'all' | 'hot' | 'official'

/**
 * Discovery state lives in the URL so a filtered view can be shared, opened in
 * a new tab and restored by the back button. Nuxt's router is the only store.
 */
export const useDiscoveryFilters = () => {
  const route = useRoute()
  const router = useRouter()
  const showAdult = useShowAdultContent()

  const readString = (key: string, fallback = ''): string => {
    const raw = route.query[key]
    return typeof raw === 'string' ? raw : fallback
  }

  const scope = computed<PackScope>(() => {
    const raw = readString('scope', 'all')
    return raw === 'hot' || raw === 'official' ? raw : 'all'
  })
  const search = computed(() => readString('q'))
  const tag = computed(() => readString('tag'))
  // linked=1 keeps only packs that declare a game; work=<id> narrows to one.
  const linked = computed(() => readString('linked') === '1')
  const work = computed(() => {
    const parsed = Number.parseInt(readString('work'), 10)
    return Number.isFinite(parsed) && parsed > 0 ? parsed : undefined
  })
  const page = computed(() => {
    const parsed = Number.parseInt(readString('page', '1'), 10)
    return Number.isFinite(parsed) && parsed > 0 ? parsed : 1
  })

  const query = computed<PackQuery>(() => ({
    page: page.value,
    sort: scope.value === 'hot' ? 'hot' : 'new',
    official: scope.value === 'official',
    q: search.value || undefined,
    tag: tag.value || undefined,
    linked: linked.value || undefined,
    work: work.value,
    // The API hides R18 unless asked; the viewer asks once in settings.
    rating: showAdult.value ? 'all' : 'sfw'
  }))

  const update = (patch: Record<string, string | number | undefined>) => {
    const next: Record<string, string> = {}
    for (const [key, value] of Object.entries({ ...route.query, ...patch })) {
      if (value === undefined || value === '' || value === null) continue
      next[key] = String(value)
    }
    // Any change but paging returns to the first page: staying on page 7 of a
    // new filter usually lands on an empty screen.
    if (!('page' in patch)) delete next.page
    router.push({ query: next })
  }

  return { scope, search, tag, linked, work, page, query, update }
}
