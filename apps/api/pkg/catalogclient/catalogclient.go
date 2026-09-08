// Package catalogclient reads the nextmoe-infra catalog /v2 face.
//
// catalog is the platform's identity registry: it answers "which work is this"
// and "which character is this", and owns nothing else -- covers, intros and
// ratings belong to the product sites. This site stores the ids it resolves
// plus a display snapshot, so the client is only needed while a user is
// picking, and for the two detail pages that show more than the snapshot.
//
// Auth is an application key (nmk_live_ / nmk_test_) sent as a bearer token.
// Responses are bare JSON -- there is no house envelope here; errors are RFC
// 9457 problem documents.
package catalogclient

import (
	"context"
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
	ErrNotConfigured = errors.New("catalog: no api key configured")
	ErrNotFound      = errors.New("catalog: not found")
	ErrUpstream      = errors.New("catalog: upstream unavailable")
)

// Missing reports that catalog answered about the object itself. Every other
// error is transport-level and says nothing about whether the object exists --
// caching a 429 as "no such character" would blank a page for a whole TTL.
func Missing(err error) bool { return errors.Is(err, ErrNotFound) }

type Config struct {
	BaseURL string
	APIKey  string
}

type Client struct {
	http   *http.Client
	origin string
	apiKey string
}

func New(cfg Config) *Client {
	return &Client{
		http:   &http.Client{Timeout: 8 * time.Second},
		origin: origin(cfg.BaseURL),
		apiKey: strings.TrimSpace(cfg.APIKey),
	}
}

// origin trims a version suffix so both "https://api.nextmoe.dev" and
// ".../v2" configure the same client; every path below carries its own /v2.
func origin(raw string) string {
	u := strings.TrimRight(strings.TrimSpace(raw), "/")
	for _, suffix := range []string{"/api/v1", "/v2", "/v1"} {
		if strings.HasSuffix(u, suffix) {
			return strings.TrimSuffix(u, suffix)
		}
	}
	return u
}

func (c *Client) Configured() bool {
	return c != nil && c.origin != "" && c.apiKey != ""
}

type problem struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Code   string `json:"code"`
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	if !c.Configured() {
		return ErrNotConfigured
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.origin+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpstream, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("%w: read: %v", ErrUpstream, err)
	}
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		var p problem
		_ = json.Unmarshal(body, &p)
		detail := p.Detail
		if detail == "" {
			detail = p.Title
		}
		return fmt.Errorf("%w: %d %s %s", ErrUpstream, resp.StatusCode, p.Code, detail)
	}
	if out == nil || len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("catalog decode: %w", err)
	}
	return nil
}

type List[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor"`
}

type LocalizedText struct {
	Value     string `json:"value"`
	IsMachine bool   `json:"is_machine"`
}

type Image struct {
	URL       string `json:"url"`
	Hash      string `json:"hash"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Thumbhash string `json:"thumbhash"`
}

type Work struct {
	ID            string                   `json:"id"`
	Medium        string                   `json:"medium"`
	DisplayName   string                   `json:"display_name"`
	Latin         *string                  `json:"latin"`
	Localized     map[string]LocalizedText `json:"localized"`
	OLang         string                   `json:"olang"`
	ContentRating string                   `json:"content_rating"`
	ReleaseDate   *string                  `json:"release_date"`
	Cover         *Image                   `json:"cover"`
}

type Character struct {
	ID          string                   `json:"id"`
	DisplayName string                   `json:"display_name"`
	Latin       *string                  `json:"latin"`
	Localized   map[string]LocalizedText `json:"localized"`
	Gender      *string                  `json:"gender"`
	Birthday    *string                  `json:"birthday"`
	BloodType   *string                  `json:"blood_type"`
	RosterRole  string                   `json:"roster_role"`
	Image       *Image                   `json:"image"`
	Traits      *[]Trait                 `json:"traits"`
	Aliases     *[]Alias                 `json:"aliases"`
}

type Trait struct {
	ID             string                   `json:"id"`
	DisplayName    string                   `json:"display_name"`
	Group          string                   `json:"group"`
	Localized      map[string]LocalizedText `json:"localized"`
	GroupLocalized map[string]LocalizedText `json:"group_localized"`
	IsSexual       bool                     `json:"is_sexual"`
}

type Alias struct {
	Name  string `json:"name"`
	Latin string `json:"latin"`
}

type Appearance struct {
	Work       Work   `json:"work"`
	RosterRole string `json:"roster_role"`
}

const searchLimit = 12

// SearchWorks backs the editor's game picker. Only galgame-shaped media are
// useful here, but catalog has no medium filter on the search lane, so the
// caller sees whatever the index ranks highest.
func (c *Client) SearchWorks(ctx context.Context, q string, limit int) ([]Work, error) {
	if limit <= 0 || limit > 50 {
		limit = searchLimit
	}
	v := url.Values{}
	v.Set("q", q)
	v.Set("limit", strconv.Itoa(limit))
	// On a collection lane localized names are gated behind include=titles;
	// without it every hit comes back with only its canonical display name and
	// a Chinese reader picks games off a list of Japanese titles.
	v.Set("include", "titles")
	var page List[Work]
	if err := c.get(ctx, "/v2/catalog/works?"+v.Encode(), &page); err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (c *Client) Work(ctx context.Context, id int64) (*Work, error) {
	var out Work
	if err := c.get(ctx, "/v2/catalog/works/"+strconv.FormatInt(id, 10)+"?include=titles", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WorkCharacters returns a work's roster. The roster is small enough that the
// picker wants it whole, so this follows next_cursor instead of exposing
// pagination -- capped, because an unbounded follow would hand a hostile or
// broken upstream an open loop.
func (c *Client) WorkCharacters(ctx context.Context, id int64) ([]Character, error) {
	const perPage, maxPages = 100, 5
	out := make([]Character, 0, perPage)
	cursor := ""
	for page := 0; page < maxPages; page++ {
		v := url.Values{}
		v.Set("limit", strconv.Itoa(perPage))
		if cursor != "" {
			v.Set("cursor", cursor)
		}
		var got List[Character]
		if err := c.get(ctx, "/v2/catalog/works/"+strconv.FormatInt(id, 10)+"/characters?"+v.Encode(), &got); err != nil {
			return nil, err
		}
		out = append(out, got.Items...)
		if got.NextCursor == "" {
			break
		}
		cursor = got.NextCursor
	}
	return out, nil
}

func (c *Client) CharacterDetail(ctx context.Context, id int64) (*Character, error) {
	v := url.Values{}
	v.Set("view", "full")
	v.Set("include", "image,traits,aliases")
	var out Character
	if err := c.get(ctx, "/v2/catalog/characters/"+strconv.FormatInt(id, 10)+"?"+v.Encode(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WorksByIDs is the batch lane: up to 100 ids, no pagination. It fills in the
// covers that the appearances lane leaves null, in one request rather than one
// per work.
func (c *Client) WorksByIDs(ctx context.Context, ids []string) (map[string]Work, error) {
	out := make(map[string]Work, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	if len(ids) > 100 {
		ids = ids[:100]
	}
	v := url.Values{}
	v.Set("ids", strings.Join(ids, ","))
	v.Set("include", "titles")
	var page List[Work]
	if err := c.get(ctx, "/v2/catalog/works?"+v.Encode(), &page); err != nil {
		return nil, err
	}
	for _, work := range page.Items {
		out[work.ID] = work
	}
	return out, nil
}

func (c *Client) CharacterAppearances(ctx context.Context, id int64, limit int) ([]Appearance, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	v := url.Values{}
	v.Set("limit", strconv.Itoa(limit))
	var page List[Appearance]
	if err := c.get(ctx, "/v2/catalog/characters/"+strconv.FormatInt(id, 10)+"/appearances?"+v.Encode(), &page); err != nil {
		return nil, err
	}
	return page.Items, nil
}
