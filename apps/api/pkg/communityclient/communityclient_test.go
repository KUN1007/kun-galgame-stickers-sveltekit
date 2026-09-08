package communityclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The three lanes do not share a payload shape, and getting it wrong fails
// silently: an unwrapped decode yields a zero-valued post that still passes a
// status==visible check, so a reply looks accepted and comes back blank. These
// bodies are copied from the live service.
func TestResponseShapesPerLane(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/community/comments/resolve":
			// resolve is the one lane whose data IS the thread-with-posts.
			_, _ = w.Write([]byte(`{"code":0,"message":"成功","data":{
				"thread":{"id":57182,"anchor_kind":2,"anchor_id":"pack-uuid","posts_count":3},
				"posts":[{"id":11048,"post_number":1,"author_id":3,"content_html":"<p>hi</p>","status":0}],
				"next_cursor":"1"}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/community/threads/57182/posts":
			_, _ = w.Write([]byte(`{"code":0,"message":"成功","data":{"post":{
				"id":11049,"thread_id":57182,"post_number":2,"author_id":3,
				"content_raw":"hey","content_html":"<p>hey</p>","status":0}}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/community/threads/57182/posts":
			_, _ = w.Write([]byte(`{"code":0,"message":"成功","data":{
				"posts":[{"id":11050,"post_number":3,"author_id":4,"content_html":"<p>p3</p>","status":0}],
				"next_cursor":"3"}}`))
		case r.Method == http.MethodPatch:
			_, _ = w.Write([]byte(`{"code":0,"message":"成功","data":{"post":{"id":11049,"content_html":"<p>edited</p>","status":0}}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, ClientID: "id", ClientSecret: "secret"})
	if !c.Configured() {
		t.Fatal("client should be configured")
	}
	ctx := context.Background()

	thread, err := c.Resolve(ctx, AnchorSiteResource, "pack-uuid", 0)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if thread.Thread.ID != 57182 || len(thread.Posts) != 1 || thread.Posts[0].ContentHTML != "<p>hi</p>" {
		t.Errorf("resolve decoded wrong: %+v", thread)
	}

	post, err := c.Reply(ctx, 57182, 3, "hey", 0)
	if err != nil {
		t.Fatalf("reply: %v", err)
	}
	if post.ID != 11049 || post.ContentHTML != "<p>hey</p>" {
		t.Errorf("reply decoded wrong (data is {\"post\": …}): %+v", post)
	}

	page, err := c.Posts(ctx, 57182, "2", 30)
	if err != nil {
		t.Fatalf("posts: %v", err)
	}
	if len(page.Posts) != 1 || page.Posts[0].ID != 11050 || page.NextCursor != "3" {
		t.Errorf("posts decoded wrong (data is {\"posts\": …}): %+v", page)
	}

	edited, err := c.Edit(ctx, 11049, 3, "edited")
	if err != nil {
		t.Fatalf("edit: %v", err)
	}
	if edited.ContentHTML != "<p>edited</p>" {
		t.Errorf("edit decoded wrong: %+v", edited)
	}
}

// A site whose client carries no tenant binding is refused on every lane, and
// that has to stay distinguishable from the service being down: one is a
// configuration mistake, the other is worth retrying.
func TestForbiddenAndRateLimitAreDistinct(t *testing.T) {
	for _, tc := range []struct {
		status int
		want   error
	}{
		{http.StatusForbidden, ErrForbidden},
		{http.StatusTooManyRequests, ErrRateLimited},
		{http.StatusNotFound, ErrNotFound},
		{http.StatusConflict, ErrConflict},
	} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte(`{"code":1,"message":"nope"}`))
		}))
		_, err := New(Config{BaseURL: srv.URL, ClientID: "id", ClientSecret: "s"}).
			Resolve(context.Background(), AnchorSiteResource, "x", 0)
		if err != tc.want {
			t.Errorf("status %d gave %v, want %v", tc.status, err, tc.want)
		}
		srv.Close()
	}
}

func TestUnconfiguredClientRefusesRatherThanCallingNowhere(t *testing.T) {
	c := New(Config{})
	if c.Configured() {
		t.Fatal("an empty config must not report configured")
	}
	if _, err := c.Resolve(context.Background(), AnchorSiteResource, "x", 0); err != ErrNotConfigured {
		t.Errorf("got %v, want ErrNotConfigured", err)
	}
}
