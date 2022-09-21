
INSERT INTO tenant (
	tenant_uuid, 
    org_name,
    top_level_domain,
    secondary_domain,
    subdomain, 
    kube_namespace_prefix, 
    lord_services_config, 
    super_services_config,
    public_services_config,
    private_access_config,
    custom_access_config,
    subscripify_deployment_cloud_location,
    liege_uuid,
    lord_uuid,
    is_lord_tenant,
    is_super_tenant,
    create_timestamp,
    created_by
    )
VALUES
	(
	 UUID_TO_BIN(UUID()), 
    'Subscripify',
    'com',
    'subscripify',
    'lord-tenant', 
    'lord-tenant', 
    UUID_TO_BIN('60135d7c-3857-11ed-907f-f5001f9bae96'), 
    UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'), 
    UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'), 
    null, 
    null,
    'azure', 
    UUID_TO_BIN(UUID()),
    UUID_TO_BIN(UUID()),
    TRUE,
    FALSE,
    CURDATE(),
    'william.ohara@subscripify.com' );
    
    -- this one should fail
    INSERT INTO tenant (
	tenant_uuid, 
    org_name,
    top_level_domain,
    secondary_domain,
    subdomain, 
    kube_namespace_prefix, 
    lord_services_config, 
    super_services_config,
    public_services_config,
    private_access_config,
    custom_access_config,
    subscripify_deployment_cloud_location,
    liege_uuid,
    lord_uuid,
    is_lord_tenant,
    is_super_tenant,
    create_timestamp,
    created_by
    )
VALUES
	(
	 UUID_TO_BIN(UUID()), 
    'Subscripify',
    'com',
    'subscripify',
    'super-tenant', 
    'super-tenant', 
    null, 
    UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'), 
    UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'), 
	UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
    UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
    'azure', 
    UUID_TO_BIN(UUID()), 
    UUID_TO_BIN(UUID()),
    FALSE,
    TRUE,
    CURDATE(),
    'william.ohara@subscripify.com' );
    
    
    INSERT INTO tenant (
	tenant_uuid, 
    org_name,
    top_level_domain,
    secondary_domain,
    subdomain, 
    kube_namespace_prefix, 
    lord_services_config, 
    super_services_config,
    public_services_config,
    private_access_config,
    custom_access_config,
    subscripify_deployment_cloud_location,
    liege_uuid,
    lord_uuid,
    is_lord_tenant,
    is_super_tenant,
    create_timestamp,
    created_by
    )
VALUES
	(
	 UUID_TO_BIN(UUID()), 
    'Subscripify',
    'com',
    'subscripify',
    'super-tenant', 
    'super-tenant', 
    null, 
    UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'), 
    UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'), 
	UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
    UUID_TO_BIN('845c5fe4-3864-11ed-907f-f5001f9bae96'),
    'azure', 
    UUID_TO_BIN(UUID()), 
    UUID_TO_BIN(UUID()),
    FALSE,
    TRUE,
    CURDATE(),
    'william.ohara@subscripify.com' );
    
-- UPDATE tenant set isLordTenant = FALSE WHERE createdBy = 'william.ohara@subscripify.com';


    
SELECT *, BIN_TO_UUID(tenant_uuid) FROM tenant;



-- SELECT BIN_TO_UUID(org_id), org_name, subdomain, kube_namespace_prefix, subscription_type FROM organizations WHERE BIN_TO_UUID(org_id) = '67e1e031-2ef5-11ed-833b-6636daa5a961';