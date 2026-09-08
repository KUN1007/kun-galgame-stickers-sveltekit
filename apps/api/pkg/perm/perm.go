package perm

import "slices"

type Permission string

// Moderation and admin surfaces are not built yet; this is the vocabulary they
// will be built against, and the owner checks in the sticker service already
// consult it so a moderator is not treated as a stranger to someone else's pack.
const (
	PackEditAny    Permission = "pack.edit_any"
	PackDeleteAny  Permission = "pack.delete_any"
	PackHide       Permission = "pack.hide"
	PackViewHidden Permission = "pack.view_hidden"

	StickerEditAny   Permission = "sticker.edit_any"
	StickerDeleteAny Permission = "sticker.delete_any"

	TagManage Permission = "tag.manage"

	AdminDashboard Permission = "admin.dashboard"
)

const (
	RoleCreator   = "creator"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
	RoleRen       = "ren"
)

var moderatorPerms = []Permission{
	PackEditAny,
	PackHide,
	PackViewHidden,
	StickerEditAny,
	StickerDeleteAny,
	TagManage,
}

var adminPerms = append(append([]Permission{}, moderatorPerms...), PackDeleteAny, AdminDashboard)

var bundles = map[string][]Permission{
	RoleModerator: moderatorPerms,
	RoleAdmin:     adminPerms,
	RoleRen:       adminPerms,
}

func Can(roles []string, p Permission) bool {
	for _, name := range roles {
		if slices.Contains(bundles[name], p) {
			return true
		}
	}
	return false
}

func Has(roles []string, name string) bool { return slices.Contains(roles, name) }

// Union folds the OP's site_roles into roles. The OP returns them separately:
// roles are global, site_roles are granted on this site only.
func Union(roles, siteRoles []string) []string {
	if len(siteRoles) == 0 {
		return roles
	}
	out := slices.Clone(roles)
	for _, r := range siteRoles {
		if !slices.Contains(out, r) {
			out = append(out, r)
		}
	}
	return out
}
