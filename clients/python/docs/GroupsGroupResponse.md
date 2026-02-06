# GroupsGroupResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**member_count** | **int** |  | [optional] 
**name** | **str** |  | [optional] 
**role** | **str** | User&#39;s role in this group | [optional] 

## Example

```python
from shorty_client.models.groups_group_response import GroupsGroupResponse

# TODO update the JSON string below
json = "{}"
# create an instance of GroupsGroupResponse from a JSON string
groups_group_response_instance = GroupsGroupResponse.from_json(json)
# print the JSON string representation of the object
print(GroupsGroupResponse.to_json())

# convert the object into a dict
groups_group_response_dict = groups_group_response_instance.to_dict()
# create an instance of GroupsGroupResponse from a dict
groups_group_response_from_dict = GroupsGroupResponse.from_dict(groups_group_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


