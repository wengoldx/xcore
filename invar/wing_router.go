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

	/* FIXME :
	 *
	 * Update the follow IsValidAdmin() and IsValidUser() methods
	 * when added same new role strings.
	 */
)

// RBAC role router keyword
const (
	WRGroupUser    = "user"
	WRGroupAdmin   = "admin"
	WRGroupComp    = "comp"
	WRGroupMachine = "mach"
	WRGroupPartner = "part"
)

// Return role router key by given role, it maybe just return
// role string when not found from defined roles
func GetRouterKey(role string) string {
	switch role {
	case WRoleSuper, WRoleAdmin:
		return WRGroupAdmin
	case WRoleUser:
		return WRGroupUser
	case WRoleMComp, WRoleSComp:
		return WRGroupComp
	case WRoleSMach:
		return WRGroupMachine
	case WRoleQKS:
		return WRGroupPartner
	}
	return role
}

// Check given role if super or admin role
func IsValidAdmin(role string) bool {
	return role != "" && (role == WRoleSuper || role == WRoleAdmin)
}

// Check given role if normal user, not admins
func IsValidUser(role string) bool {
	return role != "" && (role == WRoleUser || role == WRoleMComp ||
		role == WRoleSComp || role == WRoleSMach || role == WRoleQKS)
}
