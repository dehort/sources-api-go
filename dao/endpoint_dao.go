package dao

import (
	"encoding/json"
	"fmt"

	m "github.com/RedHatInsights/sources-api-go/model"
	"github.com/RedHatInsights/sources-api-go/util"
)

// GetEndpointDao is a function definition that can be replaced in runtime in case some other DAO provider is
// needed.
var GetEndpointDao func(*int64) EndpointDao

// getDefaultAuthenticationDao gets the default DAO implementation which will have the given tenant ID.
func getDefaultEndpointDao(tenantId *int64) EndpointDao {
	return &endpointDaoImpl{
		TenantID: tenantId,
	}
}

// init sets the default DAO implementation so that other packages can request it easily.
func init() {
	GetEndpointDao = getDefaultEndpointDao
}

type endpointDaoImpl struct {
	TenantID *int64
}

func (a *endpointDaoImpl) SubCollectionList(primaryCollection interface{}, limit int, offset int, filters []util.Filter) ([]m.Endpoint, int64, error) {
	endpoints := make([]m.Endpoint, 0, limit)
	sourceType, err := m.NewRelationObject(primaryCollection, *a.TenantID, DB.Debug())
	if err != nil {
		return nil, 0, util.NewErrNotFound("source")
	}

	query := sourceType.HasMany(&m.Endpoint{}, DB.Debug())
	query = query.Where("endpoints.tenant_id = ?", a.TenantID)

	query, err = applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	count := int64(0)
	query.Model(&m.Endpoint{}).Count(&count)

	result := query.Limit(limit).Offset(offset).Find(&endpoints)
	if result.Error != nil {
		return nil, 0, util.NewErrBadRequest(result.Error)
	}

	return endpoints, count, nil
}

func (a *endpointDaoImpl) List(limit int, offset int, filters []util.Filter) ([]m.Endpoint, int64, error) {
	endpoints := make([]m.Endpoint, 0, limit)
	query := DB.Debug().Model(&m.Endpoint{}).
		Offset(offset).
		Where("tenant_id = ?", a.TenantID)

	query, err := applyFilters(query, filters)
	if err != nil {
		return nil, 0, util.NewErrBadRequest(err)
	}

	count := int64(0)
	query.Count(&count)

	result := query.Limit(limit).Find(&endpoints)
	if result.Error != nil {
		return nil, 0, util.NewErrBadRequest(result.Error)
	}

	return endpoints, count, nil
}

func (a *endpointDaoImpl) GetById(id *int64) (*m.Endpoint, error) {
	app := &m.Endpoint{ID: *id}
	result := DB.First(&app)
	if result.Error != nil {
		return nil, util.NewErrNotFound("endpoint")
	}

	return app, nil
}

func (a *endpointDaoImpl) Create(app *m.Endpoint) error {
	app.TenantID = *a.TenantID

	result := DB.Create(app)
	return result.Error
}

func (a *endpointDaoImpl) Update(app *m.Endpoint) error {
	result := DB.Updates(app)
	return result.Error
}

func (a *endpointDaoImpl) Delete(id *int64) (*m.Endpoint, error) {
	endpt := &m.Endpoint{ID: *id}
	result := DB.Where("tenant_id = ?", a.TenantID).First(&endpt)
	if result.Error != nil {
		return nil, util.NewErrNotFound("endpoint")
	}

	if result := DB.Delete(endpt); result.Error != nil {
		return nil, fmt.Errorf("failed to delete endpoint id %v", *id)
	}

	return endpt, nil
}

func (a *endpointDaoImpl) Tenant() *int64 {
	return a.TenantID
}

func (a *endpointDaoImpl) CanEndpointBeSetAsDefaultForSource(sourceId int64) bool {
	endpoint := &m.Endpoint{}

	// add double quotes to the "default" column to avoid any clashes with postgres' "default" keyword
	result := DB.Where(`"default" = true AND source_id = ?`, sourceId).First(&endpoint)
	return result.Error != nil
}

func (a *endpointDaoImpl) IsRoleUniqueForSource(role string, sourceId int64) bool {
	endpoint := &m.Endpoint{}
	result := DB.Where("role = ? AND source_id = ?", role, sourceId).First(&endpoint)

	// If the record doesn't exist "result.Error" will have a "record not found" error
	return result.Error != nil
}

func (a *endpointDaoImpl) SourceHasEndpoints(sourceId int64) bool {
	endpoint := &m.Endpoint{}

	result := DB.Where("source_id = ?", sourceId).First(&endpoint)

	return result.Error == nil
}

func (a *endpointDaoImpl) BulkMessage(resource util.Resource) (map[string]interface{}, error) {
	endpoint := &m.Endpoint{ID: resource.ResourceID}
	result := DB.Preload("Source").Find(&endpoint)

	if result.Error != nil {
		return nil, result.Error
	}

	authentication := &m.Authentication{ResourceID: endpoint.ID, ResourceType: "Endpoint", ApplicationAuthentications: []m.ApplicationAuthentication{}}
	return BulkMessageFromSource(&endpoint.Source, authentication)
}

func (a *endpointDaoImpl) FetchAndUpdateBy(resource util.Resource, updateAttributes map[string]interface{}) error {
	result := DB.Model(&m.Endpoint{ID: resource.ResourceID}).Updates(updateAttributes)
	if result.RowsAffected == 0 {
		return fmt.Errorf("endpoint not found %v", resource)
	}

	return nil
}

func (a *endpointDaoImpl) FindWithTenant(id *int64) (*m.Endpoint, error) {
	endpoint := &m.Endpoint{ID: *id}
	result := DB.Preload("Tenant").Find(&endpoint)

	return endpoint, result.Error
}

func (a *endpointDaoImpl) ToEventJSON(resource util.Resource) ([]byte, error) {
	endpoint, err := a.FindWithTenant(&resource.ResourceID)
	data, _ := json.Marshal(endpoint.ToEvent())

	return data, err
}
