/**
 * Viewer preferences that must survive a reload and be known during SSR.
 *
 * A cookie rather than localStorage: the server renders the first page, and a
 * preference it cannot read would make the markup disagree with what the
 * client then decides -- for the R18 gate that means adult covers flashing on
 * screen before hydration takes them away.
 */
const YEAR_SECONDS = 60 * 60 * 24 * 365

export const useShowAdultContent = () => {
  const cookie = useCookie<string>('kun_sticker_r18', {
    maxAge: YEAR_SECONDS,
    sameSite: 'lax',
    default: () => '0'
  })

  return computed({
    get: () => cookie.value === '1',
    set: (value: boolean) => {
      cookie.value = value ? '1' : '0'
    }
  })
}
