# GroupsCreateGroupRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **str** |  | [optional] 
**name** | **str** |  | 
**organization_id** | **str** | Optional - defaults to org from context or global | [optional] 

## Example

```python
from shorty_client.models.groups_create_group_request import GroupsCreateGroupRequest

# TODO update the JSON string below
json = "{}"
# create an instance of GroupsCreateGroupRequest from a JSON string
groups_create_group_request_instance = GroupsCreateGroupRequest.from_json(json)
# print the JSON string representation of the object
print(GroupsCreateGroupRequest.to_json())

# convert the object into a dict
groups_create_group_request_dict = groups_create_group_request_instance.to_dict()
# create an instance of GroupsCreateGroupRequest from a dict
groups_create_group_request_from_dict = GroupsCreateGroupRequest.from_dict(groups_create_group_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


