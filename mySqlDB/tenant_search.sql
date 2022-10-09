SELECT * FROM tenants.tenant_search;

explain SELECT * from tenants.tenant_search WHERE 
(tenant_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('64f418f1-de98-4d9e-8f0c-4b1811a5b280')))
AND
(super_config_alias LIKE '%super%' OR super_config_alias LIKE '%last%');


SELECT tenant_UUID, tenant_alias, is_lord_tenant, is_super_tenant FROM tenant_search 
WHERE (custom_access_alias LIKE 'cus1' OR custom_access_alias LIKE 'cus2' ) 
AND (custom_access_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) 
AND (lord_config_alias LIKE 'lordconfigalias1' OR lord_config_alias LIKE 'lordconfigalias2' ) 
AND (lord_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) 
AND (private_access_alias LIKE 'privAC1' OR private_access_alias LIKE 'privac2' ) 
AND (private_access_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) 
AND (public_config_alias LIKE 'publicc1' OR public_config_alias LIKE 'publicc2' ) 
AND (public_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) 
AND (subdomain LIKE 'sub1' OR subdomain LIKE 'sub2' OR subdomain LIKE 'sub3' OR subdomain LIKE 'sub4' ) 
AND (super_config_alias LIKE 'superc1' OR super_config_alias LIKE 'superc2' ) 
AND (super_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) 
AND (tenant_alias LIKE 'tenantAlias1' OR tenant_alias LIKE 'tenantAlias2' ) 
AND (tenant_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1')));

select * from tenant_search;
SELECT * FROM tenant_search WHERE (custom_access_config_alias LIKE 'cus1' OR custom_access_config_alias LIKE 'cus2' ) AND (custom_access_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) AND (lord_config_alias LIKE 'lordconfigalias1' OR lord_config_alias LIKE 'lordconfigalias2' ) AND (lord_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) AND (private_access_config_alias LIKE 'privAC1' OR private_access_config_alias LIKE 'privac2' ) AND (private_access_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) AND (public_config_alias LIKE 'publicc1' OR public_config_alias LIKE 'publicc2' ) AND (public_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) AND (subdomain LIKE 'sub1' OR subdomain LIKE 'sub2' OR subdomain LIKE 'sub3' OR subdomain LIKE 'sub4' ) AND (super_config_alias LIKE 'superc1' OR super_config_alias LIKE 'superc2' ) AND (super_config_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1'))) AND (tenant_alias LIKE 'tenantAlias1' OR tenant_alias LIKE 'tenantAlias2' ) AND (tenant_uuid IN (UUID_TO_BIN('8df030a1-a057-4f53-a011-a2b1cff673a1'),UUID_TO_BIN('9df030a1-a057-4f53-a011-a2b1cff673a1')));