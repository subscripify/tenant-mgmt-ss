-- this file contains inserts for seed test data in the tenant_service_configs table
INSERT INTO tenant_service_configs (
  tenant_config_uuid,
  config_alias,
  config_location,
  config_owner_tenant,
  created_by,
  deployment_level
)
VALUES (
  UUID_TO_BIN('60135d7c-3857-11ed-907f-f5001f9bae96'),
  'AzureLord',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'lord'
);

INSERT INTO tenant_service_configs (
  tenant_config_uuid,
  config_alias,
  config_location,
  config_owner_tenant,
  created_by,
  deployment_level
)
VALUES (
  UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'),
  'FirstSuper',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'super'
);



INSERT INTO tenant_service_configs (
  tenant_config_uuid,
  config_alias,
  config_location,
  config_owner_tenant,
  created_by,
  deployment_level
)
VALUES (
  UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'),
  'FirstMain',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'public'
);
-- DELETE FROM tenant_service_configs WHERE tenant_config_uuid=UUID_TO_BIN('c234d752-3859-11ed-907f-f5001f9bae96');
SELECT *, BIN_TO_UUID(tenant_config_uuid) FROM tenant_service_configs;
