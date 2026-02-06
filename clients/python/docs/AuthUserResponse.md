# AuthUserResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | [optional] 
**has_password** | **bool** |  | [optional] 
**id** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**system_role** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.auth_user_response import AuthUserResponse

# TODO update the JSON string below
json = "{}"
# create an instance of AuthUserResponse from a JSON string
auth_user_response_instance = AuthUserResponse.from_json(json)
# print the JSON string representation of the object
print(AuthUserResponse.to_json())

# convert the object into a dict
auth_user_response_dict = auth_user_response_instance.to_dict()
# create an instance of AuthUserResponse from a dict
auth_user_response_from_dict = AuthUserResponse.from_dict(auth_user_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


