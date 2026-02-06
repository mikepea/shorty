# GroupsUpdateMemberRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**role** | **str** |  | 

## Example

```python
from shorty_client.models.groups_update_member_request import GroupsUpdateMemberRequest

# TODO update the JSON string below
json = "{}"
# create an instance of GroupsUpdateMemberRequest from a JSON string
groups_update_member_request_instance = GroupsUpdateMemberRequest.from_json(json)
# print the JSON string representation of the object
print(GroupsUpdateMemberRequest.to_json())

# convert the object into a dict
groups_update_member_request_dict = groups_update_member_request_instance.to_dict()
# create an instance of GroupsUpdateMemberRequest from a dict
groups_update_member_request_from_dict = GroupsUpdateMemberRequest.from_dict(groups_update_member_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


