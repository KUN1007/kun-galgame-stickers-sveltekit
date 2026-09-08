import { KunApiError } from '~/utils/kunFetch'

/**
 * The API answers with a distinct code per outcome so the UI can say what
 * actually went wrong. Anything unmapped falls back to the server's own
 * message, which is English but still more useful than "something failed".
 */
const messageKeys: Record<number, string> = {
  205: 'error.unauthorized',
  90001: 'error.packNotFound',
  90002: 'error.stickerNotFound',
  90003: 'error.notOwner',
  90004: 'error.packLimit',
  90005: 'error.stickerLimit',
  90006: 'error.uploadDailyLimit',
  90007: 'error.fileTooLarge',
  90008: 'error.packNeedsSticker',
  90009: 'error.imageRejected',
  90010: 'error.moderation',
  90011: 'error.imageUnavailable',
  90012: 'error.invalidParams',
  90013: 'error.tagLimit'
}

export const useApiError = () => {
  const { t, te } = useI18n()

  return (error: unknown): string => {
    if (error instanceof KunApiError) {
      const key = messageKeys[error.code]
      if (key && te(key)) return t(key)
      if (error.message === 'network') return t('error.network')
      return error.message || t('error.unknown')
    }
    return t('error.unknown')
  }
}

/** Runs a mutation and surfaces any failure as a toast. */
export const useMutation = () => {
  const describe = useApiError()

  return async <T>(run: () => Promise<T>): Promise<T | null> => {
    try {
      return await run()
    } catch (error) {
      useKunMessage(describe(error), 'error')
      return null
    }
  }
}
