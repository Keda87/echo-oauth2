package databases

import "github.com/jmoiron/sqlx"

type Manager struct {
	DB *sqlx.DB
}
