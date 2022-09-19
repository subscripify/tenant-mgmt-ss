



-- CREATE DATABASE tenants ;


USE tenants;

DROP TABLE IF EXISTS tenant;
CREATE TABLE tenant (
  tenant_uuid                             BINARY(16) NOT NULL UNIQUE PRIMARY KEY,       -- the unique id for the tenant
  org_name                                CHAR(36) NOT NULL,                            -- arbitrary alias used for search and easier search ui - this does not need to be the true name of the org
  top_level_domain                        CHAR(3)  NOT NULL,                            -- the top level domain (eg com, net, io, tv, etc.)
  secondary_domain                        CHAR(36) NOT NULL,                            -- every lord tenant must produce a secondary domain name (e.g. subscripify.com)
  subdomain                               CHAR(36) NOT NULL,                            -- sudomains must be unique for each secondary domain name
  kube_namespace_prefix                   CHAR(36) NOT NULL,                            -- tenants may need more than one k8 namespace depending upon config - but they must all start with this
  lord_services_config                BINARY(16),                                   -- only applicable when tenant type is lord
  super_services_config                   BINARY(16),                                   -- only applicable when tenant type is a lord or super 
  public_services_config                  BINARY(16) NOT NULL,                          -- all tenants have public services this field can not be null
  private_access_config                   BINARY(16),                                   -- lord services can not have a private access config - only available to super tenants
  custom_access_config                    BINARY(16),                                   -- lord services can not have a custom access config - only available to super tenants and main tenants
  subscripify_deployment_cloud_location   CHAR(36),                                     -- Azure, AWS, GCP 
  liege_uuid                              BINARY(16) NOT NULL,                          -- the owning tenant of this tenant - lord tenants this field is equal to tenantUUID
  lord_uuid                               BINARY(16) NOT NULL,                          -- the lord tenant of the secondary domain in which this tenant resides
  is_lord_tenant                          BOOL DEFAULT NULL,                            -- holds an true value if lord tenant otherwise it MUST be null
  is_super_tenant                         BOOL DEFAULT FALSE NOT NULL,                  -- holds true or false to indicate if tenant is a super tenant - this field can not be null
  create_timestamp                        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- the create date 
  created_by                              CHAR(60),                                     -- the identity of the individual or application that created the tenant
  CONSTRAINT fk_lord_services_config
  FOREIGN KEY (lord_services_config)
    REFERENCES tenant_service_configs(tenant_config_uuid),
  CONSTRAINT fk_super_services_config
  FOREIGN KEY (super_services_config)
    REFERENCES tenant_service_configs(tenant_config_uuid),
  CONSTRAINT fk_public_services_config
  FOREIGN KEY (public_services_config)
    REFERENCES tenant_service_configs(tenant_config_uuid),
  CONSTRAINT fk_private_access_config
  FOREIGN KEY (private_access_config)
    REFERENCES access_configs(access_config_uuid),
  CONSTRAINT fk_custom_access_config
  FOREIGN KEY (custom_access_config)
    REFERENCES access_configs(access_config_uuid)
);


-- this ensures that domain names are unique 
CREATE UNIQUE INDEX subdomain_unique ON tenant(subdomain, secondary_domain, top_level_domain );

-- this ensures that k8 namespaces are unique as there is a unique cluster for each secondary domain.top level domain combo
CREATE UNIQUE INDEX k8_namespace_unique ON tenant(secondary_domain, top_level_domain, kube_namespace_prefix);

-- this index requires that the isLordTenant field is either true or null
-- in combination with the force_null_is_lord_* triggers this ensures that there is only one lord tenant 
-- per secondary domain
CREATE UNIQUE INDEX one_lord_per_secondary_domain on tenant(is_lord_tenant, secondary_domain, top_level_domain);


-- these triggers ensure that the isLordTenant field is either true or null. 
-- This is to ensure the one_lord_per_secondary_domain index works as intended
DELIMITER $$
CREATE TRIGGER force_unique_lord_insert
BEFORE INSERT
ON tenant FOR EACH ROW
BEGIN
  IF NEW.is_lord_tenant IS FALSE THEN
	  SET NEW.is_lord_tenant := NULL;
  END IF;
  IF NEW.is_super_tenant IS TRUE AND NEW.is_lord_tenant IS TRUE THEN
    SIGNAL SQLSTATE '45000' SET message_text = 'can not be a Lord Tenant AND a Supper Tenant';
  END IF;
END$$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER force_unique_lord_update
BEFORE UPDATE
ON tenant FOR EACH ROW
BEGIN
  -- ensuring that is_lord_tenant is always NULL or TRUE - never FALSE - this ensures that the one_lord_per_secondary_domain index works as intended
  IF NEW.is_lord_tenant IS FALSE THEN
	  SET NEW.is_lord_tenant := NULL;
  END IF;
  IF NEW.is_super_tenant IS TRUE AND NEW.is_lord_tenant IS TRUE THEN
    SIGNAL SQLSTATE '45000' SET message_text = 'can not be a Lord Tenant AND a Supper Tenant';
  END IF;
END$$
DELIMITER ;

-- these triggers ensure that lords stay lords and supers stay supers - demotions are not allowed
DELIMITER $$
CREATE TRIGGER force_no_demotion_update
BEFORE UPDATE
ON tenant FOR EACH ROW
BEGIN
-- ensuring that lord tenants can not get demoted
  IF (NEW.is_lord_tenant IS FALSE or NEW.is_lord_tenant IS NULL) AND OLD.is_lord_tenant IS TRUE THEN
    SIGNAL SQLSTATE '45000' SET message_text = 'can not change tenant type if lord tenant - instead change config of lord tenant if needed';
  END IF;
  -- ensuring that super tenants can not get demoted
  IF (NEW.is_super_tenant IS FALSE or NEW.is_super_tenant IS NULL) AND OLD.is_super_tenant IS TRUE THEN
    SIGNAL SQLSTATE '45000' SET message_text = 'can not change tenant type if super tenant - instead change config of lord tenant if needed';
  END IF;
END$$


-- these triggers ensure that all fields are requried based upon tenant type
DELIMITER $$
CREATE TRIGGER force_proper_tenant_config_update
BEFORE UPDATE
ON tenant FOR EACH ROW
BEGIN
  
  -- ensuring that lord tenants have the proper configurations
  IF (NEW.is_lord_tenant OR OLD.is_lord_tenant) THEN
    CASE
      -- rules for lord tenatns
      WHEN NEW.lord_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'internal services config is required for lord tenant';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for lord tenant';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for lord tenant';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access configs are for super tenants only';
      WHEN NEW.custom_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs are for super tenants and main tenants only';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that super tenants have the proper configurations
  IF (NEW.is_super_tenant OR OLD.is_super_tenant) THEN
    CASE
      -- rules for super tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'internal services config is not allowed for super tenants';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for super tenant';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for super tenant';
      WHEN NEW.private_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.private_access_config) = 'private') THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config requires private access_type';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public access config requires public access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that main tenants have the proper configurations
  IF ((NEW.is_lord_tenant IS FALSE OR NEW.is_lord_tenant IS NULL) OR (OLD.is_lord_tenant IS FALSE OR OLD.is_lord_tenant IS NULL))
    AND 
    ((NEW.is_super_tenant IS FALSE OR NEW.is_super_tenant IS NULL) OR (OLD.is_super_tenant IS FALSE OR OLD.is_super_tenant IS NULL)) THEN
    CASE
      -- rules for main-tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'internal services config is not allowed for main tenants';
      WHEN NEW.super_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is not allowed for main tenants';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for main tenants';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is not allowed for main tenants';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config is required for main tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public access config requires public access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$

DELIMITER $$
CREATE TRIGGER force_proper_tenant_config_insert
BEFORE INSERT
ON tenant FOR EACH ROW
BEGIN
  
  -- ensuring that lord tenants have the proper configurations
  IF (NEW.is_lord_tenant) THEN
    CASE
      -- rules for lord tenatns
      WHEN NEW.lord_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'internal services config is required for lord tenant';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for lord tenant';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for lord tenant';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access configs are for super tenants only';
      WHEN NEW.custom_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs are for super tenants and main tenants only';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that super tenants have the proper configurations
  IF (NEW.is_super_tenant) THEN
    CASE
	  -- rules for super tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'internal services config is not allowed for super tenants';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for super tenant';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for super tenant';
      WHEN NEW.private_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.private_access_config) = 'private') THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config requires private access_type';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public access config requires public access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that main tenants have the proper configurations
  IF (NEW.is_lord_tenant IS FALSE OR NEW.is_lord_tenant IS NULL) AND (NEW.is_super_tenant IS FALSE OR NEW.is_super_tenant IS NULL) THEN
    CASE
      -- rules for main-tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'internal services config is not allowed for main tenants';
      WHEN NEW.super_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is not allowed for main tenants';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for main tenants';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is not allowed for main tenants';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config is required for main tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public access config requires public access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$

-- DELIMITER $$
-- CREATE TRIGGER force_proper_access_config_insert
-- BEFORE INSERT
-- ON tenant FOR EACH ROW
-- IF NEW.private_access_config IS NOT NULL THEN
-- 	IF !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.private_access_config) = 'private') THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config requires private access_type'; END IF;
-- END IF;
-- IF NEW.custom_access_config IS NOT NULL THEN
-- 	IF !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public access config requires public access_type'; END IF;
-- END IF;
-- END$$