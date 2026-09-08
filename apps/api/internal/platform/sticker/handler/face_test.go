package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"kun-galgame-sticker-api/pkg/problem"

	"github.com/gofiber/fiber/v3"
)

// The face's parameter parsing gets its own test because the first version of
// it was silently wrong in a way no build or vet could see: the helpers
// returned the result of problem.Write, which is nil when the write succeeds,
// so `if err != nil` never fired. Every refused parameter wrote a problem
// document, then ran the query anyway and replaced that document with a
// success body -- a 400 status line over a 200 body. These cases pin the shape
// of a refusal, not just its status.

// probe mounts the helpers on their own app so a refusal can be observed
// without a database behind it. A helper that refuses must leave a problem
// document and nothing else; one that accepts must reach the handler body.
func probe(t *testing.T, url string) (int, string, problem.Problem) {
	t.Helper()
	app := fiber.New()
	app.Get("/packs", func(c fiber.Ctx) error {
		limit, bad := faceLimit(c)
		if bad != nil {
			return bad.write(c)
		}
		page, bad := facePage(c)
		if bad != nil {
			return bad.write(c)
		}
		if _, bad = faceSort(c); bad != nil {
			return bad.write(c)
		}
		if _, bad = faceRatingFilter(c); bad != nil {
			return bad.write(c)
		}
		return c.JSON(fiber.Map{"reached": true, "limit": limit, "page": page})
	})
	app.Get("/characters/:characterId", func(c fiber.Ctx) error {
		id, bad := faceCatalogID(c, "characterId")
		if bad != nil {
			return bad.write(c)
		}
		return c.JSON(fiber.Map{"reached": true, "id": id})
	})

	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, url, nil))
	if err != nil {
		t.Fatalf("request %s: %v", url, err)
	}
	defer res.Body.Close()

	var body struct {
		problem.Problem
		Reached bool `json:"reached"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode %s: %v", url, err)
	}
	if body.Reached && res.StatusCode != fiber.StatusOK {
		t.Fatalf("%s: status %d but the handler body ran anyway", url, res.StatusCode)
	}
	return res.StatusCode, res.Header.Get(fiber.HeaderContentType), body.Problem
}

func TestFaceRefusalsAreProblemDocuments(t *testing.T) {
	cases := []struct {
		name   string
		url    string
		status int
		code   string
	}{
		{"limit above the cap", "/packs?limit=999", fiber.StatusBadRequest, problem.CodeLimitTooLarge},
		{"limit not a number", "/packs?limit=many", fiber.StatusBadRequest, problem.CodeInvalidParameter},
		{"limit zero", "/packs?limit=0", fiber.StatusBadRequest, problem.CodeInvalidParameter},
		{"page zero", "/packs?page=0", fiber.StatusBadRequest, problem.CodeInvalidParameter},
		{"unknown sort", "/packs?sort=weird", fiber.StatusBadRequest, problem.CodeInvalidParameter},
		{"nsfw is a closed vocabulary", "/packs?nsfw=maybe", fiber.StatusBadRequest, problem.CodeInvalidParameter},
		{"catalog id must be positive", "/characters/0", fiber.StatusBadRequest, problem.CodeInvalidParameter},
		{"catalog id must be a number", "/characters/abc", fiber.StatusBadRequest, problem.CodeInvalidParameter},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, ctype, doc := probe(t, tc.url)
			if status != tc.status {
				t.Errorf("status = %d, want %d", status, tc.status)
			}
			if ctype != problem.ContentType {
				t.Errorf("content-type = %q, want %q", ctype, problem.ContentType)
			}
			if doc.Code != tc.code {
				t.Errorf("code = %q, want %q", doc.Code, tc.code)
			}
			if doc.Status != tc.status {
				t.Errorf("body status = %d, want %d (it must match the status line)", doc.Status, tc.status)
			}
			if doc.Type == "" || doc.Title == "" || doc.Detail == "" {
				t.Errorf("incomplete problem document: %+v", doc)
			}
		})
	}
}

func TestFaceAcceptsWhatItShould(t *testing.T) {
	for _, url := range []string{
		"/packs",
		"/packs?limit=50&page=1000",
		"/packs?sort=hot&nsfw=true",
		"/packs?sort=new&nsfw=false",
		"/characters/7935",
	} {
		if status, _, doc := probe(t, url); status != fiber.StatusOK {
			t.Errorf("%s: status = %d (%s), want 200", url, status, doc.Code)
		}
	}
}
