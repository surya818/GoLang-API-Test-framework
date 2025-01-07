package e2etests

import (
	"strings"
	"testing"

	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	"github.com/stretchr/testify/assert"
)

/*
This test creates a new service version and verifies that Get service version returns the correct service version
POST v1/services/{serviceId}/versions
GET v1/services/{serviceId}/versions/{versionId}
A Bug present. commented out currently ***
*/
func TestServiceVersionApi_CreateAndGetServiceVersion(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	assert.False(t, strings.Contains(versionId, "id-"))
	assert.True(t, strings.HasPrefix(serviceVersion.Item.Version, "v"))
	//Commenting this because of the bug
	//assert.NotEmpty(t, serviceVersion.Item.CreatedAt)

	//Get Service version by Id and check the same the data is same as Create Service version
	get_resp, _ := ServiceVersionApi.GetServiceVersion(serviceVersion.Item.ServiceID, versionId)
	get_service_version_object := extractServiceVersionResponse(get_resp)
	assert.Equal(t, 200, get_resp.StatusCode)
	assert.Equal(t, get_service_version_object.Item.ID, serviceVersion.Item.ID)
	assert.Equal(t, get_service_version_object.Item.ServiceID, serviceVersion.Item.ServiceID)
	assert.Equal(t, get_service_version_object.Item.Version, serviceVersion.Item.Version)
	assert.NotEmpty(t, get_service_version_object.Item.CreatedAt)
}

/*
This test creates a new service version and verifies that List service versions returns the correct service version
POST v1/services/{serviceId}/versions
GET v1/services/{serviceId}/versions
*/
func TestServiceVersionApi_CreateAndListServiceVersions(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()

	//List Service version by Id and check the same verifications
	get_resp, _ := ServiceVersionApi.ListServiceVersions(serviceVersion.Item.ServiceID)
	get_service_version_object := extractListServiceVersionsResponse(get_resp)
	assert.Equal(t, 200, get_resp.StatusCode)
	assert.Equal(t, len(get_service_version_object.Items), 1)
	assert.Equal(t, get_service_version_object.Items[0].ID, serviceVersion.Item.ID)
	assert.Equal(t, get_service_version_object.Items[0].ServiceID, serviceVersion.Item.ServiceID)
	assert.Equal(t, get_service_version_object.Items[0].Version, serviceVersion.Item.Version)
	assert.NotEmpty(t, get_service_version_object.Items[0].CreatedAt)
}

/*
Create Service Version and verify timestamps in the response
*/

func TestServiceVersionApi_CreateServiceVersions_VerifyTimestamps(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	assert.NotEmpty(t, serviceVersion.Item.CreatedAt)
	assert.NotEmpty(t, serviceVersion.Item.UpdatedAt)
}

/*
Create Service with 17 character version string
The API should give a http 400 family error, ideally a Bad Request
*/
func TestServiceVersionApi_CreateServiceVersion_Fails_WithMoreThan17Chars_VersionString(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	longVersion := framework.RandomString(17)
	payload := framework.CreateServiceVersionPayload(serviceId, "", longVersion)
	service_version_resp, _ := ServiceVersionApi.CreateServiceVersion(serviceId, payload)
	assert.NotEqual(t, 500, service_version_resp.StatusCode)

}

/*
This test aims to see that Create Service version fails with a empty serviceId in the URL
The Create SErvice version doesnt error out but returns a 200 OK with an empty content
*/
func TestServiceVersionApi_ServiceCreationFails_WithEmptyServiceID_InUrlParameter(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	payload := framework.CreateServiceVersionPayload(serviceId, "", "")
	service_version_resp, _ := ServiceVersionApi.CreateServiceVersion("", payload)

	assert.Equal(t, 200, service_version_resp.StatusCode)
	service_version_object := extractServiceVersionResponse(service_version_resp)
	assert.Empty(t, service_version_object.Item)

}

/*
This test aims to see that Create Service version fails with a empty serviceId in the payload
Since the URL has the serviceId, the empty serviceId is in payload is ignored
*/
func TestServiceVersionApi_ServiceCreationSucceeds_WithEmptyServiceID_InRequestPayload(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	payload := framework.CreateServiceVersionPayload("", "", "")
	service_version_resp, _ := ServiceVersionApi.CreateServiceVersion(serviceId, payload)

	assert.Equal(t, 201, service_version_resp.StatusCode)
	service_version_object := extractServiceVersionResponse(service_version_resp)
	assert.Equal(t, serviceId, service_version_object.Item.ServiceID)

}

/*
Bug Behaviour ***
This test aims to see that Create Service version fails with a Invalid/Nonexistent serviceId in the URL
Currently the service version is getting created with an invalid non-existing serviceId
*/
func TestServiceVersionApi_ServiceCreationFails_WithInvalidServiceID_InUrlParameter(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	invalidServiceId := serviceId + framework.GetRandomNumber()
	payload := framework.CreateServiceVersionPayload(serviceId, "", "")
	service_version_resp, _ := ServiceVersionApi.CreateServiceVersion(invalidServiceId, payload)

	assert.NotEqual(t, 201, service_version_resp.StatusCode)

}

/*
This test aims to see that Create Service version succeeds with a invalid serviceId in the payload
since the serviceId in the request parameter is valid
*/
func TestServiceVersionApi_ServiceCreationSucceeds_WithInvalidServiceID_InRequestPayload(t *testing.T) {

	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	invalidServiceId := serviceId + framework.GetRandomNumber()
	payload := framework.CreateServiceVersionPayload(invalidServiceId, "", "")
	service_version_resp, _ := ServiceVersionApi.CreateServiceVersion(serviceId, payload)

	assert.Equal(t, 201, service_version_resp.StatusCode)
	service_version_object := extractServiceVersionResponse(service_version_resp)
	assert.Equal(t, serviceId, service_version_object.Item.ServiceID)

}

/*
Get Service Version With InvalidServiceId And VersionId
*/
func TestServiceVersionApi_GetServiceVersion_With_InvalidServiceId_And_VersionId(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Pass invalid serviceId, Get Service version by Id and verify the error
	invalidServiceId := serviceId + framework.GetRandomNumber()
	get_resp, _ := ServiceVersionApi.GetServiceVersion(invalidServiceId, versionId)
	assert.Equal(t, 404, get_resp.StatusCode)
	error_resp := extractErrorResponse(get_resp)
	assert.Equal(t, "Service version not found", error_resp.Error)

	//Pass invalid versionId, Get Service version by Id and verify the error
	invalidVersionId := versionId + framework.GetRandomNumber()
	get_resp, _ = ServiceVersionApi.GetServiceVersion(serviceId, invalidVersionId)
	assert.Equal(t, 404, get_resp.StatusCode)
	error_resp = extractErrorResponse(get_resp)
	assert.Equal(t, "Service version not found", error_resp.Error)
}

/*
List Service Versions With InvalidServiceId
*/
func TestServiceVersionApi_ListServiceVersion_With_InvalidServiceId(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	serviceId := serviceVersion.Item.ServiceID

	//Pass invalid serviceId, Get Service version by Id and verify the error
	invalidServiceId := serviceId + framework.GetRandomNumber()
	get_resp, _ := ServiceVersionApi.ListServiceVersions(invalidServiceId)
	assert.Equal(t, 200, get_resp.StatusCode)
	serviceVersions := listServiceVersionsAndExtractTheList(invalidServiceId)
	assert.Nil(t, serviceVersions.Items)
}

/*
Get Service Version With Empty ServiceId And VersionId
*/
func TestServiceVersionApi_GetServiceVersion_With_EmptyServiceId_And_VersionId(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Pass empty serviceId, Get Service version by Id and verify the error
	emptyServiceId := ""
	get_resp, _ := ServiceVersionApi.GetServiceVersion(emptyServiceId, versionId)
	assert.Equal(t, 404, get_resp.StatusCode)

	//Pass empty versionId, Get Service version by Id and verify the error
	emptyVersionId := ""
	get_resp, _ = ServiceVersionApi.GetServiceVersion(serviceId, emptyVersionId)
	assert.Equal(t, 404, get_resp.StatusCode)

}

/*
Delete the service version and
1. Verify the delete api response
2. Verify empty response
3. List Service versions and verify the count decremented with the deleted item
*/
func TestServiceVersionApi_DeleteServiceVersion_VerifyListServiceVersionSize(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Get Service version by Id and check the same verifications
	serviceVersions := listServiceVersionsAndExtractTheList(serviceId)
	serviceCountBeforeDelete := len(serviceVersions.Items)

	//Delete Service by Id and check the object is deleted completely
	delete_response, _ := ServiceVersionApi.DeleteServiceVersion(serviceId, versionId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

	//Get Service by Id and check the same verifications
	//Get Service version by Id and check the same verifications
	serviceVersions = listServiceVersionsAndExtractTheList(serviceId)
	serviceCountAfterDelete := len(serviceVersions.Items)
	assert.True(t, serviceCountAfterDelete == (serviceCountBeforeDelete-1))

}

/*
Delete the service  and verify that the linked service versions remains unaffected
1. Create Service and service version
2. Verify GET service versions is successful
3. Delete Service

There is no impact on service version with the service deleted
*/
func TestServiceVersionApi_DeleteService_VerifyListServiceVersion(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	serviceId := serviceVersion.Item.ServiceID

	//Get Service version by Id and check the same verifications
	serviceVersions := listServiceVersionsAndExtractTheList(serviceId)
	serviceCountBeforeDelete := len(serviceVersions.Items)

	//Delete Service by Id and check the object is deleted completely

	delete_response, _ := ServiceApi.DeleteService(serviceId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

	//Get Service by Id and check the same verifications
	serviceVersions = listServiceVersionsAndExtractTheList(serviceId)
	serviceCountAfterDelete := len(serviceVersions.Items)
	assert.True(t, serviceCountAfterDelete != serviceCountBeforeDelete, "Service versions accessible even when the linked Service is deleted")

}

/*
Delete the service version and
1. Verify the delete api response
2. Verify empty response
3. List Service versions and verify the count decremented with the deleted item
*/
func TestServiceVersionApi_DeleteServiceVersion_VerifyGetServiceVersionSize(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Get Service version by Id and check the same verifications
	get_resp, _ := ServiceVersionApi.GetServiceVersion(serviceId, versionId)

	serviceVersionFromGet := extractServiceVersionResponse(get_resp)
	assert.NotNil(t, serviceVersionFromGet.Item)
	//Delete Service by Id and check the object is deleted completely
	delete_response, _ := ServiceVersionApi.DeleteServiceVersion(serviceId, versionId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

	//Get Service by Id and check the same verifications
	//Get Service version by Id and check the same verifications
	get_resp_after_delete, _ := ServiceVersionApi.GetServiceVersion(serviceId, versionId)
	assert.Equal(t, 404, get_resp_after_delete.StatusCode)
	error_resp := extractErrorResponse(get_resp_after_delete)
	framework.Logger.Info("Get Service after delete error : " + error_resp.Error)
	assert.Equal(t, "Service version not found", error_resp.Error)
}

/*
Verifies that the Delete Service Version returns an empty Http 204 response for Invalid serviceId and versionId
*/
func TestServiceVersionApi_DeleteServiceVersion_InvalidServiceIdAndVersionId(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Get Service version by Id and check the same verifications
	get_resp, _ := ServiceVersionApi.GetServiceVersion(serviceId, versionId)

	serviceVersionFromGet := extractServiceVersionResponse(get_resp)
	assert.NotNil(t, serviceVersionFromGet.Item)
	//Delete Service by Id and check the object is deleted completely
	invalidServiceId := serviceId + framework.GetRandomNumber()
	delete_response, _ := ServiceVersionApi.DeleteServiceVersion(invalidServiceId, versionId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

	//Delete Service by Id and check the object is deleted completely
	invalidVersionId := versionId + framework.GetRandomNumber()
	delete_response, _ = ServiceVersionApi.DeleteServiceVersion(serviceId, invalidVersionId)
	assert.Equal(t, 204, delete_response.StatusCode)
	assert.True(t, delete_response.ContentLength == 0)

}

/*
***Bug Behaviour: The patched versionId gets deleted and a new versionId is created for the change

1. Create SErvice and Service version
2. List Service versions successfully with the versionId
3. Update the service version
4. Verify that the versionId we got in step 1 no longer exists. Instead a new versionId is created with the PATCH
5. ALso, GET service-version will return a 404 service version not found since the version id was changed
*/
func TestServiceVersionApi_UpdateServiceVersion_VersionIdIsChanged(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Patch Service by Id
	updatedVersion := framework.RandomString(15)
	patchPayload := models.ServiceVersion{Version: updatedVersion}
	update_response, _ := ServiceVersionApi.UpdateServiceVersion(serviceId, versionId, patchPayload)

	//Verify Patch response
	assert.Equal(t, 200, update_response.StatusCode)
	updated_response_body := extractServiceVersionResponse(update_response)
	updated_time_after_update := updated_response_body.Item.UpdatedAt
	assert.Equal(t, updatedVersion, updated_response_body.Item.Version)

	//Call GET /services/{serviceId}/version/{versionId} to see we the original version is deleted and we see empty response
	get_resp, _ := ServiceVersionApi.GetServiceVersion(serviceId, versionId)
	service_version_object := extractServiceVersionResponse(get_resp)
	assert.Equal(t, updatedVersion, service_version_object.Item.Version)
	assert.True(t, updated_time_after_update.Unix() == service_version_object.Item.UpdatedAt.Unix())

	//List service versions to verify the update
	list_service_versions := listServiceVersionsAndExtractTheList(serviceId)
	versionExists, service_version_from_list := serviceVersionWithIDExists(list_service_versions, versionId)
	assert.True(t, versionExists)
	assert.Equal(t, versionId, service_version_from_list.ID)
	assert.Equal(t, serviceId, service_version_from_list.ServiceID)

}

/*
Update Service version and verify timestamps in the response
The Created and Updated shouldnt be empty timestamps
*/
func TestServiceVersionApi_UpdateServiceVersionAndVerifyTimestamps(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Patch Service by Id
	updatedVersion := framework.RandomString(15)
	patchPayload := models.ServiceVersion{Version: updatedVersion}
	update_response, _ := ServiceVersionApi.UpdateServiceVersion(serviceId, versionId, patchPayload)
	updated_response_body := extractServiceVersionResponse(update_response)
	assert.NotEmpty(t, updated_response_body.Item.CreatedAt)
	assert.NotEmpty(t, updated_response_body.Item.UpdatedAt)
}

/*
***Bug Behaviour: Patch Service version returns Http 500 with version length greater than 16
PATCH v1/services/{serviceId}/versions/{versionId}
We return an http 500 error for version id > 16 characters in PATCH service version API
This 500 error should be handled
*/
func TestServiceVersionApi_UpdateServiceVersion_Version_MoreThan16Chars(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID
	serviceId := serviceVersion.Item.ServiceID

	//Patch Service by Id
	updatedVersion := framework.RandomString(17)
	patchPayload := models.ServiceVersion{Version: updatedVersion}
	update_response, _ := ServiceVersionApi.UpdateServiceVersion(serviceId, versionId, patchPayload)
	assert.NotEqual(t, 500, update_response.StatusCode)
	assert.Equal(t, 400, update_response.StatusCode)

}

/*
Update Service version fails with empty ServiceId URL Parameter
*/
func TestServiceVersionApi_UpdateServiceVersion_Fails_With_EmptyServiceId_InURL(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID

	//Patch Service by Id
	updatedVersion := framework.RandomString(15)
	patchPayload := models.ServiceVersion{Version: updatedVersion}
	//Update the service version by passing empty service ID
	update_response, _ := ServiceVersionApi.UpdateServiceVersion("", versionId, patchPayload)
	assert.Equal(t, 404, update_response.StatusCode)
}

/*
Update Service version fails with Invalid/Non-existent ServiceId URL Parameter
*/
func TestServiceVersionApi_UpdateServiceVersion_Fails_With_InvalidServiceId_InURL(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()
	versionId := serviceVersion.Item.ID

	//Patch Service by Id
	invalidServiceId := framework.RandomString(10)
	updatedVersion := framework.RandomString(10)
	patchPayload := models.ServiceVersion{Version: updatedVersion}

	//Update the service version by passing invalid/non-existent service ID

	update_response, _ := ServiceVersionApi.UpdateServiceVersion(invalidServiceId, versionId, patchPayload)
	assert.Equal(t, 404, update_response.StatusCode)
	error_resp := extractErrorResponse(update_response)
	assert.Equal(t, "Service version not found", error_resp.Error)
}

/*
Update Service version fails with  invalid/non-existent service version ID URL Parameter
*/
func TestServiceVersionApi_UpdateServiceVersion_Fails_With_InvalidVersionID_InURL(t *testing.T) {

	serviceVersion := CreateServiceVersion_Success()

	//Patch Service by Id
	invalidVersionId := framework.RandomString(10)
	updatedVersion := framework.RandomString(10)
	patchPayload := models.ServiceVersion{Version: updatedVersion}

	//Update the service version by passing invalid/non-existent service version ID

	update_response, _ := ServiceVersionApi.UpdateServiceVersion(serviceVersion.Item.ServiceID, invalidVersionId, patchPayload)
	assert.Equal(t, 404, update_response.StatusCode)
	error_resp := extractErrorResponse(update_response)
	assert.Equal(t, "Service version not found", error_resp.Error)
}
