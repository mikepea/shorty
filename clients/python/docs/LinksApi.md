# shorty_client.LinksApi

All URIs are relative to *http://localhost:8080/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**groups_id_links_get**](LinksApi.md#groups_id_links_get) | **GET** /groups/{id}/links | List links in a group
[**groups_id_links_post**](LinksApi.md#groups_id_links_post) | **POST** /groups/{id}/links | Create a link
[**links_get**](LinksApi.md#links_get) | **GET** /links | Search links
[**links_slug_delete**](LinksApi.md#links_slug_delete) | **DELETE** /links/{slug} | Delete a link
[**links_slug_get**](LinksApi.md#links_slug_get) | **GET** /links/{slug} | Get a link by slug
[**links_slug_put**](LinksApi.md#links_slug_put) | **PUT** /links/{slug} | Update a link


# **groups_id_links_get**
> List[LinksLinkResponse] groups_id_links_get(id, is_unread=is_unread, is_public=is_public)

List links in a group

Get all links belonging to a specific group

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.links_link_response import LinksLinkResponse
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
    api_instance = shorty_client.LinksApi(api_client)
    id = 'id_example' # str | Group ID
    is_unread = True # bool | Filter by unread status (optional)
    is_public = True # bool | Filter by public status (optional)

    try:
        # List links in a group
        api_response = api_instance.groups_id_links_get(id, is_unread=is_unread, is_public=is_public)
        print("The response of LinksApi->groups_id_links_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling LinksApi->groups_id_links_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 
 **is_unread** | **bool**| Filter by unread status | [optional] 
 **is_public** | **bool**| Filter by public status | [optional] 

### Return type

[**List[LinksLinkResponse]**](LinksLinkResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Invalid group ID |  -  |
**404** | Group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **groups_id_links_post**
> LinksLinkResponse groups_id_links_post(id, links_create_link_request)

Create a link

Create a new shortened link in a group

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.links_create_link_request import LinksCreateLinkRequest
from shorty_client.models.links_link_response import LinksLinkResponse
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
    api_instance = shorty_client.LinksApi(api_client)
    id = 'id_example' # str | Group ID
    links_create_link_request = shorty_client.LinksCreateLinkRequest() # LinksCreateLinkRequest | Link details

    try:
        # Create a link
        api_response = api_instance.groups_id_links_post(id, links_create_link_request)
        print("The response of LinksApi->groups_id_links_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling LinksApi->groups_id_links_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 
 **links_create_link_request** | [**LinksCreateLinkRequest**](LinksCreateLinkRequest.md)| Link details | 

### Return type

[**LinksLinkResponse**](LinksLinkResponse.md)

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
**404** | Group not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **links_get**
> List[LinksLinkResponse] links_get(q=q, is_unread=is_unread, is_public=is_public, group_id=group_id, tag=tag, limit=limit, offset=offset)

Search links

Search links across all groups the user has access to

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.links_link_response import LinksLinkResponse
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
    api_instance = shorty_client.LinksApi(api_client)
    q = 'q_example' # str | Search query (searches title, description, URL) (optional)
    is_unread = True # bool | Filter by unread status (optional)
    is_public = True # bool | Filter by public status (optional)
    group_id = 'group_id_example' # str | Filter by group ID (optional)
    tag = 'tag_example' # str | Filter by tag name (optional)
    limit = 56 # int | Max results (default 50, max 100) (optional)
    offset = 56 # int | Offset for pagination (optional)

    try:
        # Search links
        api_response = api_instance.links_get(q=q, is_unread=is_unread, is_public=is_public, group_id=group_id, tag=tag, limit=limit, offset=offset)
        print("The response of LinksApi->links_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling LinksApi->links_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **q** | **str**| Search query (searches title, description, URL) | [optional] 
 **is_unread** | **bool**| Filter by unread status | [optional] 
 **is_public** | **bool**| Filter by public status | [optional] 
 **group_id** | **str**| Filter by group ID | [optional] 
 **tag** | **str**| Filter by tag name | [optional] 
 **limit** | **int**| Max results (default 50, max 100) | [optional] 
 **offset** | **int**| Offset for pagination | [optional] 

### Return type

[**List[LinksLinkResponse]**](LinksLinkResponse.md)

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

# **links_slug_delete**
> Dict[str, str] links_slug_delete(slug)

Delete a link

Delete a link by slug

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
    api_instance = shorty_client.LinksApi(api_client)
    slug = 'slug_example' # str | Link slug

    try:
        # Delete a link
        api_response = api_instance.links_slug_delete(slug)
        print("The response of LinksApi->links_slug_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling LinksApi->links_slug_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 

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
**200** | Link deleted |  -  |
**404** | Link not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **links_slug_get**
> LinksLinkResponse links_slug_get(slug)

Get a link by slug

Get link details by its short slug

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.links_link_response import LinksLinkResponse
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
    api_instance = shorty_client.LinksApi(api_client)
    slug = 'slug_example' # str | Link slug

    try:
        # Get a link by slug
        api_response = api_instance.links_slug_get(slug)
        print("The response of LinksApi->links_slug_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling LinksApi->links_slug_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 

### Return type

[**LinksLinkResponse**](LinksLinkResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**404** | Link not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **links_slug_put**
> LinksLinkResponse links_slug_put(slug, links_update_link_request)

Update a link

Update an existing link by slug

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.links_link_response import LinksLinkResponse
from shorty_client.models.links_update_link_request import LinksUpdateLinkRequest
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
    api_instance = shorty_client.LinksApi(api_client)
    slug = 'slug_example' # str | Link slug
    links_update_link_request = shorty_client.LinksUpdateLinkRequest() # LinksUpdateLinkRequest | Updated link details

    try:
        # Update a link
        api_response = api_instance.links_slug_put(slug, links_update_link_request)
        print("The response of LinksApi->links_slug_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling LinksApi->links_slug_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 
 **links_update_link_request** | [**LinksUpdateLinkRequest**](LinksUpdateLinkRequest.md)| Updated link details | 

### Return type

[**LinksLinkResponse**](LinksLinkResponse.md)

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
**404** | Link not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

