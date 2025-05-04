schema "public" {}

enum "split_type" {
  schema = schema.public
  values = ["equal", "proportional", "transfer"]
}

table "expenses" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "amount_cents" {
    type = bigint
    null = false
  }
  column "refund_amount_cents" {
    type = bigint
    null = true
  }
  column "description" {
    type = varchar(255)
    null = false
  }
  column "group_id" {
    type = bigint
    null = false
  }
  column "category_id" {
    type = bigint
    null = false
  }
  column "split_ratio" {
    type = jsonb
    null = false
  }
  column "split_type" {
    type = enum.split_type
    null = false
  }
  column "payer_id" {
    type = bigint
    null = false
  }
  column "receiver_id" {
    type = bigint
    null = false
  }
  column "document_search" {
    type = tsvector
    as {
      expr = "setweight(to_tsvector('portuguese', coalesce(name, '')), 'A') || setweight(to_tsvector('portuguese', coalesce(description, '')), 'B')"
      type = STORED
    }
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id, column.version]
  }

  foreign_key "group_id_fk" {
    columns     = [column.group_id]
    ref_columns = [table.groups.column.id]
  }

  foreign_key "payer_id_fk" {
    columns     = [column.payer_id]
    ref_columns = [table.users.column.id]
  }

  foreign_key "receiver_id_fk" {
    columns     = [column.receiver_id]
    ref_columns = [table.users.column.id]
  }

  foreign_key "category_id_fk" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
  }

  index "document_search_idx" {
    type    = GIN
    columns = [column.document_search]
  }
}


table "groups" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }
}

table "categories" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "icon" {
    type = varchar(255)
    null = false
  }
  column "category_group_id" {
    type = bigint
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "category_group_id_fk" {
    columns     = [column.category_group_id]
    ref_columns = [table.category_groups.column.id]
  }
}

table "category_groups" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "icon" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }
}

table "users" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "profile_picture" {
    type = varchar(255)
    null = true
  }
  column "group_id" {
    type = bigint
    null = true
  }
  column "flags" {
    type    = sql("text[]")
    default = sql("array[]::text[]")
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  index "email_unique_idx" {
    columns = [column.email]
    unique  = true
  }
}

table "authentications" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "password" {
    type = varchar(255)
    null = true
  }
  column "provider_id" {
    type = varchar(255)
    null = true
  }
  column "type" {
    type = enum.authentication_type
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "auth_email_fk" {
    columns     = [column.email]
    ref_columns = [table.users.column.email]
  }

  index "email_type_unique_idx" {
    columns = [column.email, column.type]
    unique  = true
  }
}

enum "authentication_type" {
  schema = schema.public
  values = ["credentials", "google"]
}

table "incomes" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "user_id" {
    type = bigint
    null = false
  }
  column "amount_cents" {
    type = bigint
    null = false
  }
  column "type" {
    type = enum.income_type
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "income_user_id_fk" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
  }
}

enum "income_type" {
  schema = schema.public
  values = ["salary", "benefit", "vacation", "thirteenth_salary", "other"]
}

table "group_invites" {
  schema = schema.public
  column "id" {
    type = bigserial
    null = false
  }
  column "group_id" {
    type = bigint
    null = false
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "token" {
    type    = uuid
    default = sql("gen_random_uuid()")
    null    = false
  }
  column "status" {
    type = enum.group_invites_status
    null = false
  }
  column "expires_at" {
    type = timestamptz
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "deleted_at" {
    type = timestamptz
    null = true
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  index "group_invite_email_idx" {
    columns = [column.email]
  }

  index "group_invite_token_idx" {
    columns = [column.token]
    unique  = true
  }
}

enum "group_invites_status" {
  schema = schema.public
  values = ["pending", "sent", "accepted"]
}

table "scheduled_expenses" {
  schema = schema.public

  column "id" {
    type = bigserial
    null = false
  }
  column "name" {
    type = text
    null = false
  }
  column "amount_cents" {
    type = bigint
    null = false
  }
  column "description" {
    type = text
    null = false
  }
  column "group_id" {
    type = bigint
    null = false
  }
  column "category_id" {
    type = bigint
    null = false
  }
  column "split_type" {
    type = enum.split_type
    null = false
  }
  column "payer_id" {
    type = bigint
    null = false
  }
  column "receiver_id" {
    type = bigint
    null = false
  }
  column "frequency_in_days" {
    type = int
    null = false
  }
  column "last_generated_at" {
    type = date
    null = true
  }
  column "is_active" {
    type = boolean
    null = false
  }
  column "created_at" {
    type = timestamptz
    null = false
  }
  column "updated_at" {
    type = timestamptz
    null = false
  }
  column "version" {
    type = int
    null = false
  }

  primary_key {
    columns = [column.id]
  }
}
