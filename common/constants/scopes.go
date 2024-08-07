package constants

type SCOPE string

type SCOPE_INFO struct {
	Scope       SCOPE  `json:"scope"`
	Description string `json:"description"`
}

type API_SCOPE struct {
	API    string       `json:"api"`
	SCOPES []SCOPE_INFO `json:"scopes"`
}

const (
	// API AAA scopes
	SCOPE_BSS_USER                SCOPE = "bss:user"                // View all users of your tenant
	SCOPE_BSS_USER_EXTENSION      SCOPE = "bss:user:extension"      // Allow create/update/delete user extensions
	SCOPE_BSS_TENANT              SCOPE = "bss:tenant"              // View all tenants
	SCOPE_BSS_BUSINESS_UNIT       SCOPE = "bss:business_unit"       // View all business units of your tenant
	SCOPE_BSS_GROUP               SCOPE = "bss:group"               // View all groups of your tenant
	SCOPE_BSS_ROLE                SCOPE = "bss:role"                // View all roles
	SCOPE_BSS_PERMISSION          SCOPE = "bss:permission"          // View all permissions
	SCOPE_BSS_USER_WRITE          SCOPE = "bss:user.write"          // Allow create/update/delete users
	SCOPE_BSS_TENANT_WRITE        SCOPE = "bss:tenant.write"        // Allow create/update/delete tenants
	SCOPE_BSS_BUSINESS_UNIT_WRITE SCOPE = "bss:business_unit.write" // Allow create/update/delete business units
	SCOPE_BSS_GROUP_WRITE         SCOPE = "bss:group.write"         // Allow create/update/delete groups
	SCOPE_BSS_ROLE_WRITE          SCOPE = "bss:role.write"          // Allow create/update/delete roles
	SCOPE_BSS_PERMISSION_WRITE    SCOPE = "bss:permission.write"    // Allow create/update/delete permissions

	// External
	SCOPE_EXT_PLUGIN              SCOPE = "ext:plugin"              // View all plugins
	SCOPE_EXT_PLUGIN_WRITE        SCOPE = "ext:plugin.write"        // Allow create/update/delete plugins
	SCOPE_EXT_TENANT_PLUGIN       SCOPE = "ext:tenant_plugin"       // View all tenant plugins
	SCOPE_EXT_TENANT_PLUGIN_WRITE SCOPE = "ext:tenant_plugin.write" // Allow create/update/delete tenant plugins
)

var API_SCOPES = []API_SCOPE{
	{
		API: "API Authentication/Authorization",
		SCOPES: []SCOPE_INFO{
			{Scope: SCOPE_BSS_USER, Description: "View all users of your tenant"},
			{Scope: SCOPE_BSS_USER_EXTENSION, Description: "Allow create/update/delete user extensions"},
			{Scope: SCOPE_BSS_TENANT, Description: "View all tenants"},
			{Scope: SCOPE_BSS_BUSINESS_UNIT, Description: "View all business units of your tenant"},
			{Scope: SCOPE_BSS_GROUP, Description: "View all groups"},
			{Scope: SCOPE_BSS_ROLE, Description: "View all roles"},
			{Scope: SCOPE_BSS_PERMISSION, Description: "View all permissions"},
			{Scope: SCOPE_BSS_USER_WRITE, Description: "Allow to create and update users"},
			{Scope: SCOPE_BSS_TENANT_WRITE, Description: "Allow to create and update tenants"},
			{Scope: SCOPE_BSS_BUSINESS_UNIT_WRITE, Description: "Allow to create and update"},
			{Scope: SCOPE_BSS_GROUP_WRITE, Description: "Allow to create and update groups"},
			{Scope: SCOPE_BSS_ROLE_WRITE, Description: "Allow to create and update roles"},
			{Scope: SCOPE_BSS_PERMISSION_WRITE, Description: "Allow to create and update permissions"},
		},
	},
	{
		API: "API External",
		SCOPES: []SCOPE_INFO{
			{Scope: SCOPE_EXT_PLUGIN, Description: "View plugins"},
			{Scope: SCOPE_EXT_PLUGIN_WRITE, Description: "Allow to create and update plugin"},
			{Scope: SCOPE_EXT_TENANT_PLUGIN, Description: "View tenant plugins"},
			{Scope: SCOPE_EXT_TENANT_PLUGIN_WRITE, Description: "Allow to create and update tenant plugin"},
		},
	},
}
