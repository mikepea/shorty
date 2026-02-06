# AuthLoginRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | 
**password** | **str** |  | 

## Example

```python
from shorty_client.models.auth_login_request import AuthLoginRequest

# TODO update the JSON string below
json = "{}"
# create an instance of AuthLoginRequest from a JSON string
auth_login_request_instance = AuthLoginRequest.from_json(json)
# print the JSON string representation of the object
print(AuthLoginRequest.to_json())

# convert the object into a dict
auth_login_request_dict = auth_login_request_instance.to_dict()
# create an instance of AuthLoginRequest from a dict
auth_login_request_from_dict = AuthLoginRequest.from_dict(auth_login_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


