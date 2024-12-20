package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gapidobri/prizer/internal/database"
	"github.com/gapidobri/prizer/internal/pkg/clients/addressvalidation"
	"github.com/gapidobri/prizer/internal/pkg/clients/mandrill"
	"github.com/gapidobri/prizer/internal/pkg/clients/sheets"
	er "github.com/gapidobri/prizer/internal/pkg/errors"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
	"github.com/gapidobri/prizer/internal/pkg/util"
	"github.com/google/uuid"
	"github.com/mattbaird/gochimp"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/mail"
	"strings"
	"time"
)

type GameService struct {
	gameRepository                database.GameRepository
	prizeRepository               database.PrizeRepository
	wonPrizeRepository            database.WonPrizeRepository
	userRepository                database.UserRepository
	drawMethodRepository          database.DrawMethodRepository
	participationMethodRepository database.ParticipationMethodRepository
	participationRepository       database.ParticipationRepository
	mailTemplateRepository        database.MailTemplateRepository
	addressValidationClient       addressvalidation.Client
	mandrillClient                mandrill.Client
	sheetsClient                  sheets.Client
}

func NewGameService(
	gameRepository database.GameRepository,
	prizeRepository database.PrizeRepository,
	wonPrizeRepository database.WonPrizeRepository,
	userRepository database.UserRepository,
	drawMethodRepository database.DrawMethodRepository,
	participationMethodRepository database.ParticipationMethodRepository,
	participationRepository database.ParticipationRepository,
	mailTemplateRepository database.MailTemplateRepository,
	addressValidationClient addressvalidation.Client,
	mandrillClient mandrill.Client,
	sheetsClient sheets.Client,
) *GameService {
	return &GameService{
		gameRepository:                gameRepository,
		prizeRepository:               prizeRepository,
		wonPrizeRepository:            wonPrizeRepository,
		userRepository:                userRepository,
		drawMethodRepository:          drawMethodRepository,
		participationMethodRepository: participationMethodRepository,
		participationRepository:       participationRepository,
		mailTemplateRepository:        mailTemplateRepository,
		addressValidationClient:       addressValidationClient,
		mandrillClient:                mandrillClient,
		sheetsClient:                  sheetsClient,
	}
}

func (s *GameService) GetGames(ctx context.Context) (api.GetGamesResponse, error) {
	games, err := s.gameRepository.GetGames(ctx)
	if err != nil {
		return nil, err
	}

	apiGames := lo.Map(games, func(game dbModels.Game, index int) api.Game {
		return api.GameFromDB(game)
	})

	return apiGames, nil
}

func (s *GameService) GetGame(ctx context.Context, gameId string) (api.GetGameResponse, error) {
	err := uuid.Validate(gameId)
	if err != nil {
		return nil, er.InvalidUuid
	}

	game, err := s.gameRepository.GetGame(ctx, gameId)
	if err != nil {
		return nil, err
	}

	apiGame := api.GameFromDB(*game)

	return &apiGame, nil
}

func (s *GameService) Participate(ctx context.Context, participationMethodId string, roll api.ParticipationRequest) (*api.ParticipationResponse, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"participationMethodId": participationMethodId,
		"fields":                roll.Fields,
	})

	logger.Info("New participation")

	participationMethod, err := s.participationMethodRepository.GetParticipationMethod(ctx, participationMethodId)
	if err != nil {
		return nil, err
	}

	game, err := s.gameRepository.GetGame(ctx, participationMethod.GameId)
	if err != nil {
		return nil, err
	}

	// Sanitize and validate incoming data
	additionalUserFields := dbModels.JsonMap{}
	var (
		userFields   dbModels.UserFields
		uniqueFields dbModels.UserFields
	)
	for key, field := range participationMethod.Fields.User {
		value, exists := roll.Fields[key]
		if !exists {
			if field.Unique || field.Required {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is required", key))
			}
			continue
		}

		switch key {
		case "email":
			str, ok := value.(string)
			if !ok {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
			}
			email, err := mail.ParseAddress(str)
			if err != nil {
				return nil, er.InvalidEmail
			}

			userFields.Email = &email.Address
			if field.Unique {
				uniqueFields.Email = &email.Address
			}
			continue

		case "address":
			str, ok := value.(string)
			if !ok {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
			}

			address, err := s.addressValidationClient.NormalizeAddress(ctx, str)
			if err != nil {
				return nil, err
			}

			userFields.Address = &address
			if field.Unique {
				uniqueFields.Address = &address
			}
			continue

		case "phone":
			str, ok := value.(string)
			if !ok {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
			}

			str = strings.TrimSpace(str)

			userFields.Phone = &str
			if field.Unique {
				uniqueFields.Phone = &str
			}
			continue
		}

		value, err = s.validateField(field, key, value)
		if err != nil {
			return nil, err
		}
		additionalUserFields[key] = value
	}

	participationFields := dbModels.JsonMap{}
	uniqueParticipationFields := dbModels.JsonMap{}
	for key, field := range participationMethod.Fields.Participation {
		value, exists := roll.Fields[key]
		if !exists {
			if field.Unique || field.Required {
				return nil, er.BadRequest.With(fmt.Sprintf("Field %s is required", key))
			}
			continue
		}

		value, err = s.validateField(field, key, value)
		if err != nil {
			return nil, err
		}
		participationFields[key] = value
		if field.Unique {
			uniqueParticipationFields[key] = value
		}
	}

	user, err := s.getCreateUser(ctx, game, userFields, uniqueFields, additionalUserFields)
	if err != nil {
		logger.WithError(err).Error("Failed to get/create user")
		return nil, err
	}

	logger = logger.WithField("userId", user.Id)

	// Check if user can participate
	switch participationMethod.ParticipationLimit {
	case enums.ParticipationLimitNone:
		break
	case enums.ParticipationLimitDaily:
		participations, err := s.participationRepository.GetParticipations(ctx, dbModels.GetParticipationsFilter{
			UserId:                &user.Id,
			ParticipationMethodId: &participationMethod.Id,
			From:                  lo.ToPtr(util.StripTime(time.Now())),
			To:                    lo.ToPtr(time.Now()),
		})
		if err != nil {
			logger.WithError(err).Error("Failed to get participations")
			return nil, err
		}
		if len(participations) != 0 {
			return nil, er.AlreadyParticipated
		}
	}

	// Check participation unique fields
	if len(uniqueParticipationFields) > 0 {
		participations, err := s.participationRepository.GetParticipations(ctx, dbModels.GetParticipationsFilter{
			UserId:                &user.Id,
			ParticipationMethodId: &participationMethod.Id,
			Fields:                &uniqueParticipationFields,
		})
		if err != nil {
			logger.WithError(err).Error("Failed to get participations")
			return nil, err
		}
		if len(participations) != 0 {
			return nil, er.ParticipationDataExists
		}
	}

	wonPrizes, err := s.drawPrizes(ctx, user, game, participationMethodId)
	if err != nil {
		logger.WithError(err).Error("Failed to draw prizes")
		return nil, err
	}

	participation, err := s.participationRepository.CreateParticipation(ctx, dbModels.CreateParticipation{
		ParticipationMethodId: participationMethodId,
		UserId:                user.Id,
		Fields:                participationFields,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to create participation")
		return nil, err
	}

	err = s.sendParticipationEmail(ctx, participationMethod, user, additionalUserFields)
	if err != nil {
		logger.WithError(err).Error("Failed to send participation email")
		return nil, err
	}

	for _, prize := range wonPrizes {
		err = s.wonPrizeRepository.CreateWonPrize(ctx, dbModels.CreateWonPrize{
			ParticipationId: participation.Id,
			PrizeId:         prize.Id,
		})
		if err != nil {
			logger.WithError(err).Error("Failed to create won prize")
			return nil, err
		}

		// Send win email
		if participationMethod.WinMailTemplateId != nil && user.Email != nil {
			template, err := s.mailTemplateRepository.GetMailTemplate(ctx, *participationMethod.WinMailTemplateId)
			if err != nil {
				logger.WithError(err).Error("Failed to get win mail template")
				return nil, err
			}

			variables := []gochimp.Var{
				{Name: "PRIZE_NAME", Content: prize.Name},
				{Name: "PRIZE_DESCRIPTION", Content: prize.Description},
				{Name: "PRIZE_IMAGE_URL", Content: prize.ImageUrl},
			}

			for key, field := range participationMethod.Fields.User {
				if field.MailVariable == nil {
					continue
				}
				value, exists := additionalUserFields[key]
				if !exists {
					continue
				}
				variables = append(variables, gochimp.Var{
					Name:    *field.MailVariable,
					Content: value,
				})
			}

			err = s.mandrillClient.SendTemplate(
				template.Name,
				gochimp.Message{
					FromEmail:       template.FromEmail,
					FromName:        template.FromName,
					Subject:         template.Subject,
					GlobalMergeVars: variables,
					To: []gochimp.Recipient{
						{Email: *user.Email},
					},
				},
			)
			if err != nil {
				log.WithError(err).Error("Failed to send win email")
			}
		}
	}

	// Send lose email
	if len(wonPrizes) == 0 {
		if participationMethod.LoseMailTemplateId != nil && user.Email != nil {
			template, err := s.mailTemplateRepository.GetMailTemplate(ctx, *participationMethod.LoseMailTemplateId)
			if err != nil {
				logger.WithError(err).Error("Failed to get lose mail template")
				return nil, err
			}

			var variables []gochimp.Var
			for key, field := range participationMethod.Fields.User {
				if field.MailVariable == nil {
					continue
				}
				value, exists := additionalUserFields[key]
				if !exists {
					continue
				}
				variables = append(variables, gochimp.Var{
					Name:    *field.MailVariable,
					Content: value,
				})
			}

			err = s.mandrillClient.SendTemplate(
				template.Name,
				gochimp.Message{
					FromEmail:       template.FromEmail,
					FromName:        template.FromName,
					Subject:         template.Subject,
					GlobalMergeVars: variables,
					To: []gochimp.Recipient{
						{Email: *user.Email},
					},
				},
			)
			if err != nil {
				log.WithError(err).Error("Failed to send lose email")
				return nil, err
			}
		}

		// Append participation without prize to google sheet
		err = s.appendRowToGoogleSheets(*game, *participationMethod, *participation, *user, nil)
		if err != nil {
			logger.WithError(err).Error("Failed to append row to google sheets")
		}
	}

	for _, prize := range wonPrizes {
		err = s.appendRowToGoogleSheets(*game, *participationMethod, *participation, *user, &prize)
		if err != nil {
			logger.WithError(err).Error("Failed to append row to google sheets")
		}
	}

	publicPrizes := lo.Map(wonPrizes, func(prize dbModels.Prize, _ int) api.PublicPrize {
		return api.PublicPrizeFromDB(prize)
	})

	return &api.ParticipationResponse{
		Prizes: publicPrizes,
	}, nil
}

func (s *GameService) validateField(field dbModels.Field, key string, value any) (any, error) {
	switch field.Type {
	case enums.FieldTypeBool:
		val, ok := value.(bool)
		if !ok {
			return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a bool", key))
		}
		return val, nil

	case enums.FieldTypeString:
		val, ok := value.(string)
		if !ok {
			return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a string", key))
		}
		return val, nil

	default:
		return nil, er.BadRequest.With(fmt.Sprintf("Field %s is not a valid field", key))
	}
}

func (s *GameService) getCreateUser(
	ctx context.Context,
	game *dbModels.Game,
	userFields dbModels.UserFields,
	uniqueFields dbModels.UserFields,
	additionalUserFields dbModels.JsonMap,
) (*dbModels.User, error) {
	user, err := s.userRepository.GetUserFromFields(ctx, game.Id, dbModels.UserFields{
		Email: uniqueFields.Email,
	})
	switch {
	case err == nil:
		break
	case errors.Is(err, er.UserNotFound):
		user, err = s.userRepository.CreateUser(ctx, dbModels.CreateUser{
			GameId:           game.Id,
			UserFields:       userFields,
			AdditionalFields: additionalUserFields,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create user")
		}
	default:
		return nil, err
	}

	return user, nil
}

func (s *GameService) drawPrizes(ctx context.Context, user *dbModels.User, game *dbModels.Game, participationMethodId string) ([]dbModels.Prize, error) {
	drawMethods, err := s.drawMethodRepository.GetDrawMethods(ctx, dbModels.GetDrawMethodsFilter{
		GameId:                &game.Id,
		ParticipationMethodId: &participationMethodId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get draw methods")
	}

	var wonPrizes []dbModels.Prize

	for _, drawMethod := range drawMethods {
		var prizes []dbModels.Prize
		prizes, err = s.prizeRepository.GetPrizes(ctx, dbModels.GetPrizesFilter{
			GameId:        &game.Id,
			DrawMethodId:  &drawMethod.Id,
			UserId:        &user.Id,
			AvailableOnly: true,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get prizes")
		}

		if len(prizes) == 0 {
			continue
		}

		switch drawMethod.Method {
		case enums.DrawMethodFirstN:
			wonPrizes = append(wonPrizes, prizes[0])

		case enums.DrawMethodChance:
			var data dbModels.DrawMethodChanceData
			err = json.Unmarshal([]byte(drawMethod.Data), &data)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal draw method chance data")
			}

			if data.Chance <= rand.Float64() {
				continue
			}

			prize := prizes[rand.Intn(len(prizes))]

			wonPrizes = append(wonPrizes, prize)
		}
	}

	return wonPrizes, nil
}

func (s *GameService) sendParticipationEmail(ctx context.Context, participationMethod *dbModels.ParticipationMethod, user *dbModels.User, fields dbModels.JsonMap) error {
	if participationMethod.ParticipationMailTemplateId == nil || user.Email == nil {
		return nil
	}

	template, err := s.mailTemplateRepository.GetMailTemplate(ctx, *participationMethod.ParticipationMailTemplateId)
	if err != nil {
		return errors.Wrap(err, "failed to get participation mail template")
	}

	var variables []gochimp.Var
	for key, field := range participationMethod.Fields.User {
		if field.MailVariable == nil {
			continue
		}
		value, exists := fields[key]
		if !exists {
			continue
		}
		variables = append(variables, gochimp.Var{
			Name:    *field.MailVariable,
			Content: value,
		})
	}

	err = s.mandrillClient.SendTemplate(
		template.Name,
		gochimp.Message{
			FromEmail:       template.FromEmail,
			FromName:        template.FromName,
			Subject:         template.Subject,
			GlobalMergeVars: variables,
			To: []gochimp.Recipient{
				{Email: *user.Email},
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to send lose email")
	}

	return nil
}

func (s *GameService) appendRowToGoogleSheets(
	game dbModels.Game,
	participationMethod dbModels.ParticipationMethod,
	participation dbModels.Participation,
	user dbModels.User,
	prize *dbModels.Prize,
) error {
	if game.GoogleSheetId == nil || game.GoogleSheetTabName == nil {
		return nil
	}

	fields := []any{
		participation.Id,
		participation.CreatedAt.Format("02. 01. 2006 15:04"),
		user.Email, user.Address, user.Phone,
	}
	for _, value := range user.AdditionalFields {
		fields = append(fields, value)
	}
	fields = append(fields, participationMethod.Name)
	var participationStrings []string
	for _, value := range participation.Fields {
		participationStrings = append(participationStrings, fmt.Sprintf("%v", value))
	}
	fields = append(fields, strings.Join(participationStrings, ", "))
	if prize != nil {
		fields = append(fields, prize.Name)
	} else {
		fields = append(fields, "/")
	}

	return s.sheetsClient.AppendRow(*game.GoogleSheetId, *game.GoogleSheetTabName, fields)
}
