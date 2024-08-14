schema "public" {}

table "games" {
  schema = schema.public
  column "game_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "name" {
    type = varchar
    null = false
  }
  column "google_sheet_id" {
    type = varchar
    null = true
  }
  column "google_sheet_tab_name" {
    type = varchar
    null = true
  }
  primary_key {
    columns = [
      column.game_id
    ]
  }
}

enum "participation_limit" {
  schema = schema.public
  values = ["none", "daily"]
}

table "participation_methods" {
  schema = schema.public
  column "participation_method_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "game_id" {
    type = uuid
  }
  column "name" {
    type = varchar
  }
  column "limit" {
    type = enum.participation_limit
    default = "none"
  }
  column "fields" {
    type = json
  }
  column "win_mail_template_id" {
    type = uuid
    null = true
  }
  column "lose_mail_template_id" {
    type = uuid
    null = true
  }
  primary_key {
    columns = [
      column.participation_method_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.games.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "win_mail_template_fk" {
    columns = [column.win_mail_template_id]
    ref_columns = [table.mail_templates.column.mail_template_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "lose_mail_template_fk" {
    columns = [column.lose_mail_template_id]
    ref_columns = [table.mail_templates.column.mail_template_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "participation_methods_draw_methods" {
  schema = schema.public
  column "participation_method_id" {
    type = uuid
  }
  column "draw_method_id" {
    type = uuid
  }
  primary_key {
    columns = [
      column.participation_method_id,
      column.draw_method_id
    ]
  }
  foreign_key "participation_method_fk" {
    columns = [column.participation_method_id]
    ref_columns = [table.participation_methods.column.participation_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "draw_method_fk" {
    columns = [column.draw_method_id]
    ref_columns = [table.draw_methods.column.draw_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

enum "draw_method_enum" {
  schema = schema.public
  values = ["first_n", "chance"]
}

table "draw_methods" {
  schema = schema.public
  column "draw_method_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "game_id" {
    type = uuid
  }
  column "name" {
    type = varchar
    default = "Draw method"
  }
  column "method" {
    type = enum.draw_method_enum
  }
  column "data" {
    type = json
  }
  primary_key {
    columns = [
      column.draw_method_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.games.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "draw_methods_prizes" {
  schema = schema.public
  column "draw_method_id" {
    type = uuid
  }
  column "prize_id" {
    type = uuid
  }
  primary_key {
    columns = [
      column.draw_method_id,
      column.prize_id
    ]
  }
  foreign_key "draw_method_fk" {
    columns = [column.draw_method_id]
    ref_columns = [table.draw_methods.column.draw_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "prize_fk" {
    columns = [column.prize_id]
    ref_columns = [table.prizes.column.prize_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "prizes" {
  schema = schema.public
  column "prize_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "game_id" {
    type = uuid
  }
  column "name" {
    type = varchar
    null = false
  }
  column "description" {
    type = varchar
  }
  column "image_url" {
    type = varchar
    null = true
  }
  column "count" {
    type = integer
    null = false
  }
  primary_key {
    columns = [
      column.prize_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.games.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "won_prizes" {
  schema = schema.public
  column "prize_id" {
    type = uuid
  }
  column "participation_id" {
    type = uuid
  }
  primary_key {
    columns = [
      column.prize_id,
      column.participation_id
    ]
  }
  foreign_key "prize_fk" {
    columns = [column.prize_id]
    ref_columns = [table.prizes.column.prize_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "participation_fk" {
    columns = [column.participation_id]
    ref_columns = [table.participations.column.participation_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "users" {
  schema = schema.public
  column "user_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "email" {
    type = varchar
    null = true
  }
  column "address" {
    type = varchar
    null = true
  }
  column "phone" {
    type = varchar
    null = true
  }
  column "additional_fields" {
    type = json
  }
  column "game_id" {
    type = uuid
    null = false
  }
  primary_key {
    columns = [
      column.user_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.games.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "participations" {
  schema = schema.public
  column "participation_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "participation_method_id" {
    type = uuid
  }
  column "user_id" {
    type = uuid
  }
  column "created_at" {
    type = timestamp
    default = sql("now()")
  }
  column "fields" {
    type = json
  }
  unique "method_user_time" {
    columns = [
      column.participation_id,
      column.participation_method_id,
      column.user_id
    ]
  }
  primary_key {
    columns = [
      column.participation_id,
    ]
  }
  foreign_key "participation_method_fk" {
    columns = [column.participation_method_id]
    ref_columns = [table.participation_methods.column.participation_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "user_fk" {
    columns = [column.user_id]
    ref_columns = [table.users.column.user_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "mail_templates" {
  schema = schema.public
  column "mail_template_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "game_id" {
    type = uuid
  }
  column "name" {
    type = varchar
  }
  column "from_email" {
    type = varchar
  }
  column "from_name" {
    type = varchar
  }
  column "subject" {
    type = varchar
  }
  primary_key {
    columns = [
      column.mail_template_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.games.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}
