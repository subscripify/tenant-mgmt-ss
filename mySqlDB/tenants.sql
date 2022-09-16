CREATE DATABASE tenants ;
USE tenants;
DROP TABLE IF EXISTS tenant;
CREATE TABLE tenant (
	tenantUUID BINARY(16) NOT NULL,
    org_seq_id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    org_name char(36) NOT NULL,
    subdomain char(36) NOT NULL,
    kube_namespace_prefix char(36) NOT NULL,
    tenant_type char(36) NOT NULL,
    internal_services_config char(36),
    super_services_config char(36) NOT NULL,
    public_services_config char(36) NOT NULL,
    private_servcice_config char(36),
    custom_services_config char (36) NOT NULL,
    subscripify_deployment_cloud_location char(36),
    leigeUUID Binary(16) NOT NULL,
    createTimestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tenant
	(
    org_id, 
    org_name, 
    subdomain, 
    kube_namespace_prefix, 
    tenant_type, 
    internal_services_config, 
    super_services_config,
    public_services_config,
    custom_services_config,
    subscripify_deployment_cloud_location
    )
VALUES
	(UUID_TO_BIN(UUID()), 'Subscripify', 'main-tenant', 'main-tenant', 'lord', 'anISConfigId', 'aSSConfigId', 'aPSConfigId', 'aCSConfigId', 'azure' ),
    (UUID_TO_BIN(UUID()), 'kratos', 'kratos-sample', 'kratos-sample', 'super','anISConfigId', 'aSSConfigId', 'aPSConfigId', 'aCSConfigId', 'azure' );

CREATE INDEX org_id ON tenant (org_id);

    
SELECT *, BIN_TO_UUID(org_id) FROM tenant;



-- SELECT BIN_TO_UUID(org_id), org_name, subdomain, kube_namespace_prefix, subscription_type FROM organizations WHERE BIN_TO_UUID(org_id) = '67e1e031-2ef5-11ed-833b-6636daa5a961';