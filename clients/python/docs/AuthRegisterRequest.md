# AuthRegisterRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | 
**name** | **str** |  | 
**password** | **str** |  | 

## Example

```python
from shorty_client.models.auth_register_request import AuthRegisterRequest

# TODO update the JSON string below
json = "{}"
# create an instance of AuthRegisterRequest from a JSON string
auth_register_request_instance = AuthRegisterRequest.from_json(json)
# print the JSON string representation of the object
print(AuthRegisterRequest.to_json())

# convert the object into a dict
auth_register_request_dict = auth_register_request_instance.to_dict()
# create an instance of AuthRegisterRequest from a dict
auth_register_request_from_dict = AuthRegisterRequest.from_dict(auth_register_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


