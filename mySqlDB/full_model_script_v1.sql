-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema tenants
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema tenants
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `tenants` DEFAULT CHARACTER SET latin1 ;
USE `tenants` ;

-- -----------------------------------------------------
-- Table `tenants`.`access_configs`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tenants`.`access_configs` (
  `access_config_uuid` BINARY(16) NOT NULL,
  `config_alias` CHAR(36) NOT NULL,
  `config_location` TEXT NOT NULL,
  `config_owner_tenant` BINARY(16) NULL DEFAULT NULL,
  `created_by` CHAR(60) NOT NULL,
  `access_type` ENUM('custom', 'private') NOT NULL,
  `create_timestamp` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`access_config_uuid`),
  UNIQUE INDEX `access_config_uuid` (`access_config_uuid` ASC) VISIBLE,
  UNIQUE INDEX `config_alias` (`config_alias` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `tenants`.`tenant_service_configs`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tenants`.`tenant_service_configs` (
  `tenant_config_uuid` BINARY(16) NOT NULL,
  `config_alias` CHAR(36) NOT NULL,
  `config_location` TEXT NOT NULL,
  `config_owner_tenant` BINARY(16) NULL DEFAULT NULL,
  `created_by` CHAR(60) NOT NULL,
  `deployment_level` ENUM('lord', 'super', 'public') NOT NULL,
  `create_timestamp` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`tenant_config_uuid`),
  UNIQUE INDEX `tenant_config_uuid` (`tenant_config_uuid` ASC) VISIBLE,
  UNIQUE INDEX `config_alias` (`config_alias` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `tenants`.`tenant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tenants`.`tenant` (
  tenant_uuid                             BINARY(16) NOT NULL UNIQUE PRIMARY KEY,       -- the unique id for the tenant
  tenant_alias                            CHAR(36) NOT NULL,                            -- arbitrary alias used for search and easier search ui - this does not need to be the true name of the org
  top_level_domain                        CHAR(24)  NOT NULL,                            -- the top level domain (eg com, net, io, tv, etc.)
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
  create_timestamp                        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- the create date 
  created_by                              CHAR(60),                                     -- the identity of the individual or application that created the tenant
  PRIMARY KEY (`tenant_uuid`),
  UNIQUE INDEX `tenant_uuid` (`tenant_uuid` ASC) VISIBLE,
  -- this ensures that domain names are unique 
  UNIQUE INDEX `subdomain_unique` (`subdomain` ASC, `secondary_domain` ASC, `top_level_domain` ASC) VISIBLE,
  -- this ensures that k8 namespaces are unique as there is a unique cluster for each secondary domain.top level domain combo
  UNIQUE INDEX `k8_namespace_unique` (`secondary_domain` ASC, `top_level_domain` ASC, `kube_namespace_prefix` ASC) VISIBLE,
  -- this index requires that the isLordTenant field is either true or null
  -- in combination with the force_null_is_lord_* triggers this ensures that there is only one lord tenant 
  -- per secondary domain
  UNIQUE INDEX `one_lord_per_secondary_domain` (`is_lord_tenant` ASC, `secondary_domain` ASC, `top_level_domain` ASC) VISIBLE,
  INDEX `fk_lord_services_config` (`lord_services_config` ASC) VISIBLE,
  INDEX `fk_super_services_config` (`super_services_config` ASC) VISIBLE,
  INDEX `fk_public_services_config` (`public_services_config` ASC) VISIBLE,
  INDEX `fk_private_access_config` (`private_access_config` ASC) VISIBLE,
  INDEX `fk_custom_access_config` (`custom_access_config` ASC) VISIBLE,
  INDEX `srfk_liege_valid` (`liege_uuid` ASC) VISIBLE,
  INDEX `srfk_lord_valid` (`lord_uuid` ASC) VISIBLE,
  CONSTRAINT `fk_custom_access_config`
    FOREIGN KEY (`custom_access_config`)
    REFERENCES `tenants`.`access_configs` (`access_config_uuid`),
  CONSTRAINT `fk_lord_services_config`
    FOREIGN KEY (`lord_services_config`)
    REFERENCES `tenants`.`tenant_service_configs` (`tenant_config_uuid`),
  CONSTRAINT `fk_private_access_config`
    FOREIGN KEY (`private_access_config`)
    REFERENCES `tenants`.`access_configs` (`access_config_uuid`),
  CONSTRAINT `fk_public_services_config`
    FOREIGN KEY (`public_services_config`)
    REFERENCES `tenants`.`tenant_service_configs` (`tenant_config_uuid`),
  CONSTRAINT `fk_super_services_config`
    FOREIGN KEY (`super_services_config`)
    REFERENCES `tenants`.`tenant_service_configs` (`tenant_config_uuid`),
  CONSTRAINT `srfk_liege_valid`
    FOREIGN KEY (`liege_uuid`)
    REFERENCES `tenants`.`tenant` (`tenant_uuid`),
  CONSTRAINT `srfk_lord_valid`
    FOREIGN KEY (`lord_uuid`)
    REFERENCES `tenants`.`tenant` (`tenant_uuid`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `tenants`.`tenant_deleted`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tenants`.`tenant_deleted` (
  `tenant_uuid` BINARY(16) NOT NULL,
  `tenant_alias` CHAR(36) NOT NULL,
  `top_level_domain` CHAR(3) NOT NULL,
  `secondary_domain` CHAR(36) NOT NULL,
  `subdomain` CHAR(36) NOT NULL,
  `kube_namespace_prefix` CHAR(36) NOT NULL,
  `lord_services_config` BINARY(16) NULL DEFAULT NULL,
  `super_services_config` BINARY(16) NULL DEFAULT NULL,
  `public_services_config` BINARY(16) NOT NULL,
  `private_access_config` BINARY(16) NULL DEFAULT NULL,
  `custom_access_config` BINARY(16) NULL DEFAULT NULL,
  `subscripify_deployment_cloud_location` CHAR(36) NULL DEFAULT NULL,
  `liege_uuid` BINARY(16) NULL DEFAULT NULL,
  `lord_uuid` BINARY(16) NULL DEFAULT NULL,
  `is_lord_tenant` TINYINT(1) NULL DEFAULT NULL,
  `is_super_tenant` TINYINT(1) NOT NULL DEFAULT '0',
  `create_timestamp` TIMESTAMP NOT NULL,
  `created_by` CHAR(60) NULL DEFAULT NULL,
  `deleted_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`tenant_uuid`),
  UNIQUE INDEX `tenant_uuid` (`tenant_uuid` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;

USE `tenants` ;

-- -----------------------------------------------------
-- Placeholder table for view `tenants`.`tenant_search`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tenants`.`tenant_search` (`tenant_uuid` INT, `tenant_alias` INT, `subdomain` INT, `tenant_type` INT, `lord_config_alias` INT, `lord_config_UUID` INT, `super_config_alias` INT, `super_config_UUID` INT, `public_config_alias` INT, `public_config_UUID` INT, `private_access_config_alias` INT, `private_access_config_UUID` INT, `custom_access_config_alias` INT, `custom_access_config_UUID` INT);

-- -----------------------------------------------------
-- View `tenants`.`tenant_search`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tenants`.`tenant_search`;
USE `tenants`;
CREATE  OR REPLACE ALGORITHM=UNDEFINED DEFINER=`SfEuj2w4nrHeHoD6`@`%` SQL SECURITY DEFINER VIEW `tenants`.`tenant_search` AS 
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
  tenant.custom_access_config as custom_access_config_UUID
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
END ASC;
USE `tenants`;

DELIMITER $$
USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`copy_deleted_for_recovery`
BEFORE DELETE ON `tenants`.`tenant`
FOR EACH ROW
BEGIN
  INSERT INTO tenant_deleted SELECT *, current_timestamp() as deleted_timestamp	 FROM tenant WHERE tenant_uuid = OLD.tenant_uuid;
END$$

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_no_demotion_update`
BEFORE UPDATE ON `tenants`.`tenant`
FOR EACH ROW
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

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_proper_relationship_insert`
BEFORE INSERT ON `tenants`.`tenant`
FOR EACH ROW
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
      WHEN !((SELECT liege_uuid FROM tenant WHERE tenant_uuid = NEW.liege_uuid) = NEW.lord_uuid) THEN SIGNAL SQLSTATE '45000' SET message_text = 'liege_uuid must belong to the lord_UUID';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_proper_tenant_config_insert`
BEFORE INSERT ON `tenants`.`tenant`
FOR EACH ROW
BEGIN
  
  -- ensuring that lord tenants have the proper configurations
  IF (NEW.is_lord_tenant) THEN
    CASE
      -- rules for lord tenatns
      WHEN NEW.lord_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config is required for lord tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.lord_services_config) = 'lord') THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config requires lord services type';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for lord tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.super_services_config) = 'super') THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config requires super services type';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for lord tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.public_services_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config requires public services type';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access configs are for super tenants only';
      WHEN NEW.custom_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs are for super tenants and main tenants only';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that super tenants have the proper configurations
  IF (NEW.is_super_tenant) THEN
    CASE
	  -- rules for super tenants
	  WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config is not allowed for super tenants';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for super tenant';
	  WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.super_services_config) = 'super') THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config requires super services type';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for super tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.public_services_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config requires public services type';
      WHEN NEW.private_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.private_access_config) = 'private') THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config requires private access_type';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'custom') THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config requires custom access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that main tenants have the proper configurations
  IF (NEW.is_lord_tenant IS FALSE OR NEW.is_lord_tenant IS NULL) AND (NEW.is_super_tenant IS FALSE OR NEW.is_super_tenant IS NULL) THEN
    CASE
      -- rules for main-tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config is not allowed for main tenants';
      WHEN NEW.super_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is not allowed for main tenants';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for main tenants';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.public_services_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config requires public services type';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is not allowed for main tenants';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config is required for main tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'custom') THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config requires custom access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_proper_tenant_config_update`
BEFORE UPDATE ON `tenants`.`tenant`
FOR EACH ROW
BEGIN
  
  -- ensuring that lord tenants have the proper configurations
  IF (NEW.is_lord_tenant OR OLD.is_lord_tenant) THEN
    CASE
      -- rules for lord tenants
      WHEN NEW.lord_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config is required for lord tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.lord_services_config) = 'lord') THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config requires lord services type';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for lord tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.super_services_config) = 'super') THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config requires super services type';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for lord tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.public_services_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config requires public services type';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access configs are for super tenants only';
      WHEN NEW.custom_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs are for super tenants and main tenants only';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that super tenants have the proper configurations
  IF (NEW.is_super_tenant OR OLD.is_super_tenant) THEN
    CASE
      -- rules for super tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config is not allowed for super tenants';
      WHEN NEW.super_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is required for super tenant';
	  WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.super_services_config) = 'super') THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config requires super services type';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for super tenant';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.public_services_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config requires public services type';
      WHEN NEW.private_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.private_access_config) = 'private') THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config requires private access_type';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access configs is required for super tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'custom') THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config requires custom access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
  -- ensuring that main tenants have the proper configurations
  IF ((NEW.is_lord_tenant IS FALSE OR NEW.is_lord_tenant IS NULL) OR (OLD.is_lord_tenant IS FALSE OR OLD.is_lord_tenant IS NULL))
    AND 
    ((NEW.is_super_tenant IS FALSE OR NEW.is_super_tenant IS NULL) OR (OLD.is_super_tenant IS FALSE OR OLD.is_super_tenant IS NULL)) THEN
    CASE
      -- rules for main-tenants
      WHEN NEW.lord_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'lord services config is not allowed for main tenants';
      WHEN NEW.super_services_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'super services config is not allowed for main tenants';
      WHEN NEW.public_services_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config is required for main tenants';
      WHEN !((SELECT deployment_level FROM tenant_service_configs WHERE tenant_config_uuid = NEW.public_services_config) = 'public') THEN SIGNAL SQLSTATE '45000' SET message_text = 'public services config requires public services type';
      WHEN NEW.private_access_config IS NOT NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'private access config is not allowed for main tenants';
      WHEN NEW.custom_access_config IS NULL THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config is required for main tenants ';
      WHEN !((SELECT access_type FROM access_configs WHERE access_config_uuid = NEW.custom_access_config) = 'custom') THEN SIGNAL SQLSTATE '45000' SET message_text = 'custom access config requires custom access_type';
      ELSE BEGIN END;
    END CASE;
  END IF;
END$$

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_relationship_keep_update`
BEFORE UPDATE ON `tenants`.`tenant`
FOR EACH ROW
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

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_unique_lord_insert`
BEFORE INSERT ON `tenants`.`tenant`
FOR EACH ROW
BEGIN
  IF NEW.is_lord_tenant IS FALSE THEN
	  SET NEW.is_lord_tenant := NULL;
  END IF;
  IF NEW.is_super_tenant IS TRUE AND NEW.is_lord_tenant IS TRUE THEN
    SIGNAL SQLSTATE '45000' SET message_text = 'can not be a Lord Tenant AND a Supper Tenant';
  END IF;
END$$

USE `tenants`$$
CREATE
DEFINER=`SfEuj2w4nrHeHoD6`@`%`
TRIGGER `tenants`.`force_unique_lord_update`
BEFORE UPDATE ON `tenants`.`tenant`
FOR EACH ROW
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

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
