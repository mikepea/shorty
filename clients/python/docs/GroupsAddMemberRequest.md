# GroupsAddMemberRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | 
**role** | **str** |  | 

## Example

```python
from shorty_client.models.groups_add_member_request import GroupsAddMemberRequest

# TODO update the JSON string below
json = "{}"
# create an instance of GroupsAddMemberRequest from a JSON string
groups_add_member_request_instance = GroupsAddMemberRequest.from_json(json)
# print the JSON string representation of the object
print(GroupsAddMemberRequest.to_json())

# convert the object into a dict
groups_add_member_request_dict = groups_add_member_request_instance.to_dict()
# create an instance of GroupsAddMemberRequest from a dict
groups_add_member_request_from_dict = GroupsAddMemberRequest.from_dict(groups_add_member_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


