package service

import (
	"context"
	stderrors "errors"
	"strings"

	"kun-galgame-sticker-api/internal/platform/sticker/dto"
	"kun-galgame-sticker-api/pkg/communityclient"
	"kun-galgame-sticker-api/pkg/errors"
	"kun-galgame-sticker-api/pkg/perm"
	"kun-galgame-sticker-api/pkg/userclient"

	"github.com/google/uuid"
)

const (
	// MaxCommentRunes is this site's own cap. Community has its own limits for
	// newcomers; this one keeps a comment a comment.
	MaxCommentRunes = 2000
	commentPageSize = 30
)

// PackComments returns a pack's comment thread. The thread is addressed by the
// pack, not by an id this site stores: community gets-or-creates it from the
// anchor, so there is nothing here to keep in sync.
func (s *Service) PackComments(
	ctx context.Context,
	packID uuid.UUID,
	after string,
	v Viewer,
) (*dto.CommentPage, *errors.AppError) {
	pack, appErr := s.visiblePack(packID, v)
	if appErr != nil {
		return nil, appErr
	}
	if !s.community.Configured() {
		return &dto.CommentPage{Comments: []dto.Comment{}, Enabled: false}, nil
	}

	thread, err := s.community.Resolve(ctx, communityclient.AnchorSiteResource, pack.ID.String(), int(pack.ContentRating))
	if err != nil {
		return nil, communityError(err)
	}
	posts := thread.Posts
	cursor := thread.NextCursor
	// resolve returns the first screen; a reader paging further asks by cursor.
	if after != "" {
		page, err := s.community.Posts(ctx, thread.Thread.ID, after, commentPageSize)
		if err != nil {
			return nil, communityError(err)
		}
		posts, cursor = page.Posts, page.NextCursor
	}

	return s.commentPage(ctx, thread.Thread, posts, cursor, v), nil
}

func (s *Service) commentPage(
	ctx context.Context,
	thread communityclient.Thread,
	posts []communityclient.Post,
	cursor string,
	v Viewer,
) *dto.CommentPage {
	authorIDs := make([]int, 0, len(posts))
	for _, post := range posts {
		authorIDs = append(authorIDs, post.AuthorID)
	}
	authors := s.users.Users(ctx, authorIDs)

	out := &dto.CommentPage{
		ThreadID:   thread.ID,
		Total:      thread.PostsCount,
		NextCursor: cursor,
		Comments:   make([]dto.Comment, 0, len(posts)),
		Enabled:    true,
	}
	for _, post := range posts {
		// A held or tombstoned post is not this site's to display: it is either
		// awaiting review upstream or deliberately gone.
		if post.Status != communityclient.StatusVisible {
			continue
		}
		author, ok := authors[post.AuthorID]
		if !ok {
			author = userclient.Placeholder(post.AuthorID)
		}
		out.Comments = append(out.Comments, dto.Comment{
			ID:          post.ID,
			PostNumber:  post.PostNumber,
			ContentHTML: post.ContentHTML,
			ContentRaw:  post.ContentRaw,
			CreatedAt:   post.CreatedAt,
			EditedAt:    post.EditedAt,
			Author:      dto.Author{ID: author.ID, Name: author.Name, Avatar: author.Avatar},
			CanEdit:     v.UID > 0 && post.AuthorID == v.UID,
			CanDelete:   v.UID > 0 && (post.AuthorID == v.UID || perm.Can(v.Roles, perm.PackDeleteAny)),
		})
	}
	return out
}

func (s *Service) AddComment(
	ctx context.Context,
	packID uuid.UUID,
	v Viewer,
	body string,
	replyTo int64,
) (*dto.Comment, *errors.AppError) {
	pack, appErr := s.visiblePack(packID, v)
	if appErr != nil {
		return nil, appErr
	}
	if !s.community.Configured() {
		return nil, errors.ErrCommunityUnavailable()
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.ErrInvalidParams("a comment cannot be empty")
	}
	if len([]rune(body)) > MaxCommentRunes {
		return nil, errors.ErrInvalidParams("comment is too long")
	}

	thread, err := s.community.Resolve(ctx, communityclient.AnchorSiteResource, pack.ID.String(), int(pack.ContentRating))
	if err != nil {
		return nil, communityError(err)
	}
	// author_id is the session's, never the request body's: community trusts
	// this site's assertion of who is speaking.
	post, err := s.community.Reply(ctx, thread.Thread.ID, v.UID, body, replyTo)
	if err != nil {
		return nil, communityError(err)
	}
	page := s.commentPage(ctx, thread.Thread, []communityclient.Post{*post}, "", v)
	if len(page.Comments) == 0 {
		// A newcomer's first posts are held for review upstream; the author is
		// told rather than shown a comment that is not there.
		return nil, errors.ErrCommentHeld()
	}
	return &page.Comments[0], nil
}

func (s *Service) EditComment(
	ctx context.Context,
	postID int64,
	v Viewer,
	body string,
) (*dto.Comment, *errors.AppError) {
	if !s.community.Configured() {
		return nil, errors.ErrCommunityUnavailable()
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.ErrInvalidParams("a comment cannot be empty")
	}
	if len([]rune(body)) > MaxCommentRunes {
		return nil, errors.ErrInvalidParams("comment is too long")
	}
	post, err := s.community.Edit(ctx, postID, v.UID, body)
	if err != nil {
		return nil, communityError(err)
	}
	page := s.commentPage(ctx, communityclient.Thread{}, []communityclient.Post{*post}, "", v)
	if len(page.Comments) == 0 {
		return nil, errors.ErrCommentHeld()
	}
	return &page.Comments[0], nil
}

func (s *Service) DeleteComment(ctx context.Context, postID int64, v Viewer) *errors.AppError {
	if !s.community.Configured() {
		return errors.ErrCommunityUnavailable()
	}
	if err := s.community.Delete(ctx, postID, v.UID); err != nil {
		return communityError(err)
	}
	return nil
}

// communityError keeps the upstream's distinctions: a reader must be able to
// tell "not yours" from "the service is down", because only one of them is
// worth retrying.
func communityError(err error) *errors.AppError {
	switch {
	case communityclient.Missing(err):
		return errors.ErrCommentNotFound()
	case stderrors.Is(err, communityclient.ErrForbidden):
		return errors.ErrNotOwner()
	case stderrors.Is(err, communityclient.ErrRateLimited):
		return errors.ErrCommentRateLimited()
	case stderrors.Is(err, communityclient.ErrConflict):
		return errors.ErrInvalidParams("this comment can no longer be changed")
	default:
		return errors.ErrCommunityUnavailable()
	}
}
