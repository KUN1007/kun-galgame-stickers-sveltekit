import type { Author } from '~/features/pack/types'

export interface Comment {
  id: number
  post_number: number
  content_html: string
  content_raw: string
  created_at: string
  edited_at?: string
  author: Author
  can_edit: boolean
  can_delete: boolean
}

export interface CommentPage {
  thread_id: number
  comments: Comment[]
  total: number
  next_cursor?: string
  /** False when the community service is not configured for this deployment. */
  enabled: boolean
}

export const fetchComments = (packId: string, after?: string): Promise<CommentPage | null> =>
  kunFetchOrNull<CommentPage>(
    `/packs/${packId}/comments${after ? `?after=${encodeURIComponent(after)}` : ''}`
  )

export const addComment = (packId: string, body: string, replyTo?: number): Promise<Comment> =>
  kunFetch<Comment>(`/packs/${packId}/comments`, {
    method: 'POST',
    body: { body, reply_to: replyTo ?? 0 }
  })

export const editComment = (commentId: number, body: string): Promise<Comment> =>
  kunFetch<Comment>(`/comments/${commentId}`, { method: 'PATCH', body: { body } })

export const deleteComment = (commentId: number): Promise<unknown> =>
  kunFetch(`/comments/${commentId}`, { method: 'DELETE' })

export const MAX_COMMENT_LENGTH = 2000
