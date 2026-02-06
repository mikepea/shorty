# shorty_client.OrganizationsApi

All URIs are relative to *http://localhost:8080/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**organizations_get**](OrganizationsApi.md#organizations_get) | **GET** /organizations | List organizations
[**organizations_id_delete**](OrganizationsApi.md#organizations_id_delete) | **DELETE** /organizations/{id} | Delete an organization
[**organizations_id_get**](OrganizationsApi.md#organizations_id_get) | **GET** /organizations/{id} | Get an organization
[**organizations_id_members_get**](OrganizationsApi.md#organizations_id_members_get) | **GET** /organizations/{id}/members | List organization members
[**organizations_id_members_post**](OrganizationsApi.md#organizations_id_members_post) | **POST** /organizations/{id}/members | Add organization member
[**organizations_id_members_user_id_delete**](OrganizationsApi.md#organizations_id_members_user_id_delete) | **DELETE** /organizations/{id}/members/{userId} | Remove organization member
[**organizations_id_members_user_id_put**](OrganizationsApi.md#organizations_id_members_user_id_put) | **PUT** /organizations/{id}/members/{userId} | Update organization member
[**organizations_id_put**](OrganizationsApi.md#organizations_id_put) | **PUT** /organizations/{id} | Update an organization
[**organizations_post**](OrganizationsApi.md#organizations_post) | **POST** /organizations | Create an organization


# **organizations_get**
> List[OrganizationsOrgResponse] organizations_get()

List organizations

Get all organizations the current user is a member of

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_org_response import OrganizationsOrgResponse
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)

    try:
        # List organizations
        api_response = api_instance.organizations_get()
        print("The response of OrganizationsApi->organizations_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[OrganizationsOrgResponse]**](OrganizationsOrgResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_delete**
> Dict[str, str] organizations_id_delete(id)

Delete an organization

Delete an organization (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID

    try:
        # Delete an organization
        api_response = api_instance.organizations_id_delete(id)
        print("The response of OrganizationsApi->organizations_id_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 

### Return type

**Dict[str, str]**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Organization deleted |  -  |
**403** | Admin access required |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_get**
> OrganizationsOrgResponse organizations_id_get(id)

Get an organization

Get details of a specific organization

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_org_response import OrganizationsOrgResponse
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID

    try:
        # Get an organization
        api_response = api_instance.organizations_id_get(id)
        print("The response of OrganizationsApi->organizations_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 

### Return type

[**OrganizationsOrgResponse**](OrganizationsOrgResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**404** | Organization not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_members_get**
> List[OrganizationsMemberResponse] organizations_id_members_get(id)

List organization members

Get all members of an organization

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_member_response import OrganizationsMemberResponse
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID

    try:
        # List organization members
        api_response = api_instance.organizations_id_members_get(id)
        print("The response of OrganizationsApi->organizations_id_members_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_members_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 

### Return type

[**List[OrganizationsMemberResponse]**](OrganizationsMemberResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**404** | Organization not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_members_post**
> OrganizationsMemberResponse organizations_id_members_post(id, organizations_add_member_request)

Add organization member

Add a user to an organization (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_add_member_request import OrganizationsAddMemberRequest
from shorty_client.models.organizations_member_response import OrganizationsMemberResponse
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID
    organizations_add_member_request = shorty_client.OrganizationsAddMemberRequest() # OrganizationsAddMemberRequest | Member details

    try:
        # Add organization member
        api_response = api_instance.organizations_id_members_post(id, organizations_add_member_request)
        print("The response of OrganizationsApi->organizations_id_members_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_members_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 
 **organizations_add_member_request** | [**OrganizationsAddMemberRequest**](OrganizationsAddMemberRequest.md)| Member details | 

### Return type

[**OrganizationsMemberResponse**](OrganizationsMemberResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created |  -  |
**400** | Validation error |  -  |
**403** | Admin access required |  -  |
**404** | User not found |  -  |
**409** | User is already a member |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_members_user_id_delete**
> Dict[str, str] organizations_id_members_user_id_delete(id, user_id)

Remove organization member

Remove a member from an organization (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID
    user_id = 'user_id_example' # str | User ID

    try:
        # Remove organization member
        api_response = api_instance.organizations_id_members_user_id_delete(id, user_id)
        print("The response of OrganizationsApi->organizations_id_members_user_id_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_members_user_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 
 **user_id** | **str**| User ID | 

### Return type

**Dict[str, str]**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Member removed |  -  |
**403** | Admin access required |  -  |
**404** | Member not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_members_user_id_put**
> OrganizationsMemberResponse organizations_id_members_user_id_put(id, user_id, organizations_update_member_request)

Update organization member

Update a member's role in an organization (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_member_response import OrganizationsMemberResponse
from shorty_client.models.organizations_update_member_request import OrganizationsUpdateMemberRequest
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID
    user_id = 'user_id_example' # str | User ID
    organizations_update_member_request = shorty_client.OrganizationsUpdateMemberRequest() # OrganizationsUpdateMemberRequest | Updated member details

    try:
        # Update organization member
        api_response = api_instance.organizations_id_members_user_id_put(id, user_id, organizations_update_member_request)
        print("The response of OrganizationsApi->organizations_id_members_user_id_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_members_user_id_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 
 **user_id** | **str**| User ID | 
 **organizations_update_member_request** | [**OrganizationsUpdateMemberRequest**](OrganizationsUpdateMemberRequest.md)| Updated member details | 

### Return type

[**OrganizationsMemberResponse**](OrganizationsMemberResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Validation error |  -  |
**403** | Admin access required |  -  |
**404** | Member not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_id_put**
> OrganizationsOrgResponse organizations_id_put(id, organizations_update_org_request)

Update an organization

Update an organization (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_org_response import OrganizationsOrgResponse
from shorty_client.models.organizations_update_org_request import OrganizationsUpdateOrgRequest
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    id = 'id_example' # str | Organization ID
    organizations_update_org_request = shorty_client.OrganizationsUpdateOrgRequest() # OrganizationsUpdateOrgRequest | Updated organization details

    try:
        # Update an organization
        api_response = api_instance.organizations_id_put(id, organizations_update_org_request)
        print("The response of OrganizationsApi->organizations_id_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_id_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Organization ID | 
 **organizations_update_org_request** | [**OrganizationsUpdateOrgRequest**](OrganizationsUpdateOrgRequest.md)| Updated organization details | 

### Return type

[**OrganizationsOrgResponse**](OrganizationsOrgResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Validation error |  -  |
**403** | Admin access required |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **organizations_post**
> OrganizationsOrgResponse organizations_post(organizations_create_org_request)

Create an organization

Create a new organization with the current user as admin

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.organizations_create_org_request import OrganizationsCreateOrgRequest
from shorty_client.models.organizations_org_response import OrganizationsOrgResponse
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.OrganizationsApi(api_client)
    organizations_create_org_request = shorty_client.OrganizationsCreateOrgRequest() # OrganizationsCreateOrgRequest | Organization details

    try:
        # Create an organization
        api_response = api_instance.organizations_post(organizations_create_org_request)
        print("The response of OrganizationsApi->organizations_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling OrganizationsApi->organizations_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **organizations_create_org_request** | [**OrganizationsCreateOrgRequest**](OrganizationsCreateOrgRequest.md)| Organization details | 

### Return type

[**OrganizationsOrgResponse**](OrganizationsOrgResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created |  -  |
**400** | Validation error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

