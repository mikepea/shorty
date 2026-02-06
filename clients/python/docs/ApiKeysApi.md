# shorty_client.ApiKeysApi

All URIs are relative to *http://localhost:8080/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_keys_get**](ApiKeysApi.md#api_keys_get) | **GET** /api-keys | List API keys
[**api_keys_id_delete**](ApiKeysApi.md#api_keys_id_delete) | **DELETE** /api-keys/{id} | Delete an API key
[**api_keys_post**](ApiKeysApi.md#api_keys_post) | **POST** /api-keys | Create an API key


# **api_keys_get**
> List[ApikeysAPIKeyResponse] api_keys_get()

List API keys

Get all API keys for the authenticated user

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.apikeys_api_key_response import ApikeysAPIKeyResponse
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
    api_instance = shorty_client.ApiKeysApi(api_client)

    try:
        # List API keys
        api_response = api_instance.api_keys_get()
        print("The response of ApiKeysApi->api_keys_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ApiKeysApi->api_keys_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[ApikeysAPIKeyResponse]**](ApikeysAPIKeyResponse.md)

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

# **api_keys_id_delete**
> Dict[str, str] api_keys_id_delete(id)

Delete an API key

Delete an API key by ID

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
    api_instance = shorty_client.ApiKeysApi(api_client)
    id = 'id_example' # str | API Key ID

    try:
        # Delete an API key
        api_response = api_instance.api_keys_id_delete(id)
        print("The response of ApiKeysApi->api_keys_id_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ApiKeysApi->api_keys_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| API Key ID | 

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
**200** | API key deleted |  -  |
**404** | API key not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_keys_post**
> ApikeysCreateAPIKeyResponse api_keys_post(apikeys_create_api_key_request)

Create an API key

Create a new API key for the authenticated user

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.apikeys_create_api_key_request import ApikeysCreateAPIKeyRequest
from shorty_client.models.apikeys_create_api_key_response import ApikeysCreateAPIKeyResponse
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
    api_instance = shorty_client.ApiKeysApi(api_client)
    apikeys_create_api_key_request = shorty_client.ApikeysCreateAPIKeyRequest() # ApikeysCreateAPIKeyRequest | API key details

    try:
        # Create an API key
        api_response = api_instance.api_keys_post(apikeys_create_api_key_request)
        print("The response of ApiKeysApi->api_keys_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ApiKeysApi->api_keys_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apikeys_create_api_key_request** | [**ApikeysCreateAPIKeyRequest**](ApikeysCreateAPIKeyRequest.md)| API key details | 

### Return type

[**ApikeysCreateAPIKeyResponse**](ApikeysCreateAPIKeyResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

