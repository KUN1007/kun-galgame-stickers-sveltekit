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
  like_count: number
  is_liked: boolean
  /** What this answers, and the top-level comment the exchange hangs under. */
  reply_to?: number
  root_id?: number
  reply_to_name?: string
}

/** A top-level comment with the replies that hang under it. */
export interface CommentNode extends Comment {
  replies: Comment[]
}

export interface LikeResult {
  liked: boolean
  like_count: number
}

/**
 * community's closed vocabulary. The numbers are the wire values, so they
 * cannot be reordered to suit a menu.
 */
export const FLAG_REASONS = [0, 1, 2, 4, 3] as const
export type FlagReason = (typeof FLAG_REASONS)[number]

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

export const toggleCommentLike = (commentId: number): Promise<LikeResult> =>
  kunFetch<LikeResult>(`/comments/${commentId}/like`, { method: 'POST' })

export const reportComment = (commentId: number, reason: FlagReason, note: string): Promise<unknown> =>
  kunFetch(`/comments/${commentId}/report`, { method: 'POST', body: { reason, note } })

/**
 * Groups a flat page into one level of nesting. community threads arbitrarily
 * deep -- a reply to a reply keeps the root of the exchange -- but a comment
 * section under a sticker pack reads better as "comment, then answers to it"
 * than as an indent ladder, so everything under a root sits at one depth and
 * says who it answered.
 */
export const nestComments = (comments: Comment[]): CommentNode[] => {
  const roots: CommentNode[] = []
  const byId = new Map<number, CommentNode>()
  for (const comment of comments) {
    if (!comment.root_id) {
      const node = { ...comment, replies: [] }
      byId.set(comment.id, node)
      roots.push(node)
    }
  }
  for (const comment of comments) {
    if (!comment.root_id) continue
    const parent = byId.get(comment.root_id)
    // A reply whose root is on an earlier page has nothing to nest under, so
    // it stands on its own rather than disappearing.
    if (parent) parent.replies.push(comment)
    else roots.push({ ...comment, replies: [] })
  }
  return roots
}

export const MAX_COMMENT_LENGTH = 2000
