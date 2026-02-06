# OrganizationsMemberResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **str** |  | [optional] 
**email** | **str** |  | [optional] 
**id** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**role** | **str** |  | [optional] 
**user_id** | **str** |  | [optional] 

## Example

```python
from shorty_client.models.organizations_member_response import OrganizationsMemberResponse

# TODO update the JSON string below
json = "{}"
# create an instance of OrganizationsMemberResponse from a JSON string
organizations_member_response_instance = OrganizationsMemberResponse.from_json(json)
# print the JSON string representation of the object
print(OrganizationsMemberResponse.to_json())

# convert the object into a dict
organizations_member_response_dict = organizations_member_response_instance.to_dict()
# create an instance of OrganizationsMemberResponse from a dict
organizations_member_response_from_dict = OrganizationsMemberResponse.from_dict(organizations_member_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


