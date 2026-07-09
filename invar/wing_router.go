// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package invar

// RBAC role string
const (
	WRoleSuper = "super-admin"   // Super admin, auto add.
	WRoleAdmin = "admin"         // Normal admin, same time as super admin permissions.
	WRoleUser  = "user"          // Default normal user.
	WRoleMComp = "mall-comp"     // Access role of Mall  on 'comp' api router.
	WRoleSComp = "store-comp"    // Access role of Store on 'comp' api router.
	WRoleSMach = "store-machine" // Access role of Store on 'mach' api router.
	WRoleQKS   = "qk-partner"    // Access role of QKS   on 'part' api router.
	WRoleKL    = "kualian"       // Access role of Kualian on 'rb' api router.

	/* FIXME :
	 *
	 * Update the follow IsValidAdmin() and IsValidUser() methods
	 * when added same new role strings.
	 */
)

// RBAC role router key, one key maybe bind with multiple roles.
var _router_key_mapping = map[string]string{
	WRoleSuper: "admin", // super admin enable access admin apis.
	WRoleAdmin: "admin", // for admin role.
	WRoleUser:  "user",  // for user role.
	WRoleMComp: "comp",  // for ifsc  composer role: mall-comp.
	WRoleSComp: "comp",  // for store composer role: store-comp.
	WRoleSMach: "mach",  // for store machine  role: store-machine.
	WRoleQKS:   "part",  // for qks   partner  role: qks-partner.
	WRoleKL:    "rb",    // for kualian user   role: rb.
}

// RBAC role key mapping.
var _router_role_mapping = getRoleKeys()

// Check the given key whether valid role key.
func IsRoleKey(key string) bool {
	_, ok := _router_role_mapping[key]
	return ok
}

// Return role router key by given role, it maybe just return
// role string when not found from defined roles
func GetRoleKey(role string) string {
	if key, ok := _router_key_mapping[role]; ok {
		return key
	}
	return role
}

// Return all roles router keys .
func getRoleKeys() map[string]struct{} {
	mapping := make(map[string]struct{})
	for _, key := range _router_key_mapping {
		if _, ok := mapping[key]; !ok {
			mapping[key] = struct{}{}
		}
	}
	return mapping
}

// Check given role if super or admin role
func IsValidAdmin(role string) bool {
	switch role {
	case WRoleSuper, WRoleAdmin:
		return true
	}
	return false
}

// Check given role if normal user, not admins
func IsValidUser(role string) bool {
	switch role {
	case WRoleUser, WRoleMComp, WRoleSComp, WRoleSMach, WRoleQKS, WRoleKL:
		return true
	}
	return false
}
