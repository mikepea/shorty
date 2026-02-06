# shorty_client.GroupsApi

All URIs are relative to *http://localhost:8080/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**groups_get**](GroupsApi.md#groups_get) | **GET** /groups | List groups
[**groups_id_delete**](GroupsApi.md#groups_id_delete) | **DELETE** /groups/{id} | Delete a group
[**groups_id_get**](GroupsApi.md#groups_id_get) | **GET** /groups/{id} | Get a group
[**groups_id_members_get**](GroupsApi.md#groups_id_members_get) | **GET** /groups/{id}/members | List group members
[**groups_id_members_post**](GroupsApi.md#groups_id_members_post) | **POST** /groups/{id}/members | Add group member
[**groups_id_members_user_id_delete**](GroupsApi.md#groups_id_members_user_id_delete) | **DELETE** /groups/{id}/members/{userId} | Remove group member
[**groups_id_members_user_id_put**](GroupsApi.md#groups_id_members_user_id_put) | **PUT** /groups/{id}/members/{userId} | Update group member
[**groups_id_put**](GroupsApi.md#groups_id_put) | **PUT** /groups/{id} | Update a group
[**groups_post**](GroupsApi.md#groups_post) | **POST** /groups | Create a group


# **groups_get**
> List[GroupsGroupResponse] groups_get()

List groups

Get all groups the current user is a member of

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_group_response import GroupsGroupResponse
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
    api_instance = shorty_client.GroupsApi(api_client)

    try:
        # List groups
        api_response = api_instance.groups_get()
        print("The response of GroupsApi->groups_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[GroupsGroupResponse]**](GroupsGroupResponse.md)

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

# **groups_id_delete**
> Dict[str, str] groups_id_delete(id)

Delete a group

Delete a group (requires admin role in group)

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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID

    try:
        # Delete a group
        api_response = api_instance.groups_id_delete(id)
        print("The response of GroupsApi->groups_id_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 

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
**200** | Group deleted |  -  |
**403** | Admin access required |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **groups_id_get**
> GroupsGroupResponse groups_id_get(id)

Get a group

Get details of a specific group

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_group_response import GroupsGroupResponse
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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID

    try:
        # Get a group
        api_response = api_instance.groups_id_get(id)
        print("The response of GroupsApi->groups_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 

### Return type

[**GroupsGroupResponse**](GroupsGroupResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**404** | Group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **groups_id_members_get**
> List[GroupsMemberResponse] groups_id_members_get(id)

List group members

Get all members of a group

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_member_response import GroupsMemberResponse
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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID

    try:
        # List group members
        api_response = api_instance.groups_id_members_get(id)
        print("The response of GroupsApi->groups_id_members_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_members_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 

### Return type

[**List[GroupsMemberResponse]**](GroupsMemberResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**404** | Group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **groups_id_members_post**
> GroupsMemberResponse groups_id_members_post(id, groups_add_member_request)

Add group member

Add a user to a group (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_add_member_request import GroupsAddMemberRequest
from shorty_client.models.groups_member_response import GroupsMemberResponse
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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID
    groups_add_member_request = shorty_client.GroupsAddMemberRequest() # GroupsAddMemberRequest | Member details

    try:
        # Add group member
        api_response = api_instance.groups_id_members_post(id, groups_add_member_request)
        print("The response of GroupsApi->groups_id_members_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_members_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 
 **groups_add_member_request** | [**GroupsAddMemberRequest**](GroupsAddMemberRequest.md)| Member details | 

### Return type

[**GroupsMemberResponse**](GroupsMemberResponse.md)

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

# **groups_id_members_user_id_delete**
> Dict[str, str] groups_id_members_user_id_delete(id, user_id)

Remove group member

Remove a member from a group (requires admin role)

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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID
    user_id = 'user_id_example' # str | User ID

    try:
        # Remove group member
        api_response = api_instance.groups_id_members_user_id_delete(id, user_id)
        print("The response of GroupsApi->groups_id_members_user_id_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_members_user_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 
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

# **groups_id_members_user_id_put**
> GroupsMemberResponse groups_id_members_user_id_put(id, user_id, groups_update_member_request)

Update group member

Update a member's role in a group (requires admin role)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_member_response import GroupsMemberResponse
from shorty_client.models.groups_update_member_request import GroupsUpdateMemberRequest
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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID
    user_id = 'user_id_example' # str | User ID
    groups_update_member_request = shorty_client.GroupsUpdateMemberRequest() # GroupsUpdateMemberRequest | Updated member details

    try:
        # Update group member
        api_response = api_instance.groups_id_members_user_id_put(id, user_id, groups_update_member_request)
        print("The response of GroupsApi->groups_id_members_user_id_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_members_user_id_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 
 **user_id** | **str**| User ID | 
 **groups_update_member_request** | [**GroupsUpdateMemberRequest**](GroupsUpdateMemberRequest.md)| Updated member details | 

### Return type

[**GroupsMemberResponse**](GroupsMemberResponse.md)

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

# **groups_id_put**
> GroupsGroupResponse groups_id_put(id, groups_update_group_request)

Update a group

Update a group (requires admin role in group)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_group_response import GroupsGroupResponse
from shorty_client.models.groups_update_group_request import GroupsUpdateGroupRequest
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
    api_instance = shorty_client.GroupsApi(api_client)
    id = 'id_example' # str | Group ID
    groups_update_group_request = shorty_client.GroupsUpdateGroupRequest() # GroupsUpdateGroupRequest | Updated group details

    try:
        # Update a group
        api_response = api_instance.groups_id_put(id, groups_update_group_request)
        print("The response of GroupsApi->groups_id_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_id_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 
 **groups_update_group_request** | [**GroupsUpdateGroupRequest**](GroupsUpdateGroupRequest.md)| Updated group details | 

### Return type

[**GroupsGroupResponse**](GroupsGroupResponse.md)

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

# **groups_post**
> GroupsGroupResponse groups_post(groups_create_group_request)

Create a group

Create a new group with the current user as admin

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.groups_create_group_request import GroupsCreateGroupRequest
from shorty_client.models.groups_group_response import GroupsGroupResponse
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
    api_instance = shorty_client.GroupsApi(api_client)
    groups_create_group_request = shorty_client.GroupsCreateGroupRequest() # GroupsCreateGroupRequest | Group details

    try:
        # Create a group
        api_response = api_instance.groups_post(groups_create_group_request)
        print("The response of GroupsApi->groups_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling GroupsApi->groups_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groups_create_group_request** | [**GroupsCreateGroupRequest**](GroupsCreateGroupRequest.md)| Group details | 

### Return type

[**GroupsGroupResponse**](GroupsGroupResponse.md)

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

