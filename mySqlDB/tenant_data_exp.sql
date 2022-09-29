select bin_to_uuid(tenant_uuid) as tenant_uuid, tenant_alias,secondary_domain, subdomain, kube_namespace_prefix,
lord_services_config, super_services_config, public_services_config, private_access_config, custom_access_config, 
bin_to_uuid(liege_uuid), bin_to_uuid(lord_uuid), is_lord_tenant, is_super_tenant from tenant ORDER BY create_timestamp DESC;

-- DELETE FROM tenant WHERE tenant_uuid = UUID_TO_BIN('2f40a6e4-ede4-41f5-aeba-b8d34ca1d06a');
DELETE FROM tenant WHERE alias = "testing-123" AND is_lord_tenant IS NULL AND is_super_tenant = false;
SELECT count(tenant_uuid) as count FROM tenant WHERE tenant_uuid = UUID_TO_BIN('5c0ee1f6-6b3f-46c8-a62e-fca70d460bf0');

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
             
explain INSERT INTO tenant (
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
			 'mybadtenant8',
             top_level_domain,
		     secondary_domain,
		     subscripify_deployment_cloud_location,
			 'mybadtenant8', 
			 'mybadtenant8', 
			 UUID_TO_BIN('432ead01-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('6a24689f-385a-11ed-907f-f5001f9bae96'),
			 UUID_TO_BIN('c9057cff-3863-11ed-907f-f5001f9bae96'),
             UUID_TO_BIN('845c5fe4-3864-11ed-907f-f5001f9bae96'), 
			 true,
             UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
             uuid_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432'),
			 'william.ohara@subscripify.com'
                from tenant where tenant_uuid = UUID_to_bin('0904ba44-9d41-418a-bbd8-2c70506c3432');
    
    
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

SELECT 
	BIN_TO_UUID(tenant_uuid),
	 tenant_alias,
	 top_level_domain,
	 secondary_domain,
	 subdomain,
	 kube_namespace_prefix,
	 BIN_TO_UUID(lord_services_config),
	 BIN_TO_UUID(super_services_config),
	 BIN_TO_UUID(public_services_config),
	 BIN_TO_UUID(private_access_config),
	 BIN_TO_UUID(custom_access_config),
	 subscripify_deployment_cloud_location,
	 BIN_TO_UUID(liege_uuid),
	 BIN_TO_UUID(lord_uuid),
	 is_lord_tenant,
	 is_super_tenant,
	 create_timestamp,
	 created_by
	 FROM tenant WHERE tenant_uuid = UUID_TO_BIN('0904ba44-9d41-418a-bbd8-2c70506c343');
												  

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
					)   SELECT UUID_TO_BIN(uuid()), "super tenent", top_level_domain, secondary_domain, subscripify_deployment_cloud_location, 
                    "apicreated-super-tenant", "apicreated-super-tenant", UUID_TO_BIN("9a79e3b0-F899-f6ed-6Fd6-Fc4BE5A1E9fe"),
				 UUID_TO_BIN("44ED6Aa2-E2Cf-A5DF-6B0E-DDAED41fC5Da"), UUID_TO_BIN("Ae3cEafe-D53b-2AEc-f1C4-fd33F9f2E0df"), UUID_TO_BIN("0f9e7A1a-d7be-A460-8bEA-FAF1AFEfA9D0"), 
                 true, UUID_TO_BIN("0904ba44-9d41-418a-bbd8-2c70506c3432"), UUID_TO_BIN("0904ba44-9d41-418a-bbd8-2c70506c3432"), ? FROM tenant WHERE tenant_uuid = UUID_to_bin(?);