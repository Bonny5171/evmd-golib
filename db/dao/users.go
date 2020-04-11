package dao

import (
	"time"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

func GetUser(conn *sqlx.DB, tid int, uid string) (u model.User, err error) {
	const query = `
		SELECT user_id, username, name, firstname, lastname, email, full_photo_url, access_token, refresh_token, instance_url 
		  FROM public."user"
		 WHERE tenant_id = $1
		   AND user_id = $2
		 LIMIT 1;`

	err = conn.QueryRowx(query, tid, uid).StructScan(&u)
	if err != nil {
		err = db.WrapError(err, "conn.QueryRowx()")
		return
	}

	return u, nil
}

func GetUsersToProcess(conn *sqlx.DB, tid int) (u model.Users, err error) {
	const query = `
		SELECT user_id, access_token, refresh_token, instance_url 
		  FROM public."user"
		 WHERE tenant_id = $1
		   AND is_active = TRUE
		   AND is_deleted = FALSE;`

	err = conn.Select(&u, query, tid)
	if err != nil {
		err = db.WrapError(err, "conn.Select()")
		return
	}

	return u, nil
}

func UpdateUserAccessToken(conn *sqlx.DB, tid int, userID, accessToken string) (err error) {
	t := time.Now()

	const query = `
		UPDATE public."user" 
		   SET access_token = $3,
		       updated_at = $4
		 WHERE tenant_id = $1 
		   AND user_id = $2;`

	if _, err = conn.Exec(query, tid, userID, accessToken, t); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func SaveUser(conn *sqlx.DB, tid int, user model.User) (err error) {
	const query = `
		INSERT INTO public."user" (tenant_id, user_id, username, name, firstname, lastname, email, full_photo_url, access_token, refresh_token, instance_url) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		    ON CONFLICT (tenant_id, user_id) DO UPDATE 
		   SET username       = EXCLUDED.username, 
		       name           = EXCLUDED.name, 
		       firstname      = EXCLUDED.firstname, 
		       lastname       = EXCLUDED.lastname, 
		       email          = EXCLUDED.email, 
		       full_photo_url = EXCLUDED.full_photo_url, 
		       access_token   = EXCLUDED.access_token, 
		       refresh_token  = EXCLUDED.refresh_token, 
		       instance_url   = EXCLUDED.instance_url, 
		       updated_at     = now();`

	if _, err = conn.Exec(query, tid, user.UserID, user.UserName, user.Name, user.FirstName, user.LastName, user.Email, user.FullPhotoURL, user.AccessToken, user.RefreshToken, user.InstanceURL); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func UpdateUserFrozen(conn *sqlx.DB, user *model.User) (err error) {

	const query = `UPDATE public."user" SET sf_is_active = $1, sf_is_frozen = $2 WHERE user_id = $3;`

	if _, err := conn.Exec(query, user.SfIsActive, user.SfIsFrozen, user.UserID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
