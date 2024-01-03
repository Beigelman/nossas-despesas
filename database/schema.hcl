schema "public" {}

table "expenses" {
  schema = schema.public
  column "id" {
    type = serial
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
	column "description" {
    type = varchar(255)
    null = false
  }
	column "group_id" {
    type = int
    null = false
  }
	column "category_id" {
    type = int
    null = false
  }
	column "split_ratio" {
    type = jsonb
    null = false
  }
	column "payer_id" {
    type = int
    null = false
  }
	column "receiver_id" {
    type = int
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

  primary_key  {
    columns = [column.id, column.version]
  }

  foreign_key "group_id_fk" {
    columns = [column.group_id ]
    ref_columns = [table.groups.column.id]
  }

  foreign_key "payer_id_fk" {
    columns = [column.payer_id]
    ref_columns = [table.users.column.id]
  }

  foreign_key "receiver_id_fk" {
    columns = [column.receiver_id]
    ref_columns = [table.users.column.id]
  }
  
  foreign_key "category_id_fk" {
    columns = [column.category_id]
    ref_columns = [table.categories.column.id]
  }
}

table "groups" {
  schema = schema.public
  column "id" {
    type = serial
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

  primary_key  {
    columns = [column.id]
  }
}

table "categories" {
  schema = schema.public
  column "id" {
    type = serial
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
    type = int
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

  primary_key  {
    columns = [column.id]
  }

  foreign_key "category_group_id_fk" {
    columns = [column.category_group_id]
    ref_columns = [table.category_groups.column.id]
  }
}

table "category_groups" {
  schema = schema.public
  column "id" {
    type = serial
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

  primary_key  {
    columns = [column.id]
  }
}

table "users" {
  schema = schema.public
  column "id" {
    type = serial
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
    type = int
    null = true
  }
  column "authentication_id" {
     type = varchar(255)
     null = true
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

  primary_key  {
    columns = [column.id]
  }

  index "authentication_id_idx" {
    columns = [column.authentication_id]
    unique = true
  }
}

table "auth" {
  schema = schema.public
  column "id" {
    type =
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
    type = int
    null = true
  }
  column "" {
     type = varchar(255)
     null = true
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

  primary_key  {
    columns = [column.id]
  }

  index "authentication_id_idx" {
    columns = [column.authentication_id]
    unique = true
  }
}
