schema "public" {}

table "game" {
  schema = schema.public
  column "game_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "name" {
    type = varchar
    null = false
  }
  primary_key {
    columns = [
      column.game_id
    ]
  }
}

enum "participation_limit" {
  schema = schema.public
  values = ["daily"]
}

table "participation_method" {
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
    null = true
  }
  column "fields" {
    type = json
  }
  primary_key {
    columns = [
      column.participation_method_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.game.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "participation_method_draw_method" {
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
    ref_columns = [table.participation_method.column.participation_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "draw_method_fk" {
    columns = [column.draw_method_id]
    ref_columns = [table.draw_method.column.draw_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

enum "draw_method_enum" {
  schema = schema.public
  values = ["first_n", "chance"]
}

table "draw_method" {
  schema = schema.public
  column "draw_method_id" {
    type = uuid
    default = sql("gen_random_uuid()")
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
}

table "draw_method_prize" {
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
    ref_columns = [table.draw_method.column.draw_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "prize_fk" {
    columns = [column.prize_id]
    ref_columns = [table.prize.column.prize_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "prize" {
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
    ref_columns = [table.game.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "won_prize" {
  schema = schema.public
  column "won_prize_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "prize_id" {
    type = uuid
  }
  column "user_id" {
    type = uuid
  }
  column "created_at" {
    type = timestamp
    default = sql("now()")
  }
  primary_key {
    columns = [
      column.won_prize_id
    ]
  }
  foreign_key "prize_fk" {
    columns = [column.prize_id]
    ref_columns = [table.prize.column.prize_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "user_fk" {
    columns = [column.user_id]
    ref_columns = [table.user.column.user_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "user" {
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
    ref_columns = [table.game.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "participation" {
  schema = schema.public
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
  primary_key {
    columns = [
      column.participation_method_id,
      column.user_id,
      column.created_at,
    ]
  }
  foreign_key "participation_method_fk" {
    columns = [column.participation_method_id]
    ref_columns = [table.participation_method.column.participation_method_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
  foreign_key "user_fk" {
    columns = [column.user_id]
    ref_columns = [table.user.column.user_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}