# OrganizationsUpdateMemberRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**role** | **str** |  | 

## Example

```python
from shorty_client.models.organizations_update_member_request import OrganizationsUpdateMemberRequest

# TODO update the JSON string below
json = "{}"
# create an instance of OrganizationsUpdateMemberRequest from a JSON string
organizations_update_member_request_instance = OrganizationsUpdateMemberRequest.from_json(json)
# print the JSON string representation of the object
print(OrganizationsUpdateMemberRequest.to_json())

# convert the object into a dict
organizations_update_member_request_dict = organizations_update_member_request_instance.to_dict()
# create an instance of OrganizationsUpdateMemberRequest from a dict
organizations_update_member_request_from_dict = OrganizationsUpdateMemberRequest.from_dict(organizations_update_member_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


