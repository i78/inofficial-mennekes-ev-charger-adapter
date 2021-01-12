package persistor

import (
	"context"
	"fmt"
	"github.com/codecyclist/ev-charger-adapter/models"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

type Persistor struct {
	context  context.Context
	database *pgxpool.Pool
}

func NewPersistor(ctx context.Context) (p *Persistor) {
	var pool *pgxpool.Pool
	var err error

	pool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}

	return &Persistor{
		context:  ctx,
		database: pool,
	}
}

func (p *Persistor) Handle(readingName string, reading models.ChargerValueEnvelope) error {
	log.WithFields(log.Fields{
		"readingName": readingName,
		"value" : reading.Value,
	}).Debug("Handling new charger message")

	conn, err := p.database.Acquire(context.Background())

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Error while connecting to db")
		return err
	}
	defer conn.Release()

	sanitizer := regexp.MustCompile(`[^a-zA-Z]`)
	sanitizedReadingName := sanitizer.ReplaceAll([]byte(readingName), []byte(""))

	sql := fmt.Sprintf("INSERT INTO public.ev_readings " +
		"(time, %s) VALUES ($1, $2)", string(sanitizedReadingName))

	samplingTimeUtc := reading.SamplingTime.UTC()
	_, errInsert := conn.Exec(context.Background(), sql, samplingTimeUtc, reading.Value)

	if errInsert != nil {
		log.WithFields(log.Fields{
			"error": errInsert,
		}).Warn("Error while inserting data")
		return errInsert
	}

	return nil
}

func (r *Persistor) Run() {
	log.WithFields(log.Fields{
		"chargerIp": "na",
	}).Info("Starting Subscriber")

	for {
		select {

		case <-r.context.Done():
			log.Info("Terminating Subscriber")
			return
		}

	}
}
