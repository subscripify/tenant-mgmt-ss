



-- CREATE DATABASE tenants ;


USE tenants;

DROP TABLE IF EXISTS tenant;
CREATE TABLE tenant (
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
    REFERENCES access_configs(access_config_uuid),
  CONSTRAINT srfk_liege_valid
  FOREIGN KEY (liege_uuid)
    REFERENCES tenant(tenant_uuid),
  CONSTRAINT srfk_lord_valid
  FOREIGN KEY (lord_uuid)
    REFERENCES tenant(tenant_uuid)
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


-- these triggers ensure that all fields are required based upon tenant type
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
DELIMITER ;

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
DELIMITER ;

DELIMITER $$

CREATE TRIGGER force_proper_relationship_insert
BEFORE INSERT
ON tenant FOR EACH ROW
BEGIN
  IF (NEW.is_lord_tenant) THEN
    CASE
      -- relationship rules for lord tenants
      WHEN NEW.lord_uuid IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'when is_lord_tenant is true lord_uuid field must be null';
      WHEN NEW.liege_uuid IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'when is_lord_tenant is true liege_uuid field must be null';
      ELSE BEGIN END;
    END CASE;
  END IF;
  IF (NEW.is_super_tenant) THEN
    CASE
      -- relationship rules for super_tenants
      WHEN NEW.lord_uuid IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord_uuid must have a lord tenant UUID';
      WHEN !(SELECT is_lord_tenant FROM tenant WHERE tenant_uuid = NEW.lord_uuid) THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord_uuid must be the tenant_UUID of a lord tenant';
      WHEN NEW.liege_uuid IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'a super tenant liege_uuid field must have a lord tenant UUID';
      WHEN !(SELECT is_lord_tenant FROM tenant WHERE tenant_uuid = NEW.liege_uuid) THEN SIGNAL SQLSTATE '45000' SET message_text = 'a super tenant liege_uuid must be the tenant_UUID of a lord tenant';
      WHEN NEW.liege_uuid != NEW.lord_uuid THEN SIGNAL SQLSTATE '45000' SET message_text = 'a supper tenant liege_uuid must equal its lord_uuid';
      ELSE BEGIN END;
    END CASE;
  END IF;
  IF (NEW.is_lord_tenant IS FALSE OR NEW.is_lord_tenant IS NULL) AND (NEW.is_super_tenant IS FALSE OR NEW.is_super_tenant IS NULL) THEN
    CASE
      -- relationship rules for main_tenants
      WHEN NEW.lord_uuid IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord_uuid must have a lord tenant UUID';
      WHEN !(SELECT is_lord_tenant FROM tenant WHERE tenant_uuid = NEW.lord_uuid) THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord_uuid must be the tenant_UUID of a lord tenant'; 
      WHEN NEW.liege_uuid IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'a main tenant liege_uuid field must have a super tenant UUID';
      WHEN !(SELECT is_super_tenant FROM tenant WHERE tenant_uuid = NEW.liege_uuid) THEN SIGNAL SQLSTATE '45000' SET message_text = 'a main tenant liege_uuid must be the tenant_UUID of a super tenant';
      WHEN !((SELECT liege_uuid FROM tenants WHERE tenant_uuid = NEW.liege_uuid) = NEW.lord_uuid) THEN SIGNAL SQLSTATE '45000' SET message_text = 'liege_uuid must belong to the lord_UUID';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$
DELIMITER ;

DELIMITER $$

CREATE TRIGGER force_relationship_keep_update
BEFORE UPDATE
ON tenant FOR EACH ROW
BEGIN
  IF (NEW.is_lord_tenant OR OLD.is_lord_tenant) THEN
    CASE
      WHEN NEW.lord_uuid IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'when is_lord_tenant is true lord_uuid field must be null';
      WHEN NEW.liege_uuid IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'when is_lord_tenant is true liege_uuid field must be null';
      ELSE BEGIN END;
    END CASE;
  END IF;
  IF (NEW.is_super_tenant OR OLD.is_super_tenant) THEN
    CASE
      WHEN NEW.lord_UUID != OLD.lord_uuid THEN SIGNAL SQLSTATE '45000' SET message_text = 'can not change lord_UUID relationship - must create a new super tenant under new lord tenant and migrate resources from old to new';
      WHEN NEW.liege_uuid != OLD.liege_uuid THEN SIGNAL SQLSTATE '45000' SET message_text = 'can not change liege_UUID relationship - must create a new super tenant under new lord tenant and migrate resources from old to new';
      ELSE BEGIN END;
    END CASE;
  END IF;
  IF ((NEW.is_lord_tenant IS FALSE OR NEW.is_lord_tenant IS NULL) OR (OLD.is_lord_tenant IS FALSE OR OLD.is_lord_tenant IS NULL))
    AND 
    ((NEW.is_super_tenant IS FALSE OR NEW.is_super_tenant IS NULL) OR (OLD.is_super_tenant IS FALSE OR OLD.is_super_tenant IS NULL)) THEN
    CASE
      WHEN NEW.lord_UUID != OLD.lord_uuid THEN SIGNAL SQLSTATE '45000' SET message_text = 'can not change lord_UUID relationship - must create a new main tenant and super tenant under new lord and migrate resources from old to new';
      WHEN NEW.liege_uuid != OLD.liege_uuid THEN SIGNAL SQLSTATE '45000' SET message_text = 'can not change liege_UUID relationship - must create a new main tenant under new super tenant and migrate resources from old to new';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$
DELIMITER ;

