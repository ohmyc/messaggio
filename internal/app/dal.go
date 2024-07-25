package app

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDal() (*Dal, error) {
	connect, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s host=%s dbname=messaggio sslmode=disable",
			PostgresUser,
			PostgresPassword,
			PostgresHost,
		),
	)
	if err != nil {
		return nil, err
	}
	return &Dal{connect}, err
}

type Message struct {
	ID            string  `db:"id" json:"id"`
	Text          string  `db:"text" json:"text"`
	ProcessedText *string `db:"processed_text" json:"processed_text"`
}

type Dal struct {
	inner *sqlx.DB
}

func (d *Dal) EnsureCreated() error {
	_, err := d.inner.Query(`
CREATE TABLE IF NOT EXISTS messages 
(
    id TEXT NOT NULL PRIMARY KEY,
    text TEXT NOT NULL,
    processed_text TEXT
)`)
	return err
}

func (d *Dal) InsertNewMessage(id, text string) error {
	_, err := d.inner.Query(`
INSERT INTO messages (id, text) VALUES ($1, $2)
`, id, text)
	return err
}

func (d *Dal) UpdateProcessedMessage(id string, processedText string) error {
	_, err := d.inner.Query(`
UPDATE messages SET processed_text=$1 WHERE id=$2
`, processedText, id)
	return err
}

var ErrNotProcessed = errors.New("not processed")

func (d *Dal) GetProcessedMessage(id string) (*Message, error) {
	msg := new(Message)
	if err := d.inner.Get(`SELECT * FROM messages where id=$1`, id); err != nil {
		return nil, err
	}
	if msg.ProcessedText == nil {
		return nil, ErrNotProcessed
	}
	return msg, nil
}

func (d *Dal) getInt(query string) (int, error) {
	ret := 0
	if err := d.inner.Select(&ret, query); err != nil {
		return 0, err
	}
	return ret, nil
}

type Stats struct {
	ProcessedCount   int `json:"total_processed"`
	UnprocessedCount int `json:"total_unprocessed"`
}

func (d *Dal) GetStats() (*Stats, error) {
	proc, err := d.getInt("select count(*) from messages where processed_text != NULL")
	if err != nil {
		return nil, err
	}
	unproc, err := d.getInt("select count(*) from messages where processed_text = NULL")
	if err != nil {
		return nil, err
	}
	return &Stats{proc, unproc}, nil
}

func (d *Dal) GetAll() ([]*Message, error) {
	result := make([]*Message, 0)
	if err := d.inner.Select(&result, "select * from messages"); err != nil {
		return nil, err
	}
	return result, nil
}
