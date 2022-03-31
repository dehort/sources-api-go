package dao

import (
	"encoding/json"
	"fmt"

	m "github.com/RedHatInsights/sources-api-go/model"
	"github.com/RedHatInsights/sources-api-go/util"
)

type authenticationDaoDbImpl struct {
	TenantID *int64
}

func (add *authenticationDaoDbImpl) List(limit, offset int, filters []util.Filter) ([]m.Authentication, int64, error) {
	query := DB.
		Debug().
		Model(&m.Authentication{})

	query, err := applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	// getting the total count (filters included) for pagination
	count := int64(0)
	query.Count(&count)

	// limiting + running the actual query.
	authentications := make([]m.Authentication, 0, limit)
	err = query.
		Debug().
		Where("tenant_id = ?", add.TenantID).
		Limit(limit).
		Offset(offset).
		Find(&authentications).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	marketplaceTokenCacher = GetMarketplaceTokenCacher(add.TenantID)
	for i := 0; i < len(authentications); i++ {
		err := setMarketplaceTokenAuthExtraField(&authentications[i])
		if err != nil {
			return nil, 0, err
		}
	}

	return authentications, count, nil
}

func (add *authenticationDaoDbImpl) GetById(id string) (*m.Authentication, error) {
	authentication := &m.Authentication{}

	err := DB.
		Debug().
		Where("id = ?", id).
		Where("tenant_id = ?", add.TenantID).
		First(&authentication).
		Error

	if err != nil {
		return nil, util.NewErrNotFound("authentication")
	}

	err = setMarketplaceTokenAuthExtraField(authentication)
	if err != nil {
		return nil, err
	}

	return authentication, nil
}

func (add *authenticationDaoDbImpl) ListForSource(sourceID int64, limit, offset int, filters []util.Filter) ([]m.Authentication, int64, error) {
	// Check that the source exists before continuing.
	var sourceExists bool
	err := DB.Debug().
		Model(&m.Source{}).
		Select(`1`).
		Where(`id = ?`, sourceID).
		Where(`tenant_id = ?`, add.TenantID).
		Scan(&sourceExists).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	if !sourceExists {
		return nil, 0, util.NewErrNotFound("source")
	}

	// List and count all the authentications from the given source.
	query := DB.
		Debug().
		Model(&m.Authentication{})

	query, err = applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	// getting the total count (filters included) for pagination
	count := int64(0)
	query.
		Where("resource_id = ?", sourceID).
		Where("resource_type = 'Source'").
		Where("tenant_id = ?", add.TenantID).
		Count(&count)

	// limiting + running the actual query.
	authentications := make([]m.Authentication, 0, limit)
	err = query.
		Limit(limit).
		Offset(offset).
		Find(&authentications).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	for _, auth := range authentications {
		err := setMarketplaceTokenAuthExtraField(&auth)
		if err != nil {
			return nil, 0, err
		}
	}

	return authentications, count, nil
}

func (add *authenticationDaoDbImpl) ListForApplication(applicationID int64, limit, offset int, filters []util.Filter) ([]m.Authentication, int64, error) {
	// Check that the application exists before continuing.
	var applicationExists bool
	err := DB.Debug().
		Model(&m.Application{}).
		Select(`1`).
		Where(`id = ?`, applicationID).
		Where(`tenant_id = ?`, add.TenantID).
		Scan(&applicationExists).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	if !applicationExists {
		return nil, 0, util.NewErrNotFound("application")
	}

	// List and count all the authentications from the given application.
	query := DB.
		Debug().
		Model(&m.Authentication{})

	query, err = applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	// getting the total count (filters included) for pagination
	count := int64(0)
	query.
		Where("resource_id = ?", applicationID).
		Where("resource_type = 'Application'").
		Where("tenant_id = ?", add.TenantID).
		Count(&count)

	// limiting + running the actual query.
	authentications := make([]m.Authentication, 0, limit)
	err = query.
		Where("resource_id = ?", applicationID).
		Where("resource_type = 'Application'").
		Where("tenant_id = ?", add.TenantID).
		Limit(limit).
		Offset(offset).
		Find(&authentications).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	for _, auth := range authentications {
		err := setMarketplaceTokenAuthExtraField(&auth)
		if err != nil {
			return nil, 0, err
		}
	}

	return authentications, count, nil
}

func (add *authenticationDaoDbImpl) ListForApplicationAuthentication(appAuthID int64, limit, offset int, filters []util.Filter) ([]m.Authentication, int64, error) {
	// Check that the application authentication exists before continuing.
	var applicationAuthenticationExists bool
	err := DB.Debug().
		Model(&m.ApplicationAuthentication{}).
		Select(`1`).
		Where(`id = ?`, appAuthID).
		Where(`tenant_id = ?`, add.TenantID).
		Scan(&applicationAuthenticationExists).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	if !applicationAuthenticationExists {
		return nil, 0, util.NewErrNotFound("application authentication")
	}

	// List and count all the authentications from the given application authentications.
	query := DB.Debug().
		Model(&m.Authentication{})

	query, err = applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	// getting the total count (filters included) for pagination
	count := int64(0)
	query.
		Where("resource_id = ?", appAuthID).
		Where("resource_type = 'ApplicationAuthentication'").
		Where("tenant_id = ?", add.TenantID).
		Count(&count)

	// limiting + running the actual query.
	authentications := make([]m.Authentication, 0, limit)
	err = query.
		Where("resource_id = ?", appAuthID).
		Where("resource_type = 'ApplicationAuthentication'").
		Where("tenant_id = ?", add.TenantID).
		Limit(limit).
		Offset(offset).
		Find(&authentications).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	for _, auth := range authentications {
		err := setMarketplaceTokenAuthExtraField(&auth)
		if err != nil {
			return nil, 0, err
		}
	}

	return authentications, count, nil
}

func (add *authenticationDaoDbImpl) ListForEndpoint(endpointID int64, limit, offset int, filters []util.Filter) ([]m.Authentication, int64, error) {
	// Check that the endpoint exists before continuing.
	var endpointExists bool
	err := DB.Debug().
		Model(&m.Endpoint{}).
		Select(`1`).
		Where(`id = ?`, endpointID).
		Where(`tenant_id = ?`, add.TenantID).
		Scan(&endpointExists).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	if !endpointExists {
		return nil, 0, util.NewErrNotFound("endpoint")
	}

	// List and count all the authentications from the given endpoint.
	query := DB.Debug().
		Model(&m.Authentication{})

	query, err = applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	// getting the total count (filters included) for pagination
	count := int64(0)
	query.
		Where("resource_id = ?", endpointID).
		Where("resource_type = 'Endpoint'").
		Where("tenant_id = ?", add.TenantID).
		Count(&count)

	// limiting + running the actual query.
	authentications := make([]m.Authentication, 0, limit)
	err = query.
		Where("resource_id = ?", endpointID).
		Where("resource_type = 'Endpoint'").
		Where("tenant_id = ?", add.TenantID).
		Limit(limit).
		Offset(offset).
		Find(&authentications).
		Error

	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	for _, auth := range authentications {
		err := setMarketplaceTokenAuthExtraField(&auth)
		if err != nil {
			return nil, 0, err
		}
	}

	return authentications, count, nil
}

func (add *authenticationDaoDbImpl) Create(authentication *m.Authentication) error {
	authentication.TenantID = *add.TenantID // the TenantID gets injected in the middleware

	return DB.
		Debug().
		Create(authentication).
		Error
}

// BulkCreate method _without_ checking if the resource exists. Basically since this is the bulk-create method the
// resource doesn't exist yet and we know the source ID is set beforehand.
func (add *authenticationDaoDbImpl) BulkCreate(auth *m.Authentication) error {
	return add.Create(auth)
}

func (add *authenticationDaoDbImpl) Update(authentication *m.Authentication) error {
	return DB.
		Debug().
		Updates(authentication).
		Error
}

func (add *authenticationDaoDbImpl) Delete(id string) (*m.Authentication, error) {
	var authentication m.Authentication

	err := DB.
		Debug().
		Where("id = ?", id).
		Where("tenant_id = ?", add.TenantID).
		First(&authentication).
		Error

	if err != nil {
		return nil, util.NewErrNotFound("authentication")
	}

	err = DB.
		Debug().
		Delete(authentication).
		Error

	if err != nil {
		return nil, fmt.Errorf(`failed to delete authentication with id "%s"`, id)
	}

	return &authentication, nil
}

func (add *authenticationDaoDbImpl) Tenant() *int64 {
	return add.TenantID
}

func (add *authenticationDaoDbImpl) AuthenticationsByResource(authentication *m.Authentication) ([]m.Authentication, error) {
	var err error
	var resourceAuthentications []m.Authentication

	switch authentication.ResourceType {
	case "Source":
		resourceAuthentications, _, err = add.ListForSource(authentication.ResourceID, DEFAULT_LIMIT, DEFAULT_OFFSET, nil)
	case "Endpoint":
		resourceAuthentications, _, err = add.ListForEndpoint(authentication.ResourceID, DEFAULT_LIMIT, DEFAULT_OFFSET, nil)
	case "Application":
		resourceAuthentications, _, err = add.ListForApplication(authentication.ResourceID, DEFAULT_LIMIT, DEFAULT_OFFSET, nil)
	default:
		return nil, fmt.Errorf("unable to fetch authentications for %s", authentication.ResourceType)
	}

	if err != nil {
		return nil, err
	}

	return resourceAuthentications, nil
}

func (add *authenticationDaoDbImpl) BulkMessage(resource util.Resource) (map[string]interface{}, error) {
	add.TenantID = &resource.TenantID
	authentication, err := add.GetById(resource.ResourceUID)
	if err != nil {
		return nil, err
	}

	return BulkMessageFromSource(&authentication.Source, authentication)
}

func (add *authenticationDaoDbImpl) FetchAndUpdateBy(resource util.Resource, updateAttributes map[string]interface{}) error {
	add.TenantID = &resource.TenantID
	authentication, err := add.GetById(resource.ResourceUID)
	if err != nil {
		return err
	}

	err = authentication.UpdateBy(updateAttributes)
	if err != nil {
		return err
	}

	return add.Update(authentication)
}

func (add *authenticationDaoDbImpl) ToEventJSON(resource util.Resource) ([]byte, error) {
	add.TenantID = &resource.TenantID
	auth, err := add.GetById(resource.ResourceUID)
	if err != nil {
		return nil, err
	}

	auth.TenantID = resource.TenantID
	auth.Tenant = m.Tenant{ExternalTenant: resource.AccountNumber}
	authEvent := auth.ToEvent()
	data, err := json.Marshal(authEvent)
	if err != nil {
		return nil, err
	}

	return data, nil
}