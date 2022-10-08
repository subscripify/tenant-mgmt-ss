SELECT * FROM tenants.tenant_search;

SELECT * from tenants.tenant_search WHERE 
tenant_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('64f418f1-de98-4d9e-8f0c-4b1811a5b280'))
AND
super_config_alias LIKE '%First%'