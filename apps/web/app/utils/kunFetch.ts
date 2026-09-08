interface KunApiResponse<T> {
  code: number
  message: string
  data: T
}

const SSR_FORWARDED = ['kun_oauth_access', 'kun_oauth_refresh', 'kun_oauth_user']

export class KunApiError extends Error {
  constructor(
    readonly code: number,
    message: string,
    readonly status = 0
  ) {
    super(message)
    this.name = 'KunApiError'
  }
}

const extractForwardedCookies = (cookieHeader?: string): string | undefined => {
  if (!cookieHeader) return undefined
  const kept: string[] = []
  for (const part of cookieHeader.split(';')) {
    const trimmed = part.trim()
    if (SSR_FORWARDED.some((name) => trimmed.startsWith(`${name}=`))) kept.push(trimmed)
  }
  return kept.length > 0 ? kept.join('; ') : undefined
}

const apiBase = (): string => {
  const config = useRuntimeConfig()
  const root = import.meta.server ? config.apiBaseUrl : config.public.apiBaseUrl
  return `${root}/api/v1`
}

interface KunFetchOptions extends Record<string, unknown> {
  headers?: HeadersInit
}

/**
 * Throws KunApiError on anything but code 0. Every failure used to collapse
 * into null, so a quota rejection, a permission error and a dead network were
 * indistinguishable and the UI could only ever say "nothing here".
 */
export const kunFetch = async <T>(path: string, options?: KunFetchOptions): Promise<T> => {
  const headers = new Headers(options?.headers)
  if (import.meta.server) {
    const forwarded = extractForwardedCookies(useRequestHeaders(['cookie']).cookie)
    if (forwarded) headers.set('cookie', forwarded)
  }

  let resp: KunApiResponse<T>
  try {
    resp = await $fetch<KunApiResponse<T>>(`${apiBase()}${path}`, {
      credentials: 'include',
      ...options,
      headers
    })
  } catch (error) {
    const wire = error as { data?: KunApiResponse<unknown>; status?: number }
    if (wire?.data && typeof wire.data.code === 'number') {
      throw new KunApiError(wire.data.code, wire.data.message, wire.status ?? 0)
    }
    throw new KunApiError(0, 'network', wire?.status ?? 0)
  }

  if (!resp || resp.code !== 0) {
    throw new KunApiError(resp?.code ?? 0, resp?.message ?? 'unknown')
  }
  return resp.data
}

/** For reads that should degrade into an empty state rather than an error page. */
export const kunFetchOrNull = async <T>(
  path: string,
  options?: KunFetchOptions
): Promise<T | null> => {
  try {
    return await kunFetch<T>(path, options)
  } catch {
    return null
  }
}
