package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

// CreateAuthenticationPolicy implements schema.CreateFunc.
func CreateAuthenticationPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	objectIdentifier := sdk.NewSchemaObjectIdentifier(database, schema, name)

	createOptions := &sdk.CreateAuthenticationPolicyRequest{
		name: objectIdentifier,
		OrReplace:                 sdk.Bool(d.Get("or_replace").(bool)),
		AuthenticationMethods: d.Get("authentication_methods").([]sdk.AuthenticationMethods),
		MfaAuthenticationMethods: d.Get("mfa_authentication_methods").([]sdk.MfaAuthenticationMethods),
		MfaEnrollment: sdk.String(d.Get("mfa_enrollment").(string)),
		ClientTypes: d.Get("client_types").([]sdk.ClientTypes),
		SecurityIntegrations: d.Get("security_integrations").([]sdk.SecurityIntegrationsOption),
	}

	if v, ok := d.GetOk("comment"); ok {
		createOptions.Comment = sdk.String(v.(string))
	}

	err := client.AuthenticationPolicies.Create(ctx, createOptions)
	if err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))
	return ReadAuthenticationPolicy(d, meta)
}

// ReadAuthenticationPolicy implements schema.ReadFunc.
func ReadAuthenticationPolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if err := d.Set("qualified_name", objectIdentifier.FullyQualifiedName()); err != nil {
		return err
	}

	authenticationPolicy, err := client.AuthenticationPolicies.ShowByID(ctx, objectIdentifier)
	if err != nil {
		return err
	}

	if err := d.Set("database", authenticationPolicy.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema", authenticationPolicy.SchemaName); err != nil {
		return err
	}
	if err := d.Set("name", authenticationPolicy.Name); err != nil {
		return err
	}
	if err := d.Set("comment", authenticationPolicy.Comment); err != nil {
		return err
	}
	// passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, objectIdentifier)
	// if err != nil {
	// 	return err
	// }

	// if err := setIntProperty(d, "min_length", passwordPolicyDetails.PasswordMinLength); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "max_length", passwordPolicyDetails.PasswordMaxLength); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "min_upper_case_chars", passwordPolicyDetails.PasswordMinUpperCaseChars); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "min_lower_case_chars", passwordPolicyDetails.PasswordMinLowerCaseChars); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "min_numeric_chars", passwordPolicyDetails.PasswordMinNumericChars); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "min_special_chars", passwordPolicyDetails.PasswordMinSpecialChars); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "min_age_days", passwordPolicyDetails.PasswordMinAgeDays); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "max_age_days", passwordPolicyDetails.PasswordMaxAgeDays); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "max_retries", passwordPolicyDetails.PasswordMaxRetries); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "lockout_time_mins", passwordPolicyDetails.PasswordLockoutTimeMins); err != nil {
	// 	return err
	// }
	// if err := setIntProperty(d, "history", passwordPolicyDetails.PasswordHistory); err != nil {
	// 	return err
	// }

	return nil
}

