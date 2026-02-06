# OrganizationsAddMemberRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **str** |  | 
**role** | **str** |  | 

## Example

```python
from shorty_client.models.organizations_add_member_request import OrganizationsAddMemberRequest

# TODO update the JSON string below
json = "{}"
# create an instance of OrganizationsAddMemberRequest from a JSON string
organizations_add_member_request_instance = OrganizationsAddMemberRequest.from_json(json)
# print the JSON string representation of the object
print(OrganizationsAddMemberRequest.to_json())

# convert the object into a dict
organizations_add_member_request_dict = organizations_add_member_request_instance.to_dict()
# create an instance of OrganizationsAddMemberRequest from a dict
organizations_add_member_request_from_dict = OrganizationsAddMemberRequest.from_dict(organizations_add_member_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


