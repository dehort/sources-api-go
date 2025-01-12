package dao

import (
	"fmt"

	m "github.com/RedHatInsights/sources-api-go/model"
	"github.com/RedHatInsights/sources-api-go/util"
)

// GetApplicationAuthenticationDao is a function definition that can be replaced in runtime in case some other DAO
// provider is needed.
var GetApplicationAuthenticationDao func(*int64) ApplicationAuthenticationDao

// getDefaultApplicationAuthenticationDao gets the default DAO implementation which will have the given tenant ID.
func getDefaultApplicationAuthenticationDao(tenantId *int64) ApplicationAuthenticationDao {
	return &applicationAuthenticationDaoImpl{
		TenantID: tenantId,
	}
}

// init sets the default DAO implementation so that other packages can request it easily.
func init() {
	GetApplicationAuthenticationDao = getDefaultApplicationAuthenticationDao
}

type applicationAuthenticationDaoImpl struct {
	TenantID *int64
}

func (a *applicationAuthenticationDaoImpl) ApplicationAuthenticationsByApplications(applications []m.Application) ([]m.ApplicationAuthentication, error) {
	var applicationAuthentications []m.ApplicationAuthentication

	applicationIDs := make([]int64, 0)
	for _, value := range applications {
		applicationIDs = append(applicationIDs, value.ID)
	}

	err := DB.Preload("Tenant").Where("application_id IN ?", applicationIDs).Find(&applicationAuthentications).Error
	if err != nil {
		return nil, err
	}

	return applicationAuthentications, nil
}

func (a *applicationAuthenticationDaoImpl) ApplicationAuthenticationsByAuthentications(authentications []m.Authentication) ([]m.ApplicationAuthentication, error) {
	var applicationAuthentications []m.ApplicationAuthentication

	authenticationUIDs := make([]string, 0)
	for _, value := range authentications {
		authenticationUIDs = append(authenticationUIDs, value.ID)
	}

	result := DB.Preload("Tenant").Where("authentication_uid IN ?", authenticationUIDs).Find(&applicationAuthentications)
	if result.Error != nil {
		return nil, result.Error
	}

	return applicationAuthentications, nil
}

func (a *applicationAuthenticationDaoImpl) ApplicationAuthenticationsByResource(resourceType string, applications []m.Application, authentications []m.Authentication) ([]m.ApplicationAuthentication, error) {
	if resourceType == "Source" {
		return a.ApplicationAuthenticationsByApplications(applications)
	}

	return a.ApplicationAuthenticationsByAuthentications(authentications)
}

func (a *applicationAuthenticationDaoImpl) List(limit int, offset int, filters []util.Filter) ([]m.ApplicationAuthentication, int64, error) {
	appAuths := make([]m.ApplicationAuthentication, 0, limit)
	query := DB.Debug().Model(&m.ApplicationAuthentication{}).
		Offset(offset).
		Where("tenant_id = ?", a.TenantID)

	query, err := applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	count := int64(0)
	query.Count(&count)

	result := query.Limit(limit).Find(&appAuths)
	if result.Error != nil {
		return nil, 0, util.NewErrBadRequest(result.Error)
	}
	return appAuths, count, nil
}

func (a *applicationAuthenticationDaoImpl) GetById(id *int64) (*m.ApplicationAuthentication, error) {
	appAuth := &m.ApplicationAuthentication{ID: *id}
	result := DB.First(&appAuth)
	if result.Error != nil {
		return nil, util.NewErrNotFound("application authentication")
	}
	return appAuth, nil
}

func (a *applicationAuthenticationDaoImpl) Create(appAuth *m.ApplicationAuthentication) error {
	result := DB.Create(appAuth)
	return result.Error
}

func (a *applicationAuthenticationDaoImpl) Update(appAuth *m.ApplicationAuthentication) error {
	result := DB.Updates(appAuth)
	return result.Error
}

func (a *applicationAuthenticationDaoImpl) Delete(id *int64) error {
	appAuth := &m.ApplicationAuthentication{ID: *id}
	if result := DB.Delete(appAuth); result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete application id %v", *id)
	}

	return nil
}

func (a *applicationAuthenticationDaoImpl) Tenant() *int64 {
	return a.TenantID
}
