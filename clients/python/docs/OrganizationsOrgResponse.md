# OrganizationsOrgResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**is_global** | **bool** |  | [optional] 
**member_count** | **int** |  | [optional] 
**name** | **str** |  | [optional] 
**role** | **str** |  | [optional] 
**slug** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.organizations_org_response import OrganizationsOrgResponse

# TODO update the JSON string below
json = "{}"
# create an instance of OrganizationsOrgResponse from a JSON string
organizations_org_response_instance = OrganizationsOrgResponse.from_json(json)
# print the JSON string representation of the object
print(OrganizationsOrgResponse.to_json())

# convert the object into a dict
organizations_org_response_dict = organizations_org_response_instance.to_dict()
# create an instance of OrganizationsOrgResponse from a dict
organizations_org_response_from_dict = OrganizationsOrgResponse.from_dict(organizations_org_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


