package domain

type (
	Token struct {
		Id        int64  `db:"id"`
		Number    int64  `db:"number"`
		Purpose   int    `db:"purpose"`
		Secret    string `db:"secret"`
		ExpiresAt int64  `db:"expires_at"`
	}
)
