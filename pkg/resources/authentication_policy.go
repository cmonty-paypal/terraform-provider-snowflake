package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var authenticationPolicySchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database this authentication policy belongs to.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema this authentication policy belongs to.",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the authentication policy; must be unique for your account.",
	},
	"or_replace": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Whether to override a previous authentication policy with the same name.",
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"if_not_exists": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Prevent overwriting a previous authentication policy with the same name.",
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"authentication_methods": {
		Type:         schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:     true,
		Default: []string{"all"},
		Description:  "A list of authentication methods that are allowed during login. Allowed values: `all`, `saml`, `password`, `oauth`, `keypair`",
		// ValidateFunc: validation.MapValueMatch()
	},
	"mfa_authentication_methods": {
		Type:         schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:     true,
		Description:  "A list of authentication methods that enforce multi-factor authentication (MFA) during login. Allowed values: `all`, `saml`, `password`",
		// ValidateFunc: validation.IntBetween(8, 256),
	},
	"mfa_enrollment": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "optional",
		Description:  "Determines whether a user must enroll in multi-factor authentication. Allowed values: `optional`, `required`",
		// ValidateFunc: validation.IntBetween(0, 256),
	},
	"client_types": {
		Type:         schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:     true,
		Default: []string{"all"},
		Description:  "A list of clients that can authenticate with Snowflake. Allowed values: `all`, `snowflake_ui`, `drivers`, `snowsql`",
		// ValidateFunc: validation.IntBetween(0, 256),
	},
	"security_integrations": {
		Type:         schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:     true,
		Default: []string{"all"},
		Description:  "A list of security integrations the authentication policy is associated with.",
		// ValidateFunc: validation.IntBetween(0, 256),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Adds a comment or overwrites an existing comment for the authentication policy.",
	},
	"qualified_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The qualified name for the authentication policy.",
	},
}

func AuthenticationPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "An authentication policy specifies the requirements that must be met to authenticate to Snowflake.",
		Create:      CreatePasswordPolicy,
		Read:        ReadPasswordPolicy,
		Update:      UpdatePasswordPolicy,
		Delete:      DeletePasswordPolicy,

		Schema: authenticationPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
