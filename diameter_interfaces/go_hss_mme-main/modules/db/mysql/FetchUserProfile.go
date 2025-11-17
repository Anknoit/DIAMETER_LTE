package mysql

import (
	"database/sql"
	"fmt"
)

// UserProfile represents a user profile in the database
type UserProfile struct {
	IMSI       string
	RAND       []byte
	MSISDN     sql.NullString
	IMEI       sql.NullString
	MSPSStatus sql.NullString
}

// FetchRDSUserProfile retrieves the full user profile from the `users` table in the database
func FetchRDSUserProfile(imsi string) (*UserProfile, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection pool is not initialized")
	}

	query := `
	SELECT rand, msisdn, imei, ms_ps_status
	FROM users
	WHERE imsi = ?`

	var user UserProfile

	err := DB.QueryRow(query, imsi).Scan(
		&user.RAND, &user.MSISDN, &user.IMEI, &user.MSPSStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
