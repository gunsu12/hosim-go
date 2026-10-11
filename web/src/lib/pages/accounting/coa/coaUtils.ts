import type { AccountRecord, AccountTreeNode, AccountType } from '$lib/types/accounting/account';

export function filterTreeRecursive(
  nodes: AccountTreeNode[],
  searchQuery: string,
  selectedTypeFilter: AccountType | 'ALL',
  activeOnlyFilter: boolean,
  postableFilter: 'ALL' | 'POSTABLE' | 'HEADER'
): AccountTreeNode[] {
  if (selectedTypeFilter === 'ALL' && !searchQuery.trim() && !activeOnlyFilter && postableFilter === 'ALL') {
    return nodes;
  }

  const query = searchQuery.trim().toLowerCase();

  function matchNode(node: AccountTreeNode): boolean {
    const matchSearch = !query ||
      node.code.toLowerCase().includes(query) ||
      node.name.toLowerCase().includes(query);
    const matchType = selectedTypeFilter === 'ALL' || node.type === selectedTypeFilter;
    const matchActive = !activeOnlyFilter || node.is_active;
    const matchPostable = postableFilter === 'ALL' ||
      (postableFilter === 'POSTABLE' && node.is_postable) ||
      (postableFilter === 'HEADER' && !node.is_postable);

    return matchSearch && matchType && matchActive && matchPostable;
  }

  function recurse(list: AccountTreeNode[]): AccountTreeNode[] {
    const result: AccountTreeNode[] = [];
    for (const node of list) {
      const filteredChildren = node.children && node.children.length > 0
        ? recurse(node.children)
        : [];
      const isSelfMatch = matchNode(node);

      if (isSelfMatch || filteredChildren.length > 0) {
        result.push({
          ...node,
          children: filteredChildren
        });
      }
    }
    return result;
  }

  return recurse(nodes);
}

export function filterFlatAccounts(
  accounts: AccountRecord[],
  searchQuery: string,
  selectedTypeFilter: AccountType | 'ALL',
  activeOnlyFilter: boolean,
  postableFilter: 'ALL' | 'POSTABLE' | 'HEADER'
): AccountRecord[] {
  const query = searchQuery.trim().toLowerCase();

  return accounts.filter(acc => {
    const matchSearch = !query ||
      acc.code.toLowerCase().includes(query) ||
      acc.name.toLowerCase().includes(query) ||
      (acc.bank_name && acc.bank_name.toLowerCase().includes(query)) ||
      (acc.bank_account_number && acc.bank_account_number.toLowerCase().includes(query));
    const matchType = selectedTypeFilter === 'ALL' || acc.type === selectedTypeFilter;
    const matchActive = !activeOnlyFilter || acc.is_active;
    const matchPostable = postableFilter === 'ALL' ||
      (postableFilter === 'POSTABLE' && acc.is_postable) ||
      (postableFilter === 'HEADER' && !acc.is_postable);

    return matchSearch && matchType && matchActive && matchPostable;
  });
}

export function collectExpandedNodeIds(nodes: AccountTreeNode[], maxLevel?: number): Set<string> {
  const set = new Set<string>();
  function traverse(list: AccountTreeNode[]) {
    for (const n of list) {
      if (maxLevel === undefined || n.account_level <= maxLevel) {
        set.add(n.id);
      }
      if (n.children && n.children.length > 0) {
        traverse(n.children);
      }
    }
  }
  traverse(nodes);
  return set;
}
