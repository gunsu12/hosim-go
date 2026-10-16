// TypeScript Definitions for Finance Price Schema (HOSIM-GO)

export type AdjustmentType = 'NONE' | 'PERCENT' | 'FIXED_AMOUNT' | 'FIXED_PRICE';

export type MatchLevel =
  | 'SCHEMA_DEFAULT'
  | 'DEPT'
  | 'UNIT'
  | 'ITEM'
  | 'ITEM_DEPT'
  | 'ITEM_UNIT';

export interface PriceSchemaRuleRecord {
  id: string;
  price_schema_id: string;
  level: MatchLevel;
  departement_id?: string | null;
  service_unit_id?: string | null;
  item_id?: string | null;
  tariff_adjustment_type: AdjustmentType;
  tariff_adjustment_value: number | string;
  medication_adjustment_type?: AdjustmentType;
  medication_adjustment_value?: number | string;
  pay_by_patient_percent: number | string;
  pay_by_customer_percent: number | string;
  tax_percent: number | string;
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface PriceSchemaRecord {
  id: string;
  code: string;
  name: string;
  description?: string | null;
  base_price_plan_id?: string | null;
  base_price_plan_code?: string | null;
  base_price_plan_name?: string | null;
  is_active: boolean;
  rules_count?: number;
  rules?: PriceSchemaRuleRecord[];
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface CreatePriceSchemaDTO {
  code: string;
  name: string;
  description?: string;
  base_price_plan_id?: string | null;
  rules?: UpsertRuleDTO[];
}

export interface UpdatePriceSchemaDTO {
  name: string;
  description?: string;
  base_price_plan_id?: string | null;
  is_active?: boolean;
}

export interface UpsertRuleDTO {
  departement_id?: string | null;
  service_unit_id?: string | null;
  item_id?: string | null;
  tariff_adjustment_type: AdjustmentType;
  tariff_adjustment_value: number;
  medication_adjustment_type?: AdjustmentType;
  medication_adjustment_value?: number;
  pay_by_patient_percent: number;
  pay_by_customer_percent: number;
  tax_percent: number;
}

export interface SimulateTariffDTO {
  price_schema_id?: string;
  price_plan_id?: string;
  tariff_class_id: string;
  item_id: string;
  departement_id?: string;
  service_unit_id?: string;
  service_date?: string;
  is_cito?: boolean;
}

export interface SimulatedTariffComponent {
  component_id: string;
  component_code?: string;
  component_type?: string;
  component_name: string;
  base_amount?: number | string;
  amount?: number | string;
  final_amount?: number | string;
  coa_code: string;
}

export interface SimulateTariffResult {
  price_plan_id: string;
  price_plan_code: string;
  base_price: number | string;
  cito_surcharge: number | string;
  is_cito: boolean;
  price_schema_id?: string;
  price_schema_rule_id?: string;
  match_level?: MatchLevel;
  tariff_adjustment_type: AdjustmentType;
  tariff_adjustment_value: number | string;
  final_price: number | string;
  tax_percent: number | string;
  tax_amount: number | string;
  gross_amount: number | string;
  pay_by_patient_percent: number | string;
  pay_by_customer_percent: number | string;
  patient_amount: number | string;
  customer_amount: number | string;
  components: SimulatedTariffComponent[];
}

// Tree View Node Interface for Hierarchical Tree (ala CoaTreeView)
export interface PriceSchemaTreeNode {
  id: string; // Real rule ID if override exists, or synthetic ID (e.g., 'dept_123', 'unit_456')
  nodeType: 'DEFAULT' | 'DEPT' | 'UNIT' | 'ITEM';
  scopeType: MatchLevel;
  name: string;
  code?: string;
  departementId?: string | null;
  serviceUnitId?: string | null;
  itemId?: string | null;
  isInherited: boolean; // True if this node inherits rule from parent
  parentRuleSummary?: string;

  // Actual or effective rule values
  ruleId?: string; // Set only if this node has an explicit rule in database
  tariffAdjustmentType: AdjustmentType;
  tariffAdjustmentValue: number;
  medicationAdjustmentType: AdjustmentType;
  medicationAdjustmentValue: number;
  payByPatientPercent: number;
  payByCustomerPercent: number;
  taxPercent: number;

  children?: PriceSchemaTreeNode[];
}
