-- DROP TABLE tenant_deleted;
CREATE TABLE tenant_deleted (
  tenant_uuid                             BINARY(16) NOT NULL UNIQUE PRIMARY KEY,       -- the unique id for the tenant
  tenant_alias                            CHAR(36) NOT NULL,                            -- arbitrary alias used for search and easier search ui - this does not need to be the true name of the org
  top_level_domain                        CHAR(3)  NOT NULL,                            -- the top level domain (eg com, net, io, tv, etc.)
  secondary_domain                        CHAR(36) NOT NULL,                            -- every lord tenant must produce a secondary domain name (e.g. subscripify.com)
  subdomain                               CHAR(36) NOT NULL,                            -- sudomains must be unique for each secondary domain name
  kube_namespace_prefix                   CHAR(36) NOT NULL,                            -- tenants may need more than one k8 namespace depending upon config - but they must all start with this
  lord_services_config                    BINARY(16),                                   -- only applicable when tenant type is lord
  super_services_config                   BINARY(16),                                   -- only applicable when tenant type is a lord or super 
  public_services_config                  BINARY(16) NOT NULL,                          -- all tenants have public services this field can not be null
  private_access_config                   BINARY(16),                                   -- lord services can not have a private access config - only available to super tenants
  custom_access_config                    BINARY(16),                                   -- lord services can not have a custom access config - only available to super tenants and main tenants
  subscripify_deployment_cloud_location   CHAR(36),                                     -- Azure, AWS, GCP 
  liege_uuid                              BINARY(16),                                   -- the owning tenant of this tenant - lord tenants this field is equal to tenantUUID
  lord_uuid                               BINARY(16),                                   -- the lord tenant of the secondary domain in which this tenant resides
  is_lord_tenant                          BOOL DEFAULT NULL,                            -- holds an true value if lord tenant otherwise it MUST be null
  is_super_tenant                         BOOL DEFAULT FALSE NOT NULL,                  -- holds true or false to indicate if tenant is a super tenant - this field can not be null
  create_timestamp                        TIMESTAMP NOT NULL,                           -- the create date 
  created_by                              CHAR(60),                                     -- the identity of the individual or application that created the tenant
  deleted_timestamp						  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  -- the time this record was deteted from tenants
);  

select * from tenant_deleted;