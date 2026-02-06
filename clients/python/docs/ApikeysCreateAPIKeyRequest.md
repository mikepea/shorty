# ApikeysCreateAPIKeyRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.apikeys_create_api_key_request import ApikeysCreateAPIKeyRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ApikeysCreateAPIKeyRequest from a JSON string
apikeys_create_api_key_request_instance = ApikeysCreateAPIKeyRequest.from_json(json)
# print the JSON string representation of the object
print(ApikeysCreateAPIKeyRequest.to_json())

# convert the object into a dict
apikeys_create_api_key_request_dict = apikeys_create_api_key_request_instance.to_dict()
# create an instance of ApikeysCreateAPIKeyRequest from a dict
apikeys_create_api_key_request_from_dict = ApikeysCreateAPIKeyRequest.from_dict(apikeys_create_api_key_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


