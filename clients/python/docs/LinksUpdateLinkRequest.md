# LinksUpdateLinkRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 
**is_public** | **bool** |  | [optional] 
**is_unread** | **bool** |  | [optional] 
**slug** | **str** |  | [optional] 
**title** | **str** |  | [optional] 
**url** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.links_update_link_request import LinksUpdateLinkRequest

# TODO update the JSON string below
json = "{}"
# create an instance of LinksUpdateLinkRequest from a JSON string
links_update_link_request_instance = LinksUpdateLinkRequest.from_json(json)
# print the JSON string representation of the object
print(LinksUpdateLinkRequest.to_json())

# convert the object into a dict
links_update_link_request_dict = links_update_link_request_instance.to_dict()
# create an instance of LinksUpdateLinkRequest from a dict
links_update_link_request_from_dict = LinksUpdateLinkRequest.from_dict(links_update_link_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


