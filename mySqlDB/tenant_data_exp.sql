select *, bin_to_uuid(tenant_uuid) as tenant_uuid_r from tenant;
 
INSERT INTO tenant (
			tenant_uuid, 
			tenant_alias,
			top_level_domain,
			secondary_domain,
			subscripify_deployment_cloud_location,
			subdomain, 
			kube_namespace_prefix,
			super_services_config, 
			public_services_config,
			private_access_config,
			custom_access_config,
            is_super_tenant,
            liege_uuid,
            lord_uuid,
		  created_by
				)
			VALUES (UUID_TO_BIN(UUID()),
			 'mybadtenant5',
			 (SELECT top_level_domain FROM (SELECT top_level_domain FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')) as b ), 
			 (SELECT secondary_domain FROM (SELECT secondary_domain FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')) as b ),  
			 (SELECT subscripify_deployment_cloud_location FROM (SELECT subscripify_deployment_cloud_location FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')) as b ),
			 'mybadtenant5', 
			 'mybadtenant5', 
			 UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
             UUID_TO_BIN('845c5fe4-3864-11ed-907f-f5001f9bae96'), 
			 true,
             UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
             uuid_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
			 'william.ohara@subscripify.com');
             
INSERT INTO tenant (
			tenant_uuid, 
			tenant_alias,
			top_level_domain,
			secondary_domain,
			subscripify_deployment_cloud_location,
			subdomain, 
			kube_namespace_prefix,
			super_services_config, 
			public_services_config,
			private_access_config,
			custom_access_config,
            is_super_tenant,
            liege_uuid,
            lord_uuid,
		  created_by
				)   
                select
                UUID_TO_BIN(UUID()),
			 'mybadtenant6',
             top_level_domain,
		     secondary_domain,
		     subscripify_cloud_deployment_location,
			 'mybadtenant6', 
			 'mybadtenant6', 
			 UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
             UUID_TO_BIN('845c5fe4-3864-11ed-907f-f5001f9bae96'), 
			 true,
             UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
             uuid_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
			 'william.ohara@subscripify.com'
                from tenants where tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432');
    
    
 INSERT INTO tenant (
			tenant_uuid, 
			tenant_alias,
			top_level_domain,
			secondary_domain,
			subscripify_deployment_cloud_location,
			subdomain, 
			kube_namespace_prefix,
			super_services_config, 
			public_services_config,
			private_access_config,
			custom_access_config,
            is_super_tenant,
            liege_uuid,
            lord_uuid,
		  created_by
				)
			VALUES (UUID_TO_BIN(UUID()),
			 'mybadtenant',
			 (SELECT top_level_domain FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')), 
			 (SELECT secondary_domain FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')),  
			 (SELECT subscripify_deployment_cloud_location FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')),
			 'mybadtenant', 
			 'mybadtenant', 
			 UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
             UUID_TO_BIN('845c5fe4-3864-11ed-907f-f5001f9bae96'), 
			 true,
             UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
             uuid_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
			 'william.ohara@subscripify.com');   

EXPLAIN INSERT INTO tenant (tenant_uuid, secondary_domain)
VALUES (
  UUID_TO_BIN(UUID()),
  (SELECT top_level_domain FROM (SELECT top_level_domain FROM tenant WHERE tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')) as b ));
  
explain insert iNTO tenant (tenant_uuid, secondary_domain)
VALUES (
  UUID_TO_BIN(UUID()),
  (SELECT b.secondary_domain FROM tenant a INNER JOIN tenant b on a.tenant_uuid = b.tenant_uuid WHERE b.tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432')));
-- SELECT BIN_TO_UUID(org_id), org_name, subdomain, kube_namespace_prefix, subscription_type FROM organizations WHERE BIN_TO_UUID(org_id) = '67e1e031-2ef5-11ed-833b-6636daa5a961';