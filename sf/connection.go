package sf

import (
	"database/sql"
	"fmt"

	force "bitbucket.org/everymind/gforce/lib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-golib/db/dao"
	"bitbucket.org/everymind/evmd-golib/db/model"
	"bitbucket.org/everymind/evmd-golib/logger"
)

func NewDummyForce(conn *sqlx.DB, tid int, pType dao.ParameterType) (*force.Force, error) {
	return new(force.Force), nil
}

func NewForce(conn *sqlx.DB, tid int, pType dao.ParameterType) (f *force.Force, err error) {
	p, err := dao.GetParameters(conn, tid, dao.EnumParamNil)
	if err != nil {
		err = errors.Wrap(err, "dao.GetParameters()")
		return
	}

	if len(p) == 0 {
		err = errors.New("parameters not found")
		return
	}

	var (
		creds            force.ForceSession
		clientID         = p.ByName("SF_CLIENT_ID")
		endpoint         = GetEndpoint(p.ByName("SF_ENVIRONMENT"))
		userID           = p.ByName("SF_USER_ID")
		instanceURL      = p.ByName("SF_INSTANCE_URL")
		accessToken      = p.ByName("SF_ACCESS_TOKEN")
		refreshToken     = p.ByName("SF_REFRESH_TOKEN")
		username         = p.ByName("SF_USERNAME")
		password         = fmt.Sprintf("%s%s", p.ByName("SF_PASSWORD"), p.ByName("SF_SECURITY_TOKEN"))
		accessTokenLogin = p.ByName("SF_LOGIN_MODE") == "ACCESS-TOKEN"
	)

	force.CustomEndpoint = instanceURL

	if accessTokenLogin {
		creds = force.ForceSession{
			ClientId:      clientID,
			AccessToken:   accessToken,
			RefreshToken:  refreshToken,
			InstanceUrl:   instanceURL,
			ForceEndpoint: endpoint,
			UserInfo: &force.UserInfo{
				OrgId:  p[0].OrgID,
				UserId: userID,
			},
			SessionOptions: &force.SessionOptions{
				ApiVersion:    force.ApiVersion(),
				RefreshMethod: force.RefreshOauth,
			},
		}
	} else {
		if accessToken == "" {
			creds, err = force.ForceSoapLogin(endpoint, username, password)
			if err != nil {
				err = errors.Wrap(err, "force.ForceSoapLogin()")
				return
			}

			if e := UpdateOrgCredentials(conn, tid, f.Credentials); e != nil {
				e = errors.Wrap(e, "UpdateOrgCredentials()")
				logger.Errorln(e)
			}
		} else {
			creds = force.ForceSession{
				ClientId:      clientID,
				AccessToken:   accessToken,
				InstanceUrl:   instanceURL,
				ForceEndpoint: endpoint,
				UserInfo: &force.UserInfo{
					OrgId:  p[0].OrgID,
					UserId: userID,
				},
				SessionOptions: &force.SessionOptions{
					ApiVersion: force.ApiVersion(),
				},
			}
		}
	}

	f = force.NewForce(&creds)

	return f, nil
}

func NewForceByUser(orgID, clientID, userID, accessToken, refreshToken, instanceURL string) (f *force.Force, err error) {
	creds := force.ForceSession{
		ClientId:      clientID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		InstanceUrl:   instanceURL,
		ForceEndpoint: force.EndpointInstace,
		UserInfo: &force.UserInfo{
			OrgId:  orgID,
			UserId: userID,
		},
		SessionOptions: &force.SessionOptions{
			ApiVersion:    force.ApiVersion(),
			RefreshMethod: force.RefreshOauth,
		},
	}

	f = force.NewForce(&creds)

	return f, nil
}

func UpdateOrgCredentials(conn *sqlx.DB, tid int, f *force.ForceSession) error {
	params := []model.Parameter{}

	// access token
	accessToken := model.Parameter{
		TenantID: tid,
		Name:     "SF_ACCESS_TOKEN",
		Value:    f.AccessToken,
		Description: sql.NullString{
			Valid:  true,
			String: "Salesforce access token",
		},
		Type: "s",
	}
	params = append(params, accessToken)

	// orgID
	orgID := model.Parameter{
		TenantID: tid,
		Name:     "SF_ORG_ID",
		Value:    f.UserInfo.OrgId,
		Description: sql.NullString{
			Valid:  true,
			String: "Salesforce Org ID",
		},
		Type: "s",
	}
	params = append(params, orgID)

	// userID
	userID := model.Parameter{
		TenantID: tid,
		Name:     "SF_USER_ID",
		Value:    f.UserInfo.UserId,
		Description: sql.NullString{
			Valid:  true,
			String: "Salesforce User ID",
		},
		Type: "s",
	}
	params = append(params, userID)

	// instanceUrl
	instanceURL := model.Parameter{
		TenantID: tid,
		Name:     "SF_INSTANCE_URL",
		Value:    f.InstanceUrl,
		Description: sql.NullString{
			Valid:  true,
			String: "Salesforce Instance URL",
		},
		Type: "s",
	}
	params = append(params, instanceURL)

	if err := dao.UpdateParameters(conn, params); err != nil {
		return errors.Wrap(err, "dao.UpdateParameters()")
	}

	return nil
}
