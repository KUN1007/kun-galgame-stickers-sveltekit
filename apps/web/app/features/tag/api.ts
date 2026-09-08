import type { Tag } from '~/features/pack/types'

export const fetchTags = (): Promise<Tag[] | null> => kunFetchOrNull<Tag[]>('/tags')
