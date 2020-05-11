package sf

import (
	"errors"
	"fmt"
	"log"
	"os"

	"bitbucket.org/everymind/gforce"
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db/dao"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//NewForce func
func NewForce(conn *sqlx.DB, tid int, pType dao.ParameterType) (f *gforce.Force, err error) {
	p, err := dao.GetParameters(conn, tid, dao.EnumParamNil)
	if err != nil {
		err = fmt.Errorf("dao.GetParameters(): %w", err)
		return
	}

	if len(p) == 0 {
		err = errors.New("parameters not found")
		return
	}

	var (
		creds        gforce.ForceSession
		endpoint     = GetEndpoint(p.ByName("SF_ENVIRONMENT"))
		userID       = p.ByName("SF_USER_ID")
		instanceURL  = p.ByName("SF_INSTANCE_URL")
		accessToken  = p.ByName("SF_ACCESS_TOKEN")
		refreshToken = p.ByName("SF_REFRESH_TOKEN")
	)

	gforce.CustomEndpoint = instanceURL

	creds = gforce.ForceSession{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		InstanceUrl:   instanceURL,
		ForceEndpoint: endpoint,
		UserInfo: &gforce.UserInfo{
			OrgId:  p[0].OrgID,
			UserId: userID,
		},
		SessionOptions: &gforce.SessionOptions{
			ApiVersion:    gforce.ApiVersion(),
			RefreshMethod: gforce.RefreshOauth,
		},
	}

	if len(os.Getenv("SF_CLIENT_ID")) > 0 {
		creds.ClientId = os.Getenv("SF_CLIENT_ID")
	}

	f = gforce.NewForce(&creds)

	return f, nil
}

//NewJobForce func
func NewJobForce(conn *sqlx.DB, tid int, uid string, pType dao.ParameterType) (f *gforce.Force, err error) {
	p, err := dao.GetParameters(conn, tid, dao.EnumParamNil)
	if err != nil {
		err = fmt.Errorf("dao.GetParameters(): %w", err)
	}

	if len(p) == 0 {
		err = errors.New("parameters not found")
		return
	}

	var user model.User
	tenant, err := dao.GetTenantByID(conn, tid)
	if err != nil {
		return nil, err
	}

	log.Printf("UID - USER: %v - %v", uid, p.ByName("SF_USER_ID"))

	if len(uid) == 0 {
		user, err = dao.GetUser(conn, tid, p.ByName("SF_USER_ID"))
		if err != nil {
			return nil, err
		}
	} else {
		user, err = dao.GetUser(conn, tid, uid)
		if err != nil {
			return nil, err
		}
	}

	instanceEndpoint := p.ByName("SF_INSTANCE_URL")

	var authURL string
	if p.ByName("SF_ENVIRONMENT") == "PRODUCTION" {
		authURL = "https://login.salesforce.com"
	} else {
		authURL = "https://test.salesforce.com"
	}

	session, err := gforce.GetServerAuthorization(p[0].OrgID, tenant.SfClientID, user.Email, authURL, instanceEndpoint)
	if err != nil {
		return nil, err
	}
	log.Printf("Session: %v", session)

	creds := gforce.ForceSession{
		AccessToken:   session.AccessToken,
		RefreshToken:  session.RefreshToken,
		InstanceUrl:   session.InstanceUrl,
		ForceEndpoint: GetEndpoint(p.ByName("SF_ENVIRONMENT")),
		UserInfo: &gforce.UserInfo{
			OrgId:  p[0].OrgID,
			UserId: user.UserID,
		},
		SessionOptions: &gforce.SessionOptions{
			ApiVersion:    gforce.ApiVersion(),
			RefreshMethod: gforce.RefreshUnavailable,
		},
		ClientId: tenant.SfClientID,
	}

	f = gforce.NewForce(&creds)

	log.Printf("Force -> %v", f)

	return f, nil
}

//NewForceByUser func
func NewForceByUser(orgID, userID, accessToken, refreshToken, instanceURL string) (f *gforce.Force, err error) {
	creds := gforce.ForceSession{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		InstanceUrl:   instanceURL,
		ForceEndpoint: gforce.EndpointInstace,
		UserInfo: &gforce.UserInfo{
			OrgId:  orgID,
			UserId: userID,
		},
		SessionOptions: &gforce.SessionOptions{
			ApiVersion:    gforce.ApiVersion(),
			RefreshMethod: gforce.RefreshOauth,
		},
	}

	if len(os.Getenv("SF_CLIENT_ID")) > 0 {
		creds.ClientId = os.Getenv("SF_CLIENT_ID")
	}

	f = gforce.NewForce(&creds)

	return f, nil
}

//UpdateOrgCredentials func
func UpdateOrgCredentials(conn *sqlx.DB, tid int, f *gforce.ForceSession) error {
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
		return fmt.Errorf("dao.UpdateParameters(): %w", err)
	}

	return nil
}
