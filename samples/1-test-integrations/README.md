This folder represents a 2 part test:

a-test-resource is meant to be run first and will create integrations
b-test-datasource is meant to be run second and will pull the integrations created in the first step

Test cases:
* Two integrations with the same name but different types should return properly from the integration data source
* Integrations with different sets of settings should be created and destroyed properly
