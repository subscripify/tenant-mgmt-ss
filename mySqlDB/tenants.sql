-- CREATE DATABASE tenants ;
USE tenants;
DROP TABLE IF EXISTS tenant;
CREATE TABLE tenant (
  tenantUUID                              BINARY(16) NOT NULL UNIQUE PRIMARY KEY,       -- the uniue id for the tenant
  org_name                                CHAR(36) NOT NULL,                            -- arbitray alias used for search and easier search ui - this does not need to be the true name of the org
  top_level_domain                        CHAR(3)  NOT NULL,                            -- the top level domain (eg com, net, io, tv, etc.)
  secondary_domain                        CHAR(36) NOT NULL,                            -- every lord tenant must produce a secondary domain name (e.g. subscripify.com)
  subdomain                               CHAR(36) NOT NULL,                            -- sudomains must be unique for each secondary domain name
  kube_namespace_prefix                   CHAR(36) NOT NULL,                            -- teants may need more than one k8 namespace depending upon config - but they must all start with this
  tenant_type                             CHAR(36) NOT NULL,                            -- lord, super, or main
  internal_services_config                BINARY(16),                                   -- only applicable when tenant type is lord
  super_services_config                   BINARY(16),                                   -- only applicable when tenant type is a lord or super 
  public_services_config                  BINARY(16) NOT NULL,                          -- all tenants have public services this field can not be null
  private_services_config                 BINARY(16),                                   -- lord services can not have a private config - only aavialable to super tenants
  custom_services_config                  BINARY(16),                                   -- lord services can not have a private config - only aavialable to super tenants and main tenants
  subscripify_deployment_cloud_location   CHAR(36),                                     -- Azure, AWS, GCP 
  leigeUUID                               BINARY(16) NOT NULL,                          -- the owning tenant of this tenant - lord tenants this field is equal to tenantUUID
  lordUUID                                BINARY(16) NOT NULL,                          -- the lord tenant of the secondary domain in which this tenant resides
  isLordTenant                            BOOL DEFAULT NULL,                            -- holds an true value if lord tennat owerwise it MUST be null
  isSuperTenant                           BOOL DEFAULT FALSE NOT NULL,                  -- holds true or false to indicate if tenant is a super tenant - this field can not be null
  createTimestamp                         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- the create date 
  createdBy                               CHAR(60)                                      -- the identity of the individual or applicaiton that created the tenant
);


-- this ensures that domain names are uniquie 
CREATE UNIQUE INDEX subdomain_unique ON tenant(subdomain, secondary_domain, top_level_domain );

-- this ensures that k8 namespaces are uniqe as there is a uniqe cluster for each secondary domain.top levle domain combo
CREATE UNIQUE INDEX k8_namespace_unique ON tenant(secondary_domain, top_level_domain, kube_namespace_prefix);

-- this index requires that the isLordTenant field is either true or null
-- in combination with the force_null_is_lord_* triggers this ensures that there is only one lord tenant 
-- per secondary domain
CREATE UNIQUE INDEX lord_sedondary_domain on tenant(isLordTenant, secondary_domain, top_level_domain);


-- these triggers ensure that the isLordTenant field is either true or null. 
-- This is to ensure the lord_secondary_domain index works as intended
DELIMITER $$
CREATE TRIGGER force_null_is_lord_insert
BEFORE INSERT
ON tenant FOR EACH ROW
BEGIN
  IF NEW.isLordTenant IS FALSE THEN
	SET NEW.isLordTenant := NULL;
  END IF;
END$$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER force_null_is_lord_update
BEFORE UPDATE
ON tenant FOR EACH ROW
BEGIN
  IF NEW.isLordTenant IS FALSE THEN
	SET NEW.isLordTenant := NULL;
  END IF;
END$$
DELIMITER ;

INSERT INTO tenant
	(
    tenantUUID, 
    org_name,
    top_level_domain,
    secondary_domain,
    subdomain, 
    kube_namespace_prefix, 
    tenant_type, 
    internal_services_config, 
    super_services_config,
    public_services_config,
    private_services_config,
    custom_services_config,
    subscripify_deployment_cloud_location,
    leigeUUID,
    lordUUID,
    isLordTenant,
    isSuperTenant,
    createTimestamp,
    createdBy
    )
VALUES
	(UUID_TO_BIN(UUID()), 
    'Subscripify',
    'com',
    'subscripify',
    'lord-tenant', 
    'lord-tenant', 
    'lord', 
    UUID_TO_BIN(UUID()), 
    UUID_TO_BIN(UUID()), 
    UUID_TO_BIN(UUID()), 
    null, 
    null,
    'azure', 
    UUID_TO_BIN(UUID()),
    UUID_TO_BIN(UUID()),
    TRUE,
    FALSE,
    CURDATE(),
    'william.ohara@subscripify.com' );
    
-- UPDATE tenant set isLordTenant = FALSE WHERE createdBy = 'william.ohara@subscripify.com';


    
-- SELECT *, BIN_TO_UUID(tenantUUID) FROM tenant;



-- SELECT BIN_TO_UUID(org_id), org_name, subdomain, kube_namespace_prefix, subscription_type FROM organizations WHERE BIN_TO_UUID(org_id) = '67e1e031-2ef5-11ed-833b-6636daa5a961';