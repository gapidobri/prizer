package database

type MailTemplate struct {
	Id        string `db:"mail_template_id"`
	GameId    string `db:"game_id"`
	Name      string `db:"name"`
	FromEmail string `db:"from_email"`
	FromName  string `db:"from_name"`
	Subject   string `db:"subject"`
}
