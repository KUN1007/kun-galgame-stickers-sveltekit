// Package problem writes RFC 9457 problem details.
//
// This is the error language of the public developer-platform faces, not of
// this site's own BFF: /api/v1 keeps the house {code, message, data} envelope,
// while /v1/sticker answers like catalog's /v2 does. A third party holding one
// nmk_ key across both faces should not have to decode two error dialects, and
// the gateway's own refusals (401/403/429 from ForwardAuth) are the only ones
// this service never gets to shape.
package problem

import (
	"github.com/gofiber/fiber/v3"
)

// Codes are reused verbatim from infra's closed registry
// (apps/api/internal/platform/apiv2/problem/registry.go) so a client can share
// one decoder. Only the subset a read-only face can actually produce is here.
const (
	CodeInvalidParameter   = "INVALID_PARAMETER"
	CodeLimitTooLarge      = "LIMIT_TOO_LARGE"
	CodeNotFound           = "NOT_FOUND"
	CodeInternalError      = "INTERNAL_ERROR"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// TypeBase is the platform's problem-type namespace. The URI is a stable
// identifier, not a promise that the page exists yet.
const TypeBase = "https://developer.nextmoe.dev/problems/sticker/"

type definition struct {
	title  string
	status int
	slug   string
}

var registry = map[string]definition{
	CodeInvalidParameter:   {"Invalid parameter", fiber.StatusBadRequest, "invalid-parameter"},
	CodeLimitTooLarge:      {"Limit too large", fiber.StatusBadRequest, "limit-too-large"},
	CodeNotFound:           {"Not found", fiber.StatusNotFound, "not-found"},
	CodeInternalError:      {"Internal error", fiber.StatusInternalServerError, "internal-error"},
	CodeServiceUnavailable: {"Service unavailable", fiber.StatusServiceUnavailable, "service-unavailable"},
}

// Problem is the wire shape. Field names and their meanings match infra's
// Problem so the two faces decode the same way; the members that only make
// sense for a write face (errors[], object, current_id) are left out rather
// than emitted permanently empty.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance"`
	Code     string `json:"code"`
}

// Write answers with a problem document. An unregistered code is an internal
// error rather than a half-filled document, because the alternative is a 200
// shaped like a failure.
func Write(c fiber.Ctx, code, detail string) error {
	def, ok := registry[code]
	if !ok {
		code = CodeInternalError
		def = registry[CodeInternalError]
	}
	// The content type is JSON's second argument, not a header set beforehand:
	// c.JSON overwrites Content-Type, which quietly turned every problem
	// document back into application/json.
	return c.Status(def.status).JSON(Problem{
		Type:     TypeBase + def.slug,
		Title:    def.title,
		Status:   def.status,
		Detail:   detail,
		Instance: instance(c),
		Code:     code,
	}, ContentType)
}

// ContentType is RFC 9457's media type. Exported because a caller that builds
// its own document (a test, most often) should not spell it again.
const ContentType = "application/problem+json"

// instance carries the query string: for a read face, which filter combination
// failed is most of the diagnostic value.
func instance(c fiber.Ctx) string {
	path := c.Path()
	if raw := string(c.Request().URI().QueryString()); raw != "" {
		return path + "?" + raw
	}
	return path
}
