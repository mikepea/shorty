# LinksCreateLinkRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 
**is_public** | **bool** |  | [optional] 
**is_unread** | **bool** |  | [optional] 
**slug** | **str** |  | [optional] 
**title** | **str** |  | [optional] 
**url** | **str** |  | 

## Example

```python
from shorty_client.models.links_create_link_request import LinksCreateLinkRequest

# TODO update the JSON string below
json = "{}"
# create an instance of LinksCreateLinkRequest from a JSON string
links_create_link_request_instance = LinksCreateLinkRequest.from_json(json)
# print the JSON string representation of the object
print(LinksCreateLinkRequest.to_json())

# convert the object into a dict
links_create_link_request_dict = links_create_link_request_instance.to_dict()
# create an instance of LinksCreateLinkRequest from a dict
links_create_link_request_from_dict = LinksCreateLinkRequest.from_dict(links_create_link_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


