# AuthChangePasswordRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**current_password** | **str** |  | 
**new_password** | **str** |  | 

## Example

```python
from shorty_client.models.auth_change_password_request import AuthChangePasswordRequest

# TODO update the JSON string below
json = "{}"
# create an instance of AuthChangePasswordRequest from a JSON string
auth_change_password_request_instance = AuthChangePasswordRequest.from_json(json)
# print the JSON string representation of the object
print(AuthChangePasswordRequest.to_json())

# convert the object into a dict
auth_change_password_request_dict = auth_change_password_request_instance.to_dict()
# create an instance of AuthChangePasswordRequest from a dict
auth_change_password_request_from_dict = AuthChangePasswordRequest.from_dict(auth_change_password_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


