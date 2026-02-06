# TagsTagResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **str** |  | [optional] 
**link_count** | **int** |  | [optional] 
**name** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.tags_tag_response import TagsTagResponse

# TODO update the JSON string below
json = "{}"
# create an instance of TagsTagResponse from a JSON string
tags_tag_response_instance = TagsTagResponse.from_json(json)
# print the JSON string representation of the object
print(TagsTagResponse.to_json())

# convert the object into a dict
tags_tag_response_dict = tags_tag_response_instance.to_dict()
# create an instance of TagsTagResponse from a dict
tags_tag_response_from_dict = TagsTagResponse.from_dict(tags_tag_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


