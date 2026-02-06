# TagsSetTagsRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | **List[str]** |  | 

## Example

```python
from shorty_client.models.tags_set_tags_request import TagsSetTagsRequest

# TODO update the JSON string below
json = "{}"
# create an instance of TagsSetTagsRequest from a JSON string
tags_set_tags_request_instance = TagsSetTagsRequest.from_json(json)
# print the JSON string representation of the object
print(TagsSetTagsRequest.to_json())

# convert the object into a dict
tags_set_tags_request_dict = tags_set_tags_request_instance.to_dict()
# create an instance of TagsSetTagsRequest from a dict
tags_set_tags_request_from_dict = TagsSetTagsRequest.from_dict(tags_set_tags_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


