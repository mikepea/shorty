# ApikeysAPIKeyResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**key_prefix** | **str** |  | [optional] 
**last_used_at** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.apikeys_api_key_response import ApikeysAPIKeyResponse

# TODO update the JSON string below
json = "{}"
# create an instance of ApikeysAPIKeyResponse from a JSON string
apikeys_api_key_response_instance = ApikeysAPIKeyResponse.from_json(json)
# print the JSON string representation of the object
print(ApikeysAPIKeyResponse.to_json())

# convert the object into a dict
apikeys_api_key_response_dict = apikeys_api_key_response_instance.to_dict()
# create an instance of ApikeysAPIKeyResponse from a dict
apikeys_api_key_response_from_dict = ApikeysAPIKeyResponse.from_dict(apikeys_api_key_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


