package sf

import (
	"os"

	force "bitbucket.org/everymind/gforce/lib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-golib/db/dao"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

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
		creds        force.ForceSession
		endpoint     = GetEndpoint(p.ByName("SF_ENVIRONMENT"))
		userID       = p.ByName("SF_USER_ID")
		instanceURL  = p.ByName("SF_INSTANCE_URL")
		accessToken  = p.ByName("SF_ACCESS_TOKEN")
		refreshToken = p.ByName("SF_REFRESH_TOKEN")
	)

	force.CustomEndpoint = instanceURL

	creds = force.ForceSession{
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

	if len(os.Getenv("SF_CLIENT_ID")) > 0 {
		creds.ClientId = os.Getenv("SF_CLIENT_ID")
	}

	f = force.NewForce(&creds)

	return f, nil
}

func NewForceByUser(orgID, userID, accessToken, refreshToken, instanceURL string) (f *force.Force, err error) {
	creds := force.ForceSession{
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

	if len(os.Getenv("SF_CLIENT_ID")) > 0 {
		creds.ClientId = os.Getenv("SF_CLIENT_ID")
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
	}
	params = append(params, accessToken)

	// refresh token
	orgID := model.Parameter{
		TenantID: tid,
		Name:     "SF_REFRESH_TOKEN",
		Value:    f.RefreshToken,
	}
	params = append(params, orgID)

	// userID
	userID := model.Parameter{
		TenantID: tid,
		Name:     "SF_USER_ID",
		Value:    f.UserInfo.UserId,
	}
	params = append(params, userID)

	// instanceUrl
	instanceURL := model.Parameter{
		TenantID: tid,
		Name:     "SF_INSTANCE_URL",
		Value:    f.InstanceUrl,
	}
	params = append(params, instanceURL)

	if err := dao.UpdateParameters(conn, params); err != nil {
		return errors.Wrap(err, "dao.UpdateParameters()")
	}

	return nil
}
