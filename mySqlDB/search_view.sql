
CREATE VIEW tenant_search AS
SELECT 
  tenant.tenant_uuid as tenant_uuid, 
  tenant.tenant_alias,
  tenant.subdomain,
  CASE
  WHEN tenant.is_lord_tenant IS TRUE AND tenant.is_super_tenant IS FALSE THEN 'lord'
  WHEN tenant.is_lord_tenant IS NULL AND tenant.is_super_tenant IS TRUE THEN 'super'
  WHEN tenant.is_lord_tenant IS NULL AND tenant.is_super_tenant IS FALSE THEN 'main'
  ELSE 'error'
  END AS tenant_type,
  lord.config_alias AS lord_config_alias, 
  tenant.lord_services_config as lord_config_UUID,
  sup.config_alias AS super_config_alias, 
  tenant.super_services_config as super_config_UUID,
  public.config_alias AS public_config_alias,
  tenant.public_services_config as public_config_UUID,
  private.config_alias as private_access_config_alias,
  tenant.private_access_config as private_access_config_UUID,
  custom.config_alias as custom_access_config_alias,
  tenant.custom_access_config as custom_asscess_config_UUID
  FROM tenant
LEFT JOIN (
SELECT tenant_service_configs.tenant_config_uuid AS config_uuid, tenant_service_configs.config_alias AS config_alias FROM tenant_service_configs WHERE tenant_service_configs.deployment_level = 'lord'
 ) AS lord ON tenant.lord_services_config = lord.config_uuid
LEFT JOIN (
SELECT tenant_service_configs.tenant_config_uuid AS config_uuid, tenant_service_configs.config_alias AS config_alias FROM tenant_service_configs WHERE tenant_service_configs.deployment_level = 'super'
) AS sup ON tenant.super_services_config = sup.config_uuid
LEFT JOIN (
SELECT tenant_service_configs.tenant_config_uuid AS config_uuid, tenant_service_configs.config_alias AS config_alias FROM tenant_service_configs WHERE tenant_service_configs.deployment_level = 'public'
) AS public ON tenant.public_services_config = public.config_uuid
LEFT JOIN (
SELECT access_configs.access_config_uuid AS config_uuid, access_configs.config_alias AS config_alias FROM access_configs WHERE access_configs.access_type = 'custom'
) AS custom ON tenant.custom_access_config = custom.config_uuid
LEFT JOIN (
SELECT access_configs.access_config_uuid AS config_uuid, access_configs.config_alias AS config_alias FROM access_configs WHERE access_configs.access_type = 'private'
) AS private ON tenant.private_access_config = private.config_uuid
ORDER BY tenant.tenant_alias DESC, 
CASE
  WHEN tenant_type = 'lord' THEN 1
  WHEN tenant_type = 'super' THEN 2
  WHEN tenant_type = 'main' THEN 3
  ELSE 4
END ASC ;




