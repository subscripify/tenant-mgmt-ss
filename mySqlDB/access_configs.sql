CREATE TABLE access_configs (
  access_config_uuid                      BINARY(16) NOT NULL UNIQUE PRIMARY KEY, 
  config_alias                            CHAR(36) NOT NULL UNIQUE,
  config_location                         TEXT NOT NULL,
  config_owner_tenant                     BINARY(16),
  created_by                              CHAR(60) NOT NULL,
  access_type                             ENUM("private", "public"),
  create_timestamp                        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP              
);

INSERT INTO access_configs (
  access_config_uuid,
  config_alias,
  config_location,
  config_owner_tenant,
  created_by,
  access_type
)
VALUES (
  UUID_TO_BIN(UUID()),
  'FirstPrivate',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'private'
);

INSERT INTO access_configs (
  access_config_uuid,
  config_alias,
  config_location,
  config_owner_tenant,
  created_by,
  access_type
)
VALUES (
  UUID_TO_BIN(UUID()),
  'FirstPublic',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'public'
);
SELECT *, BIN_TO_UUID(access_config_uuid) FROM access_configs;