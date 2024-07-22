resource "snowflake_account_role" "test" {
  name = var.account_role_name
}

resource "snowflake_account_role" "other_role" {
  name = var.other_account_role_name
}

resource "snowflake_database" "test" {
  name = var.database_name
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_account_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}