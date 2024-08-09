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
  column "win_percentage" {
    type = float
    null = false
  }
  column "unique_collaborator_data" {
    type = boolean
  }
  primary_key {
    columns = [
      column.game_id
    ]
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
  column "collaborator_id" {
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
  foreign_key "collaborator_fk" {
    columns = [column.collaborator_id]
    ref_columns = [table.collaborator.column.collaborator_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}

table "collaborator" {
  schema = schema.public
  column "collaborator_id" {
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "email" {
    type = varchar
    null = false
  }
  column "address" {
    type = varchar
    null = true
  }
  column "game_id" {
    type = uuid
    null = false
  }
  column "last_roll_time" {
    type = timestamp
    null = true
  }
  unique "email_game_id" {
    columns = [column.email, column.game_id]
  }
  primary_key {
    columns = [
      column.collaborator_id
    ]
  }
  foreign_key "game_fk" {
    columns = [column.game_id]
    ref_columns = [table.game.column.game_id]
    on_delete = CASCADE
    on_update = CASCADE
  }
}