USE tenants;
-- DROP TABLE IF EXISTS tenant_service_configs;
CREATE TABLE tenant_service_configs (
  tenant_config_uuid                      BINARY(16) NOT NULL UNIQUE PRIMARY KEY,
  config_alias                            CHAR(36) NOT NULL UNIQUE,
  config_location                         TEXT NOT NULL,
  config_owner_tenant                     BINARY(16),
  created_by                              CHAR(60) NOT NULL,
  deployment_level                        ENUM('lord','super','public') NOT NULL,
  create_timestamp   					DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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
  UUID_TO_BIN(UUID()),
  'FirstSuper',
  'somewhere',
  UUID_TO_BIN('60135d7c-3857-11ed-907f-f5001f9bae96'),
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
  UUID_TO_BIN(UUID()),
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
  UUID_TO_BIN(UUID()),
  'FirstMain',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'main'
);
-- DELETE FROM tenant_service_configs WHERE tenant_config_uuid=UUID_TO_BIN('c234d752-3859-11ed-907f-f5001f9bae96');
SELECT *, BIN_TO_UUID(tenant_config_uuid) FROM tenant_service_configs;
