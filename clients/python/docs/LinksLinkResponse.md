# LinksLinkResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**click_count** | **int** |  | [optional] 
**created_at** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**group_id** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**is_public** | **bool** |  | [optional] 
**is_unread** | **bool** |  | [optional] 
**slug** | **str** |  | [optional] 
**title** | **str** |  | [optional] 
**updated_at** | **str** |  | [optional] 
**url** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.links_link_response import LinksLinkResponse

# TODO update the JSON string below
json = "{}"
# create an instance of LinksLinkResponse from a JSON string
links_link_response_instance = LinksLinkResponse.from_json(json)
# print the JSON string representation of the object
print(LinksLinkResponse.to_json())

# convert the object into a dict
links_link_response_dict = links_link_response_instance.to_dict()
# create an instance of LinksLinkResponse from a dict
links_link_response_from_dict = LinksLinkResponse.from_dict(links_link_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


