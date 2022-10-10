-- this file contians insert commands for seed data into the access_configs table
INSERT INTO access_configs (
  access_config_uuid,
  config_alias,
  config_location,
  config_owner_tenant,
  created_by,
  access_type
)
VALUES (
  UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
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
  UUID_TO_BIN('845c5fe4-3864-11ed-907f-f5001f9bae96'),
  'FirstPublic',
  'somewhere',
  null,
  'william.ohara@subscripify.com',
  'custom'
);
SELECT *, BIN_TO_UUID(access_config_uuid) FROM access_configs;