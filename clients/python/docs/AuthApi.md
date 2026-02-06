# shorty_client.AuthApi

All URIs are relative to *http://localhost:8080/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**auth_login_post**](AuthApi.md#auth_login_post) | **POST** /auth/login | Login
[**auth_logout_post**](AuthApi.md#auth_logout_post) | **POST** /auth/logout | Logout
[**auth_me_get**](AuthApi.md#auth_me_get) | **GET** /auth/me | Get current user
[**auth_password_put**](AuthApi.md#auth_password_put) | **PUT** /auth/password | Change password
[**auth_register_post**](AuthApi.md#auth_register_post) | **POST** /auth/register | Register a new user


# **auth_login_post**
> AuthAuthResponse auth_login_post(auth_login_request)

Login

Authenticate with email and password to receive a JWT token

### Example


```python
import shorty_client
from shorty_client.models.auth_auth_response import AuthAuthResponse
from shorty_client.models.auth_login_request import AuthLoginRequest
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)


# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.AuthApi(api_client)
    auth_login_request = shorty_client.AuthLoginRequest() # AuthLoginRequest | Login credentials

    try:
        # Login
        api_response = api_instance.auth_login_post(auth_login_request)
        print("The response of AuthApi->auth_login_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->auth_login_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **auth_login_request** | [**AuthLoginRequest**](AuthLoginRequest.md)| Login credentials | 

### Return type

[**AuthAuthResponse**](AuthAuthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Validation error |  -  |
**401** | Invalid credentials |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **auth_logout_post**
> Dict[str, str] auth_logout_post()

Logout

Logout the current user (client-side token invalidation)

### Example


```python
import shorty_client
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)


# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.AuthApi(api_client)

    try:
        # Logout
        api_response = api_instance.auth_logout_post()
        print("The response of AuthApi->auth_logout_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->auth_logout_post: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

**Dict[str, str]**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Logged out successfully |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **auth_me_get**
> AuthUserResponse auth_me_get()

Get current user

Get the authenticated user's profile

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.auth_user_response import AuthUserResponse
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
    api_instance = shorty_client.AuthApi(api_client)

    try:
        # Get current user
        api_response = api_instance.auth_me_get()
        print("The response of AuthApi->auth_me_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->auth_me_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**AuthUserResponse**](AuthUserResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Authentication required |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **auth_password_put**
> Dict[str, str] auth_password_put(auth_change_password_request)

Change password

Change the password for the authenticated user (requires existing password)

### Example

* Api Key Authentication (BearerAuth):

```python
import shorty_client
from shorty_client.models.auth_change_password_request import AuthChangePasswordRequest
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
    api_instance = shorty_client.AuthApi(api_client)
    auth_change_password_request = shorty_client.AuthChangePasswordRequest() # AuthChangePasswordRequest | Password change details

    try:
        # Change password
        api_response = api_instance.auth_password_put(auth_change_password_request)
        print("The response of AuthApi->auth_password_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->auth_password_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **auth_change_password_request** | [**AuthChangePasswordRequest**](AuthChangePasswordRequest.md)| Password change details | 

### Return type

**Dict[str, str]**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Password changed successfully |  -  |
**400** | Validation error or SSO-only user |  -  |
**401** | Authentication required or wrong current password |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **auth_register_post**
> AuthAuthResponse auth_register_post(auth_register_request)

Register a new user

Create a new user account and receive a JWT token

### Example


```python
import shorty_client
from shorty_client.models.auth_auth_response import AuthAuthResponse
from shorty_client.models.auth_register_request import AuthRegisterRequest
from shorty_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api
# See configuration.py for a list of all supported configuration parameters.
configuration = shorty_client.Configuration(
    host = "http://localhost:8080/api"
)


# Enter a context with an instance of the API client
with shorty_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = shorty_client.AuthApi(api_client)
    auth_register_request = shorty_client.AuthRegisterRequest() # AuthRegisterRequest | Registration details

    try:
        # Register a new user
        api_response = api_instance.auth_register_post(auth_register_request)
        print("The response of AuthApi->auth_register_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AuthApi->auth_register_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **auth_register_request** | [**AuthRegisterRequest**](AuthRegisterRequest.md)| Registration details | 

### Return type

[**AuthAuthResponse**](AuthAuthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Created |  -  |
**400** | Validation error |  -  |
**409** | Email already registered |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

