package database

import (
	"context"
	"github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/jmoiron/sqlx"
)

type MailTemplateRepository interface {
	GetMailTemplate(ctx context.Context, mailTemplateId string) (*database.MailTemplate, error)
}

type mailTemplateRepository struct {
	db *sqlx.DB
}

func NewMailTemplateRepository(db *sqlx.DB) MailTemplateRepository {
	return &mailTemplateRepository{db}
}

func (r *mailTemplateRepository) GetMailTemplate(ctx context.Context, mailTemplateId string) (*database.MailTemplate, error) {
	var mailTemplate database.MailTemplate
	err := r.db.GetContext(ctx, &mailTemplate, `
		SELECT *
		FROM mail_templates
		WHERE mail_template_id = $1
	`, mailTemplateId)
	if err != nil {
		return nil, err
	}
	return &mailTemplate, nil
}
