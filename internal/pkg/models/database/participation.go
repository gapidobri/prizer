package database

import "time"

type Participation struct {
	Id                    string    `db:"participation_id"`
	UserId                string    `db:"user_id"`
	ParticipationMethodId string    `db:"participation_method_id"`
	CreatedAt             time.Time `db:"created_at"`
	Fields                JsonMap
}

type CreateParticipation struct {
	UserId                string  `db:"user_id"`
	ParticipationMethodId string  `db:"participation_method_id"`
	Fields                JsonMap `db:"fields"`
}

type GetParticipationsFilter struct {
	UserId                *string
	ParticipationMethodId *string
	From                  *time.Time
	To                    *time.Time
	Fields                *JsonMap
}
