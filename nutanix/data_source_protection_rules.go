package nutanix

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func dataSourceNutanixProtectionRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixProtectionRulesRead,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_update_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"spec_hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"categories": categoriesSchema(),
						"owner_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type: schema.TypeString,
									},
									"uuid": {
										Type: schema.TypeString,
									},
									"name": {
										Type: schema.TypeString,
									},
								},
							},
						},
						"project_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type: schema.TypeString,
									},
									"uuid": {
										Type: schema.TypeString,
									},
									"name": {
										Type: schema.TypeString,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_connectivity_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_availability_zone_index": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"source_availability_zone_index": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshot_schedule_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"recovery_point_objective_secs": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"local_snapshot_retention_policy": {
													Type:     schema.TypeMap,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"num_snapshots": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"rollup_retention_policy_multiple": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"rollup_retention_policy_snapshot_interval_type": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"auto_suspend_timeout_secs": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"snapshot_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"remote_snapshot_retention_policy": {
													Type:     schema.TypeMap,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"num_snapshots": {
																Type: schema.TypeInt,
															},
															"rollup_retention_policy_multiple": {
																Type: schema.TypeInt,
															},
															"rollup_retention_policy_snapshot_interval_type": {
																Type: schema.TypeInt,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"ordered_availability_zone_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"availability_zone_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"category_filter": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"params": {
										Type:     schema.TypeSet,
										Computed: true,
										Set:      filterParamsHash,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"values": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixProtectionRulesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllProtectionRules()
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}
	if err := d.Set("entities", flattenProtectionRuleEntities(resp.Entities)); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenProtectionRuleEntities(ProtectionRules []*v3.ProtectionRuleResponse) []map[string]interface{} {
	entities := make([]map[string]interface{}, len(ProtectionRules))

	for i, protectionRule := range ProtectionRules {
		metadata, categories := setRSEntityMetadata(protectionRule.Metadata)

		entities[i] = map[string]interface{}{
			"name":                                protectionRule.Status.Name,
			"metadata":                            metadata,
			"categories":                          categories,
			"project_reference":                   flattenReferenceValues(protectionRule.Metadata.ProjectReference),
			"owner_reference":                     flattenReferenceValues(protectionRule.Metadata.OwnerReference),
			"start_time":                          protectionRule.Status.Resources.StartTime,
			"availability_zone_connectivity_list": flattenAvailabilityZoneConnectivityList(protectionRule.Spec.Resources.AvailabilityZoneConnectivityList),
			"ordered_availability_zone_list":      flattenOrderAvailibilityList(protectionRule.Spec.Resources.OrderedAvailabilityZoneList),
			"category_filter":                     flattenCategoriesFilter(protectionRule.Spec.Resources.CategoryFilter),
			"state":                               protectionRule.Status.State,
			"api_version":                         protectionRule.APIVersion,
		}
	}
	return entities
}
