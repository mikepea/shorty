# OrganizationsUpdateOrgRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.organizations_update_org_request import OrganizationsUpdateOrgRequest

# TODO update the JSON string below
json = "{}"
# create an instance of OrganizationsUpdateOrgRequest from a JSON string
organizations_update_org_request_instance = OrganizationsUpdateOrgRequest.from_json(json)
# print the JSON string representation of the object
print(OrganizationsUpdateOrgRequest.to_json())

# convert the object into a dict
organizations_update_org_request_dict = organizations_update_org_request_instance.to_dict()
# create an instance of OrganizationsUpdateOrgRequest from a dict
organizations_update_org_request_from_dict = OrganizationsUpdateOrgRequest.from_dict(organizations_update_org_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


