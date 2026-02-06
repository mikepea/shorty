# GroupsUpdateGroupRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 
**name** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.groups_update_group_request import GroupsUpdateGroupRequest

# TODO update the JSON string below
json = "{}"
# create an instance of GroupsUpdateGroupRequest from a JSON string
groups_update_group_request_instance = GroupsUpdateGroupRequest.from_json(json)
# print the JSON string representation of the object
print(GroupsUpdateGroupRequest.to_json())

# convert the object into a dict
groups_update_group_request_dict = groups_update_group_request_instance.to_dict()
# create an instance of GroupsUpdateGroupRequest from a dict
groups_update_group_request_from_dict = GroupsUpdateGroupRequest.from_dict(groups_update_group_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


