# GroupsMemberResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**role** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.groups_member_response import GroupsMemberResponse

# TODO update the JSON string below
json = "{}"
# create an instance of GroupsMemberResponse from a JSON string
groups_member_response_instance = GroupsMemberResponse.from_json(json)
# print the JSON string representation of the object
print(GroupsMemberResponse.to_json())

# convert the object into a dict
groups_member_response_dict = groups_member_response_instance.to_dict()
# create an instance of GroupsMemberResponse from a dict
groups_member_response_from_dict = GroupsMemberResponse.from_dict(groups_member_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


