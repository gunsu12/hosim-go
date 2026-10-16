import type {
  PriceSchemaRuleRecord,
  PriceSchemaTreeNode,
  AdjustmentType
} from '$lib/types/finance/price_schema';
import type { DepartementRecord } from '$lib/types/master/departement';
import type { ServiceUnitRecord } from '$lib/types/master/service_unit';
import type { ItemSummaryRecord } from '$lib/types/master/item';

export function formatAdjustmentText(
  type: AdjustmentType,
  value: number | string
): string {
  const num = typeof value === 'string' ? parseFloat(value) || 0 : value || 0;
  switch (type) {
    case 'NONE':
      return 'Buku Tarif (0%)';
    case 'PERCENT':
      if (num > 0) return `Markup +${num}%`;
      if (num < 0) return `Diskon ${num}%`;
      return 'Tarif Normal (0%)';
    case 'FIXED_AMOUNT':
      if (num > 0) return `+Rp ${num.toLocaleString('id-ID')}`;
      if (num < 0) return `-Rp ${Math.abs(num).toLocaleString('id-ID')}`;
      return 'Rp 0';
    case 'FIXED_PRICE':
      return `Rp ${num.toLocaleString('id-ID')} (Harga Tetap)`;
    default:
      return String(type);
  }
}

export function buildPriceSchemaTree(
  defaultRule: PriceSchemaRuleRecord | null,
  rules: PriceSchemaRuleRecord[],
  departments: DepartementRecord[],
  serviceUnits: ServiceUnitRecord[],
  itemsMap: Map<string, ItemSummaryRecord>
): PriceSchemaTreeNode[] {
  // Index rules by scope keys
  const deptRulesMap = new Map<string, PriceSchemaRuleRecord>();
  const unitRulesMap = new Map<string, PriceSchemaRuleRecord>();
  const itemUnitRules: PriceSchemaRuleRecord[] = [];
  const itemDeptRules: PriceSchemaRuleRecord[] = [];
  const globalItemRules: PriceSchemaRuleRecord[] = [];

  for (const r of rules) {
    if (r.level === 'SCHEMA_DEFAULT') {
      // handled separately
    } else if (r.level === 'DEPT' && r.departement_id) {
      deptRulesMap.set(r.departement_id, r);
    } else if (r.level === 'UNIT' && r.service_unit_id) {
      unitRulesMap.set(r.service_unit_id, r);
    } else if (r.level === 'ITEM_UNIT') {
      itemUnitRules.push(r);
    } else if (r.level === 'ITEM_DEPT') {
      itemDeptRules.push(r);
    } else if (r.level === 'ITEM') {
      globalItemRules.push(r);
    }
  }

  // Fallback defaults if defaultRule is missing
  const fallbackType: AdjustmentType = defaultRule?.tariff_adjustment_type || 'NONE';
  const fallbackVal = Number(defaultRule?.tariff_adjustment_value) || 0;
  const fallbackMedType: AdjustmentType = defaultRule?.medication_adjustment_type || 'NONE';
  const fallbackMedVal = Number(defaultRule?.medication_adjustment_value) || 0;
  const fallbackPatientPct = Number(defaultRule?.pay_by_patient_percent) || 0;
  const fallbackCustomerPct = Number(defaultRule?.pay_by_customer_percent) || 100;
  const fallbackTax = Number(defaultRule?.tax_percent) || 0;

  // Group service units by department
  const unitsByDept = new Map<string, ServiceUnitRecord[]>();
  for (const u of serviceUnits) {
    const dId = u.departement_id || 'UNASSIGNED';
    if (!unitsByDept.has(dId)) {
      unitsByDept.set(dId, []);
    }
    unitsByDept.get(dId)!.push(u);
  }

  const resultTree: PriceSchemaTreeNode[] = [];

  // 1. Build Department Branches
  for (const dept of departments) {
    const deptRule = deptRulesMap.get(dept.id);
    const hasDeptRule = !!deptRule;

    const deptTariffType = hasDeptRule ? deptRule.tariff_adjustment_type : fallbackType;
    const deptTariffVal = hasDeptRule ? Number(deptRule.tariff_adjustment_value) : fallbackVal;
    const deptMedType = hasDeptRule ? deptRule.medication_adjustment_type || 'NONE' : fallbackMedType;
    const deptMedVal = hasDeptRule ? Number(deptRule.medication_adjustment_value) : fallbackMedVal;
    const deptPatient = hasDeptRule ? Number(deptRule.pay_by_patient_percent) : fallbackPatientPct;
    const deptCustomer = hasDeptRule ? Number(deptRule.pay_by_customer_percent) : fallbackCustomerPct;
    const deptTax = hasDeptRule ? Number(deptRule.tax_percent) : fallbackTax;

    const deptChildren: PriceSchemaTreeNode[] = [];

    // Service units under this department
    const deptUnits = unitsByDept.get(dept.id) || [];
    for (const unit of deptUnits) {
      const unitRule = unitRulesMap.get(unit.id);
      const hasUnitRule = !!unitRule;

      const unitTariffType = hasUnitRule ? unitRule.tariff_adjustment_type : deptTariffType;
      const unitTariffVal = hasUnitRule ? Number(unitRule.tariff_adjustment_value) : deptTariffVal;
      const unitMedType = hasUnitRule ? unitRule.medication_adjustment_type || 'NONE' : deptMedType;
      const unitMedVal = hasUnitRule ? Number(unitRule.medication_adjustment_value) : deptMedVal;
      const unitPatient = hasUnitRule ? Number(unitRule.pay_by_patient_percent) : deptPatient;
      const unitCustomer = hasUnitRule ? Number(unitRule.pay_by_customer_percent) : deptCustomer;
      const unitTax = hasUnitRule ? Number(unitRule.tax_percent) : deptTax;

      // Items under this service unit (ITEM_UNIT)
      const unitItems = itemUnitRules.filter(r => r.service_unit_id === unit.id);
      const unitItemChildren: PriceSchemaTreeNode[] = unitItems.map(ir => {
        const itemObj = ir.item_id ? itemsMap.get(ir.item_id) : undefined;
        return {
          id: ir.id,
          nodeType: 'ITEM',
          scopeType: 'ITEM_UNIT',
          name: itemObj ? itemObj.name : `Item Tindakan (${ir.item_id || 'Unknown'})`,
          code: itemObj?.code,
          serviceUnitId: unit.id,
          itemId: ir.item_id,
          isInherited: false,
          ruleId: ir.id,
          tariffAdjustmentType: ir.tariff_adjustment_type,
          tariffAdjustmentValue: Number(ir.tariff_adjustment_value),
          medicationAdjustmentType: ir.medication_adjustment_type || 'NONE',
          medicationAdjustmentValue: Number(ir.medication_adjustment_value),
          payByPatientPercent: Number(ir.pay_by_patient_percent),
          payByCustomerPercent: Number(ir.pay_by_customer_percent),
          taxPercent: Number(ir.tax_percent)
        };
      });

      deptChildren.push({
        id: hasUnitRule ? unitRule.id : `synthetic_unit_${unit.id}`,
        nodeType: 'UNIT',
        scopeType: 'UNIT',
        name: unit.name,
        code: unit.code,
        departementId: dept.id,
        serviceUnitId: unit.id,
        isInherited: !hasUnitRule,
        ruleId: hasUnitRule ? unitRule.id : undefined,
        parentRuleSummary: !hasUnitRule
          ? `Mewarisi aturan ${dept.name} (${formatAdjustmentText(deptTariffType, deptTariffVal)})`
          : undefined,
        tariffAdjustmentType: unitTariffType,
        tariffAdjustmentValue: unitTariffVal,
        medicationAdjustmentType: unitMedType,
        medicationAdjustmentValue: unitMedVal,
        payByPatientPercent: unitPatient,
        payByCustomerPercent: unitCustomer,
        taxPercent: unitTax,
        children: unitItemChildren
      });
    }

    // Also check for ITEM_DEPT rules under this department
    const deptItems = itemDeptRules.filter(r => r.departement_id === dept.id);
    for (const ir of deptItems) {
      const itemObj = ir.item_id ? itemsMap.get(ir.item_id) : undefined;
      deptChildren.push({
        id: ir.id,
        nodeType: 'ITEM',
        scopeType: 'ITEM_DEPT',
        name: itemObj ? itemObj.name : `Item Spesifik Dept (${ir.item_id || 'Unknown'})`,
        code: itemObj?.code,
        departementId: dept.id,
        itemId: ir.item_id,
        isInherited: false,
        ruleId: ir.id,
        tariffAdjustmentType: ir.tariff_adjustment_type,
        tariffAdjustmentValue: Number(ir.tariff_adjustment_value),
        medicationAdjustmentType: ir.medication_adjustment_type || 'NONE',
        medicationAdjustmentValue: Number(ir.medication_adjustment_value),
        payByPatientPercent: Number(ir.pay_by_patient_percent),
        payByCustomerPercent: Number(ir.pay_by_customer_percent),
        taxPercent: Number(ir.tax_percent)
      });
    }

    resultTree.push({
      id: hasDeptRule ? deptRule.id : `synthetic_dept_${dept.id}`,
      nodeType: 'DEPT',
      scopeType: 'DEPT',
      name: dept.name,
      code: dept.code,
      departementId: dept.id,
      isInherited: !hasDeptRule,
      ruleId: hasDeptRule ? deptRule.id : undefined,
      parentRuleSummary: !hasDeptRule
        ? `Mewarisi Default Skema (${formatAdjustmentText(fallbackType, fallbackVal)})`
        : undefined,
      tariffAdjustmentType: deptTariffType,
      tariffAdjustmentValue: deptTariffVal,
      medicationAdjustmentType: deptMedType,
      medicationAdjustmentValue: deptMedVal,
      payByPatientPercent: deptPatient,
      payByCustomerPercent: deptCustomer,
      taxPercent: deptTax,
      children: deptChildren
    });
  }

  // 2. Global Item Branch (if any ITEM rules exist at global level)
  if (globalItemRules.length > 0) {
    const globalItemChildren: PriceSchemaTreeNode[] = globalItemRules.map(ir => {
      const itemObj = ir.item_id ? itemsMap.get(ir.item_id) : undefined;
      return {
        id: ir.id,
        nodeType: 'ITEM',
        scopeType: 'ITEM',
        name: itemObj ? itemObj.name : `Item Tindakan (${ir.item_id || 'Unknown'})`,
        code: itemObj?.code,
        itemId: ir.item_id,
        isInherited: false,
        ruleId: ir.id,
        tariffAdjustmentType: ir.tariff_adjustment_type,
        tariffAdjustmentValue: Number(ir.tariff_adjustment_value),
        medicationAdjustmentType: ir.medication_adjustment_type || 'NONE',
        medicationAdjustmentValue: Number(ir.medication_adjustment_value),
        payByPatientPercent: Number(ir.pay_by_patient_percent),
        payByCustomerPercent: Number(ir.pay_by_customer_percent),
        taxPercent: Number(ir.tax_percent)
      };
    });

    resultTree.push({
      id: 'synthetic_global_items_group',
      nodeType: 'DEPT',
      scopeType: 'ITEM',
      name: 'Pengecualian Tindakan Global (Lintas Unit)',
      code: 'GLOBAL',
      isInherited: false,
      tariffAdjustmentType: fallbackType,
      tariffAdjustmentValue: fallbackVal,
      medicationAdjustmentType: fallbackMedType,
      medicationAdjustmentValue: fallbackMedVal,
      payByPatientPercent: fallbackPatientPct,
      payByCustomerPercent: fallbackCustomerPct,
      taxPercent: fallbackTax,
      children: globalItemChildren
    });
  }

  return resultTree;
}

export function filterTreeRecursive(
  nodes: PriceSchemaTreeNode[],
  searchQuery: string,
  filterMode: 'ALL' | 'OVERRIDE' | 'INHERITED'
): PriceSchemaTreeNode[] {
  const q = searchQuery.trim().toLowerCase();

  return nodes
    .map((node): PriceSchemaTreeNode | null => {
      // Recursively filter children first
      const filteredChildren = node.children
        ? filterTreeRecursive(node.children, searchQuery, filterMode)
        : [];

      // Check self match
      const nameMatches = !q || node.name.toLowerCase().includes(q) || (node.code && node.code.toLowerCase().includes(q));

      let modeMatches = true;
      if (filterMode === 'OVERRIDE') {
        modeMatches = !node.isInherited;
      } else if (filterMode === 'INHERITED') {
        modeMatches = node.isInherited;
      }

      const selfMatches = nameMatches && modeMatches;
      const hasMatchingChildren = filteredChildren.length > 0;

      if (selfMatches || hasMatchingChildren) {
        return {
          ...node,
          children: filteredChildren
        };
      }
      return null;
    })
    .filter((n): n is PriceSchemaTreeNode => n !== null);
}

export function collectExpandedNodeIds(
  nodes: PriceSchemaTreeNode[],
  maxDepth: number = 2,
  currentDepth: number = 0
): Set<string> {
  const ids = new Set<string>();

  function traverse(n: PriceSchemaTreeNode, depth: number) {
    if (depth < maxDepth) {
      ids.add(n.id);
    }
    if (n.children) {
      for (const child of n.children) {
        traverse(child, depth + 1);
      }
    }
  }

  for (const root of nodes) {
    traverse(root, currentDepth);
  }

  return ids;
}
