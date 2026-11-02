package app

import (
	"context"
	"fmt"
	"time"

	"github.com/sunsetsavorer/climbing-elo/internal/adapter/clock"
	"github.com/sunsetsavorer/climbing-elo/internal/adapter/idgenerator"
	"github.com/sunsetsavorer/climbing-elo/internal/adapter/postgres"
	"github.com/sunsetsavorer/climbing-elo/internal/adapter/postgres/repository"
	"github.com/sunsetsavorer/climbing-elo/internal/config"
	"github.com/sunsetsavorer/climbing-elo/internal/delivery/telegram"
	"github.com/sunsetsavorer/climbing-elo/internal/usecase"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

const pollerTimeout = 10 * time.Second

type Bot struct {
	db     *gorm.DB
	server *telegram.Server
}

func NewBot(cfg *config.Config) (*Bot, error) {
	db, err := postgres.NewDB(cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	server, err := newServer(cfg, db)
	if err != nil {
		_ = closeDB(db)
		return nil, err
	}

	return &Bot{db: db, server: server}, nil
}

// Run blocks until ctx is cancelled, then stops the bot and releases resources.
func (b *Bot) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		b.server.Stop()
	}()

	b.server.Start()

	return closeDB(b.db)
}

func newServer(cfg *config.Config, db *gorm.DB) (*telegram.Server, error) {
	idGenerator := idgenerator.New()
	clk := clock.New()
	climberRepository := repository.NewClimberRepository(db)
	ascentRepository := repository.NewAscentRepository(db)
	txManager := repository.NewTxManager(db)

	getOrCreateClimber, err := usecase.NewGetOrCreateClimber(usecase.GetOrCreateClimberDeps{
		IDGenerator:       idGenerator,
		Clock:             clk,
		ClimberRepository: climberRepository,
	})
	if err != nil {
		return nil, fmt.Errorf("init get or create climber use case: %w", err)
	}

	recordAscent, err := usecase.NewRecordAscent(usecase.RecordAscentDeps{
		IDGenerator:       idGenerator,
		Clock:             clk,
		ClimberRepository: climberRepository,
		AscentRepository:  ascentRepository,
		TxManager:         txManager,
	})
	if err != nil {
		return nil, fmt.Errorf("init record ascent use case: %w", err)
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tele.LongPoller{Timeout: pollerTimeout},
	})
	if err != nil {
		return nil, fmt.Errorf("init telegram bot: %w", err)
	}

	server, err := telegram.NewServer(bot, getOrCreateClimber, recordAscent)
	if err != nil {
		return nil, fmt.Errorf("init telegram server: %w", err)
	}

	return server, nil
}

func closeDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}
