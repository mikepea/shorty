# OrganizationsCreateOrgRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** |  | 
**slug** | **str** |  | 

## Example

```python
from shorty_client.models.organizations_create_org_request import OrganizationsCreateOrgRequest

# TODO update the JSON string below
json = "{}"
# create an instance of OrganizationsCreateOrgRequest from a JSON string
organizations_create_org_request_instance = OrganizationsCreateOrgRequest.from_json(json)
# print the JSON string representation of the object
print(OrganizationsCreateOrgRequest.to_json())

# convert the object into a dict
organizations_create_org_request_dict = organizations_create_org_request_instance.to_dict()
# create an instance of OrganizationsCreateOrgRequest from a dict
organizations_create_org_request_from_dict = OrganizationsCreateOrgRequest.from_dict(organizations_create_org_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


