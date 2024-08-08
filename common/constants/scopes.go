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
	SCOPE_IBK_USER                SCOPE = "ibk:user"                // View all users of your tenant
	SCOPE_IBK_USER_EXTENSION      SCOPE = "ibk:user:extension"      // Allow create/update/delete user extensions
	SCOPE_IBK_TENANT              SCOPE = "ibk:tenant"              // View all tenants
	SCOPE_IBK_BUSINESS_UNIT       SCOPE = "ibk:business_unit"       // View all business units of your tenant
	SCOPE_IBK_ROLE                SCOPE = "ibk:role"                // View all roles
	SCOPE_IBK_PERMISSION          SCOPE = "ibk:permission"          // View all permissions
	SCOPE_IBK_USER_WRITE          SCOPE = "ibk:user.write"          // Allow create/update/delete users
	SCOPE_IBK_TENANT_WRITE        SCOPE = "ibk:tenant.write"        // Allow create/update/delete tenants
	SCOPE_IBK_BUSINESS_UNIT_WRITE SCOPE = "ibk:business_unit.write" // Allow create/update/delete business units
	SCOPE_IBK_ROLE_WRITE          SCOPE = "ibk:role.write"          // Allow create/update/delete roles
	SCOPE_IBK_PERMISSION_WRITE    SCOPE = "ibk:permission.write"    // Allow create/update/delete permissions
)

var API_SCOPES = []API_SCOPE{
	{
		API: "API Authentication/Authorization",
		SCOPES: []SCOPE_INFO{
			{Scope: SCOPE_IBK_USER, Description: "View all users of your tenant"},
			{Scope: SCOPE_IBK_TENANT, Description: "View all tenants"},
			{Scope: SCOPE_IBK_BUSINESS_UNIT, Description: "View all business units of your tenant"},
			{Scope: SCOPE_IBK_ROLE, Description: "View all roles"},
			{Scope: SCOPE_IBK_PERMISSION, Description: "View all permissions"},
			{Scope: SCOPE_IBK_USER_WRITE, Description: "Allow to create and update users"},
			{Scope: SCOPE_IBK_TENANT_WRITE, Description: "Allow to create and update tenants"},
			{Scope: SCOPE_IBK_BUSINESS_UNIT_WRITE, Description: "Allow to create and update"},
			{Scope: SCOPE_IBK_ROLE_WRITE, Description: "Allow to create and update roles"},
			{Scope: SCOPE_IBK_PERMISSION_WRITE, Description: "Allow to create and update permissions"},
		},
	},
}
