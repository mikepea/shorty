# AuthAuthResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**token** | **str** |  | [optional] 
**user** | [**AuthUserResponse**](AuthUserResponse.md) |  | [optional] 

## Example

```python
from shorty_client.models.auth_auth_response import AuthAuthResponse

# TODO update the JSON string below
json = "{}"
# create an instance of AuthAuthResponse from a JSON string
auth_auth_response_instance = AuthAuthResponse.from_json(json)
# print the JSON string representation of the object
print(AuthAuthResponse.to_json())

# convert the object into a dict
auth_auth_response_dict = auth_auth_response_instance.to_dict()
# create an instance of AuthAuthResponse from a dict
auth_auth_response_from_dict = AuthAuthResponse.from_dict(auth_auth_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


