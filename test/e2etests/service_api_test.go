package e2etests

import (
	"testing"

	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	"github.com/stretchr/testify/assert"
)

/*
This test creates a new service and verifies that Get servic returns the correct service
POST v1/services
GET v1/services/{serviceId}
*/
func TestServiceApi_CreateAndGetService(t *testing.T) {

	service_object := CreateService_Success()
	serviceName := service_object.Item.Name

	//Verify Service Response Id is not what we pass in the request, but an auto generated GUID
	assert.NotEqual(t, serviceName, service_object.Item.ID)
	//Verify Service name is same as we pass in the response
	assert.Equal(t, serviceName, service_object.Item.Name)
	assert.Empty(t, service_object.Item.CreatedAt)
	serviceId := service_object.Item.ID
	assert.NotNil(t, serviceId)

	//Get Service by Id and check the same verifications
	get_resp, _ := ServiceApi.GetService(serviceId)
	service_object = extractServiceResponse(get_resp)
	assert.Equal(t, 200, get_resp.StatusCode)
	assert.Equal(t, serviceName, service_object.Item.Name)
	assert.NotEmpty(t, service_object.Item.CreatedAt)
	assert.Equal(t, serviceId, service_object.Item.ID)

}

/*
Bug Behaviour
This test creates a new service and verifies that Get servic returns the correct service
POST v1/services
GET v1/services/{serviceId}
*/
func TestServiceApi_CreateAndVerifyTimestamps(t *testing.T) {

	service_object := CreateService_Success()
	assert.NotEmpty(t, service_object.Item.CreatedAt)
	assert.NotEmpty(t, service_object.Item.UpdatedAt)

}

/*
This test creates a new service and verifies that the new service is listed in the list of services
POST v1/services
GET v1/services
*/
func TestServiceApi_CreateServiceAndVerifyListServices(t *testing.T) {

	service_object := CreateService_Success()
	serviceName := service_object.Item.Name

	//Verify Service Response Id is not what we pass in the request, but an auto generated GUID
	assert.NotEqual(t, serviceName, service_object.Item.ID)
	//Verify Service name is same as we pass in the response
	assert.Equal(t, serviceName, service_object.Item.Name)
	assert.Empty(t, service_object.Item.CreatedAt)
	serviceId := service_object.Item.ID
	assert.NotNil(t, serviceId)

	//Get Service by Id and check the same verifications
	get_resp, _ := ServiceApi.ListServices()
	list_services := extractListServicesResponse(get_resp)
	serviceExistsInListServices, targetService := serviceWithIDExists(list_services, serviceId)
	assert.True(t, serviceExistsInListServices)
	assert.Equal(t, 200, get_resp.StatusCode)
	assert.Equal(t, serviceName, targetService.Name)
	assert.NotEmpty(t, targetService.CreatedAt)
	assert.Equal(t, serviceId, targetService.ID)

}

/*
This test aims to see that POST /v1/services fails with a payload with an empty name
This test is in place because of the name being required field
*/
func TestServiceApi_ServiceCreationFailsWithEmptyName(t *testing.T) {

	serviceId := framework.GetRandomName("service")
	payload := framework.CreateServicePayload(serviceId, "", "test service")
	service_response, _ := CreateService(payload)
	assert.NotEqual(t, 201, service_response.StatusCode)
}

/*
This test aims to see that POST /v1/services fails with a payload with an empty ServiceID in the payload
This test is in place because of the ID being required field
*/
func TestServiceApi_ServiceCreationFailsWithEmptyId(t *testing.T) {

	serviceName := framework.GetRandomName("service")
	payload := framework.CreateServicePayload("", serviceName, "test service")
	service_response, _ := CreateService(payload)
	assert.NotEqual(t, 201, service_response.StatusCode)
}

// Invoke Get Service with non existent or invalid service ID and expect a 200/404??
func TestServiceApi_GetServiceWithNonExistentId(t *testing.T) {
	serviceId := "invalid"
	service_response, _ := ServiceApi.GetService(serviceId)
	assert.Equal(t, 200, service_response.StatusCode)
	assert.True(t, service_response.ContentLength == 0)
}

// Invoke Get Service with non existent or invalid service ID and expect a 200/404??
func TestServiceApi_GetServiceWithEmptyId(t *testing.T) {
	serviceId := " "
	service_response, _ := ServiceApi.GetService(serviceId)
	assert.Equal(t, 200, service_response.StatusCode)
	assert.True(t, service_response.ContentLength == 0)

}

/*
Delete the service and
1. Verify the delete api response
2. Verify empty response
3. List Services and verify the count decremented with the deleted item
4. Get Deleted Service and verify empty response
*/
func TestServiceApi_DeleteService(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	services := listServicesAndExtractTheList()
	serviceCountBeforeDelete := len(services.Items)

	//Delete Service by Id and check the object is deleted completely
	delete_response, _ := ServiceApi.DeleteService(serviceId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

	//Call GET /services/{serviceId} to see we get an empty response
	//Get Service by Id and check the same verifications
	get_resp, _ := ServiceApi.GetService(serviceId)
	service_object = extractServiceResponse(get_resp)
	assert.Equal(t, 200, get_resp.StatusCode)
	assert.True(t, get_resp.ContentLength == 0)

	//Cal GET /services to see the count decrements by 1
	serviceId = service_object.Item.ID
	services = listServicesAndExtractTheList()
	serviceCountAfterDelete := len(services.Items)
	assert.True(t, serviceCountAfterDelete == serviceCountBeforeDelete-1)
	serviceExists, _ := serviceWithIDExists(services, serviceId)
	assert.False(t, serviceExists)

}

func TestServiceApi_DeleteServiceTwiceAndVerifyHttp204(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID

	//Delete Service by Id and check the object is deleted completely
	delete_response, _ := ServiceApi.DeleteService(serviceId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

	//Delete again and verify the idempotency
	delete_response, _ = ServiceApi.DeleteService(serviceId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

}

// Invoke Get Service with non existent or invalid service ID and expect a 204??
func TestServiceApi_DeleteServiceWithNonExistentId(t *testing.T) {
	serviceId := "invalid"
	service_response, _ := ServiceApi.DeleteService(serviceId)
	assert.Equal(t, 204, service_response.StatusCode)
	assert.True(t, service_response.ContentLength == 0)
}

// Invoke Delete Service with non existent or invalid service ID and expect a 204
func TestServiceApi_DeleteServiceWithEmptyId(t *testing.T) {
	serviceId := " "
	service_response, _ := ServiceApi.DeleteService(serviceId)
	assert.Equal(t, 204, service_response.StatusCode)
	assert.True(t, service_response.ContentLength == 0)

}

/*
Update the service via PATCH /service/{}serviceId  and
1. Verify the Patch api response
2. List Services and verify the update
4. Get Updated Service and verify the updated properties in the response

Bug Behavior:
The Patch response shows the update but on GET and LIST SErvice the changes are not reflected
*/
func TestServiceApi_UpdateService(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	serviceName := service_object.Item.Name
	updated_time := service_object.Item.UpdatedAt

	//Patch Service by Id
	updatedserviceName := "updated-" + serviceName
	updatedDescription := "Updated decription"
	patchPayload := models.Service{Name: updatedserviceName, Description: updatedDescription}
	update_response, _ := ServiceApi.UpdateService(serviceId, patchPayload)

	//Verify Patch response
	assert.Equal(t, 200, update_response.StatusCode)
	updated_response_body := extractServiceResponse(update_response)
	updated_time_after_update := updated_response_body.Item.UpdatedAt
	assert.True(t, updated_time_after_update.Unix() > updated_time.Unix())
	assert.Equal(t, updatedserviceName, updated_response_body.Item.Name)
	assert.Equal(t, updatedDescription, updated_response_body.Item.Description)

	//Call GET /services/{serviceId} to see we get an empty response
	//Get Service by Id and check the same verifications
	get_resp, _ := ServiceApi.GetService(serviceId)
	service_object = extractServiceResponse(get_resp)
	assert.True(t, updated_time_after_update.Unix() == service_object.Item.UpdatedAt.Unix())
	assert.Equal(t, updatedserviceName, service_object.Item.Name)
	assert.Equal(t, updatedDescription, service_object.Item.Description)

	//Cal GET /services to see the service has updates
	serviceId = service_object.Item.ID
	services := listServicesAndExtractTheList()
	serviceExists, service := serviceWithIDExists(services, serviceId)
	assert.True(t, serviceExists)
	assert.True(t, updated_time_after_update.Unix() == service.UpdatedAt.Unix())
	assert.Equal(t, updatedserviceName, service.Name)
	assert.Equal(t, updatedDescription, service.Description)
}
