# shorty_client.TagsApi

All URIs are relative to *http://localhost:8080/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**groups_id_tags_get**](TagsApi.md#groups_id_tags_get) | **GET** /groups/{id}/tags | List tags in a group
[**links_slug_tags_get**](TagsApi.md#links_slug_tags_get) | **GET** /links/{slug}/tags | Get link tags
[**links_slug_tags_put**](TagsApi.md#links_slug_tags_put) | **PUT** /links/{slug}/tags | Set link tags
[**links_slug_tags_tag_delete**](TagsApi.md#links_slug_tags_tag_delete) | **DELETE** /links/{slug}/tags/{tag} | Remove tag from link
[**links_slug_tags_tag_post**](TagsApi.md#links_slug_tags_tag_post) | **POST** /links/{slug}/tags/{tag} | Add tag to link
[**tags_get**](TagsApi.md#tags_get) | **GET** /tags | List tags


# **groups_id_tags_get**
> List[TagsTagResponse] groups_id_tags_get(id)

List tags in a group

Get all tags used in a specific group

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.tags_tag_response import TagsTagResponse
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
    api_instance = shorty_client.TagsApi(api_client)
    id = 'id_example' # str | Group ID

    try:
        # List tags in a group
        api_response = api_instance.groups_id_tags_get(id)
        print("The response of TagsApi->groups_id_tags_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TagsApi->groups_id_tags_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Group ID | 

### Return type

[**List[TagsTagResponse]**](TagsTagResponse.md)

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

# **links_slug_tags_get**
> List[TagsTagResponse] links_slug_tags_get(slug)

Get link tags

Get all tags for a specific link

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.tags_tag_response import TagsTagResponse
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
    api_instance = shorty_client.TagsApi(api_client)
    slug = 'slug_example' # str | Link slug

    try:
        # Get link tags
        api_response = api_instance.links_slug_tags_get(slug)
        print("The response of TagsApi->links_slug_tags_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TagsApi->links_slug_tags_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 

### Return type

[**List[TagsTagResponse]**](TagsTagResponse.md)

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

# **links_slug_tags_put**
> List[TagsTagResponse] links_slug_tags_put(slug, tags_set_tags_request)

Set link tags

Replace all tags on a link

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.tags_set_tags_request import TagsSetTagsRequest
from shorty_client.models.tags_tag_response import TagsTagResponse
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
    api_instance = shorty_client.TagsApi(api_client)
    slug = 'slug_example' # str | Link slug
    tags_set_tags_request = shorty_client.TagsSetTagsRequest() # TagsSetTagsRequest | Tags to set

    try:
        # Set link tags
        api_response = api_instance.links_slug_tags_put(slug, tags_set_tags_request)
        print("The response of TagsApi->links_slug_tags_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TagsApi->links_slug_tags_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 
 **tags_set_tags_request** | [**TagsSetTagsRequest**](TagsSetTagsRequest.md)| Tags to set | 

### Return type

[**List[TagsTagResponse]**](TagsTagResponse.md)

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

# **links_slug_tags_tag_delete**
> Dict[str, str] links_slug_tags_tag_delete(slug, tag)

Remove tag from link

Remove a tag from a link

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
    api_instance = shorty_client.TagsApi(api_client)
    slug = 'slug_example' # str | Link slug
    tag = 'tag_example' # str | Tag name

    try:
        # Remove tag from link
        api_response = api_instance.links_slug_tags_tag_delete(slug, tag)
        print("The response of TagsApi->links_slug_tags_tag_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TagsApi->links_slug_tags_tag_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 
 **tag** | **str**| Tag name | 

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
**200** | Tag removed |  -  |
**404** | Link or tag not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **links_slug_tags_tag_post**
> TagsTagResponse links_slug_tags_tag_post(slug, tag)

Add tag to link

Add a single tag to a link

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.tags_tag_response import TagsTagResponse
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
    api_instance = shorty_client.TagsApi(api_client)
    slug = 'slug_example' # str | Link slug
    tag = 'tag_example' # str | Tag name

    try:
        # Add tag to link
        api_response = api_instance.links_slug_tags_tag_post(slug, tag)
        print("The response of TagsApi->links_slug_tags_tag_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TagsApi->links_slug_tags_tag_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **slug** | **str**| Link slug | 
 **tag** | **str**| Tag name | 

### Return type

[**TagsTagResponse**](TagsTagResponse.md)

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

# **tags_get**
> List[TagsTagResponse] tags_get()

List tags

Get all tags used across the user's groups

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.tags_tag_response import TagsTagResponse
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
    api_instance = shorty_client.TagsApi(api_client)

    try:
        # List tags
        api_response = api_instance.tags_get()
        print("The response of TagsApi->tags_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TagsApi->tags_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[TagsTagResponse]**](TagsTagResponse.md)

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

