---
page_title: "snowflake_grant_application_role Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# snowflake_grant_application_role (Resource)



## Example Usage

```terraform
locals {
  application_role_identifier = "\"my_appplication\".\"app_role_1\""
}

##################################
### grant application role to account role
##################################


resource "snowflake_account_role" "role" {
  name = "my_role"
}

resource "snowflake_grant_application_role" "g" {
  application_role_name    = local.application_role_identifier
  parent_account_role_name = snowflake_account_role.role.name
}

##################################
### grant application role to application
##################################

resource "snowflake_grant_application_role" "g" {
  application_role_name = local.application_role_identifier
  application_name      = "my_second_application"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_role_name` (String) Specifies the identifier for the application role to grant.

### Optional

- `application_name` (String) The fully qualified name of the application on which application role will be granted.
- `parent_account_role_name` (String) The fully qualified name of the account role on which application role will be granted.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# format is application_role_name (string) | object_type (ACCOUNT_ROLE|APPLICATION) | grantee_name (string)
terraform import "\"my_application\".\"app_role_1\"|ACCOUNT_ROLE|\"my_role\""
```