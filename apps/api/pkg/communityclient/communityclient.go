// Package communityclient reads and writes the nextmoe-infra community service.
//
// community is the platform's discussion primitive: one unit, a thread, whose
// anchor decides its shape. A sticker pack's comment section is the
// "entity-resource comments" shape -- anchor kind site_resource, anchor id the
// pack's uuid -- and `comments/resolve` gets-or-creates that thread, so this
// site never stores a thread id of its own.
//
// Auth is S2S Basic with the site's OAuth client credentials, and the tenant is
// NOT on the wire: community derives it from the calling client's
// oauth_clients.catalog_site binding, which is what stops one site writing into
// another's threads. The acting user is different -- community trusts this site
// to have authenticated them, so every write carries an author_id this site
// vouches for. That means author_id must never come from the request body a
// browser sent; it comes from the session.
package communityclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotConfigured = errors.New("community: not configured")
	ErrNotFound      = errors.New("community: not found")
	ErrForbidden     = errors.New("community: forbidden")
	ErrRateLimited   = errors.New("community: rate limited")
	ErrConflict      = errors.New("community: conflict")
	ErrUpstream      = errors.New("community: upstream unavailable")
)

func Missing(err error) bool { return errors.Is(err, ErrNotFound) }

// Anchor kinds, from the service's closed vocabulary. A sticker pack is a
// resource this site owns.
const (
	AnchorBoard         = 0
	AnchorSiteGame      = 1
	AnchorSiteResource  = 2
	AnchorCatalogWork   = 3
	AnchorCatalogPerson = 4
)

// Post statuses. Anything but visible is either awaiting review or gone, and
// this site shows neither.
const (
	StatusVisible = 0
	StatusHeld    = 1
	StatusDeleted = 2
)

type Config struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
}

type Client struct {
	http   *http.Client
	origin string
	auth   string
}

func New(cfg Config) *Client {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" || cfg.ClientID == "" || cfg.ClientSecret == "" {
		return &Client{}
	}
	return &Client{
		http:   &http.Client{Timeout: 8 * time.Second},
		origin: strings.TrimSuffix(base, "/api/v1/community"),
		auth:   base64.StdEncoding.EncodeToString([]byte(cfg.ClientID + ":" + cfg.ClientSecret)),
	}
}

func (c *Client) Configured() bool { return c != nil && c.origin != "" && c.auth != "" }

// envelope is the house shape; community uses it on every response including
// errors, so a non-zero code is the failure even when the status is 200.
type envelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	if !c.Configured() {
		return ErrNotConfigured
	}
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.origin+"/api/v1/community"+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+c.auth)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpstream, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("%w: read: %v", ErrUpstream, err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusConflict:
		return ErrConflict
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var env envelope[json.RawMessage]
		_ = json.Unmarshal(raw, &env)
		return fmt.Errorf("%w: %d %s", ErrUpstream, resp.StatusCode, env.Message)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("community decode: %w", err)
	}
	return nil
}

type Thread struct {
	ID                int64  `json:"id"`
	AnchorKind        int    `json:"anchor_kind"`
	AnchorID          string `json:"anchor_id"`
	PostsCount        int    `json:"posts_count"`
	ParticipantsCount int    `json:"participants_count"`
	Status            int    `json:"status"`
}

type Post struct {
	ID                int64  `json:"id"`
	ThreadID          int64  `json:"thread_id"`
	PostNumber        int    `json:"post_number"`
	AuthorID          int    `json:"author_id"`
	ContentHTML       string `json:"content_html"`
	ContentRaw        string `json:"content_raw"`
	Status            int    `json:"status"`
	ReplyToPostID     int64  `json:"reply_to_post_id"`
	CreatedAt         string `json:"created_at"`
	EditedAt          string `json:"edited_at"`
	EditedByModerator bool   `json:"edited_by_moderator"`
}

type ThreadWithPosts struct {
	Thread     Thread `json:"thread"`
	Posts      []Post `json:"posts"`
	NextCursor string `json:"next_cursor"`
}

// Resolve gets-or-creates the comments thread for an anchor and returns its
// first page. It is idempotent per anchor, so a reader arriving on a pack that
// nobody has commented on costs one call and creates one empty thread.
func (c *Client) Resolve(ctx context.Context, anchorKind int, anchorID string, contentRating int) (*ThreadWithPosts, error) {
	var env envelope[ThreadWithPosts]
	body := map[string]any{
		"anchor_kind":    anchorKind,
		"anchor_id":      anchorID,
		"content_rating": contentRating,
	}
	if err := c.do(ctx, http.MethodPost, "/comments/resolve", body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// PostList is what the posts lane returns: a page of posts, without the
// thread. resolve is the only lane whose data is a ThreadWithPosts.
type PostList struct {
	Posts      []Post `json:"posts"`
	NextCursor string `json:"next_cursor"`
}

func (c *Client) Posts(ctx context.Context, threadID int64, after string, limit int) (*PostList, error) {
	v := url.Values{}
	if after != "" {
		v.Set("after", after)
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	path := "/threads/" + strconv.FormatInt(threadID, 10) + "/posts"
	if q := v.Encode(); q != "" {
		path += "?" + q
	}
	var env envelope[PostList]
	if err := c.do(ctx, http.MethodGet, path, nil, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

func (c *Client) Reply(ctx context.Context, threadID int64, authorID int, body string, replyTo int64) (*Post, error) {
	payload := map[string]any{"author_id": authorID, "body": body}
	if replyTo > 0 {
		payload["reply_to_post_id"] = replyTo
	}
	// The write lanes wrap the post: data is {"post": …}, not the post itself.
	var env envelope[postWrapper]
	path := "/threads/" + strconv.FormatInt(threadID, 10) + "/posts"
	if err := c.do(ctx, http.MethodPost, path, payload, &env); err != nil {
		return nil, err
	}
	return &env.Data.Post, nil
}

type postWrapper struct {
	Post Post `json:"post"`
}

func (c *Client) Edit(ctx context.Context, postID int64, authorID int, body string) (*Post, error) {
	var env envelope[postWrapper]
	payload := map[string]any{"author_id": authorID, "body": body}
	if err := c.do(ctx, http.MethodPatch, "/posts/"+strconv.FormatInt(postID, 10), payload, &env); err != nil {
		return nil, err
	}
	return &env.Data.Post, nil
}

// Delete tombstones a post. author_id is a query param upstream: the request
// stays body-free, matching the artifact service's DELETE.
func (c *Client) Delete(ctx context.Context, postID int64, authorID int) error {
	path := fmt.Sprintf("/posts/%d?author_id=%d", postID, authorID)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}
