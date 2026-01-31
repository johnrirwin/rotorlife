import type { EquipmentCategory, EquipmentSearchParams, SellerInfo, InventorySummary } from '../equipmentTypes';
import { EQUIPMENT_CATEGORIES, ITEM_CONDITIONS } from '../equipmentTypes';
import type { AppSection } from '../equipmentTypes';

interface EquipmentSidebarProps {
  activeSection: AppSection;
  onSectionChange: (section: AppSection) => void;
  searchParams: EquipmentSearchParams;
  onSearchChange: (params: Partial<EquipmentSearchParams>) => void;
  sellers: SellerInfo[];
  inventorySummary: InventorySummary | null;
  inventoryCategory: EquipmentCategory | null;
  inventoryCondition: string | null;
  onInventoryFilterChange: (category: EquipmentCategory | null, condition: string | null) => void;
}

export function EquipmentSidebar({
  activeSection,
  onSectionChange,
  searchParams,
  onSearchChange,
  inventorySummary,
  inventoryCategory,
  inventoryCondition,
  onInventoryFilterChange,
}: EquipmentSidebarProps) {
  const handleCategorySelect = (category: EquipmentCategory | undefined) => {
    if (activeSection === 'equipment') {
      onSearchChange({ category });
    } else {
      onInventoryFilterChange(category || null, inventoryCondition);
    }
  };

  const selectedCategory = activeSection === 'equipment' 
    ? searchParams.category 
    : inventoryCategory;

  return (
    <aside className="w-64 flex-shrink-0 bg-slate-900 border-r border-slate-800 flex flex-col overflow-hidden">
      {/* Section Switcher */}
      <div className="p-4 border-b border-slate-800">
        <nav className="flex flex-col gap-1">
          <button
            onClick={() => onSectionChange('news')}
            className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-colors ${
              activeSection === 'news' 
                ? 'bg-primary-600/20 text-primary-400' 
                : 'text-slate-400 hover:text-white hover:bg-slate-800'
            }`}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
            </svg>
            <span className="font-medium">News Feed</span>
          </button>
          <button
            onClick={() => onSectionChange('equipment')}
            className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-colors ${
              activeSection === 'equipment' 
                ? 'bg-primary-600/20 text-primary-400' 
                : 'text-slate-400 hover:text-white hover:bg-slate-800'
            }`}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
            <span className="font-medium">Shop</span>
          </button>
          <button
            onClick={() => onSectionChange('inventory')}
            className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-colors ${
              activeSection === 'inventory' 
                ? 'bg-primary-600/20 text-primary-400' 
                : 'text-slate-400 hover:text-white hover:bg-slate-800'
            }`}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
            <span className="font-medium">My Gear</span>
            {inventorySummary && inventorySummary.totalItems > 0 && (
              <span className="ml-auto px-2 py-0.5 bg-slate-700 rounded-full text-xs text-slate-300">
                {inventorySummary.totalItems}
              </span>
            )}
          </button>
        </nav>
      </div>

      {/* Inventory filters (only for inventory section) */}
      {activeSection === 'inventory' && (
        <div className="p-4 border-b border-slate-800">
          <div className="space-y-3">
            {/* Condition filter */}
            <div>
              <label className="block text-xs font-medium text-slate-500 uppercase mb-1.5">
                Condition
              </label>
              <select
                value={inventoryCondition || ''}
                onChange={(e) => onInventoryFilterChange(inventoryCategory, e.target.value || null)}
                className="w-full px-3 py-1.5 bg-slate-800 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-primary-500"
              >
                <option value="">All Conditions</option>
                {ITEM_CONDITIONS.map(cond => (
                  <option key={cond.value} value={cond.value}>{cond.label}</option>
                ))}
              </select>
            </div>

            {/* Inventory summary */}
            {inventorySummary && (
              <div className="pt-2 border-t border-slate-800">
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div className="bg-slate-800 rounded-lg p-2">
                    <div className="text-slate-400 text-xs">Total Items</div>
                    <div className="text-white font-semibold">{inventorySummary.totalItems}</div>
                  </div>
                  <div className="bg-slate-800 rounded-lg p-2">
                    <div className="text-slate-400 text-xs">Total Value</div>
                    <div className="text-primary-400 font-semibold">
                      ${inventorySummary.totalValue.toFixed(0)}
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Categories (only for inventory section) */}
      {activeSection === 'inventory' && (
        <div className="flex-1 overflow-y-auto">
          <div className="p-4">
            <h3 className="text-xs font-medium text-slate-500 uppercase mb-2">
              Categories
            </h3>
            <nav className="space-y-0.5">
              <button
                onClick={() => handleCategorySelect(undefined)}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left text-sm transition-colors ${
                  !selectedCategory
                    ? 'bg-slate-800 text-white'
                    : 'text-slate-400 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                <span className="w-1.5 h-1.5 rounded-full bg-slate-500" />
                All Categories
              </button>
              {EQUIPMENT_CATEGORIES.map(cat => (
                <button
                  key={cat.value}
                  onClick={() => handleCategorySelect(cat.value)}
                  className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left text-sm transition-colors ${
                    selectedCategory === cat.value
                      ? 'bg-slate-800 text-white'
                      : 'text-slate-400 hover:text-white hover:bg-slate-800/50'
                  }`}
                >
                  <span className="w-1.5 h-1.5 rounded-full bg-primary-500" />
                  {cat.label}
                </button>
              ))}
            </nav>
          </div>
        </div>
      )}
    </aside>
  );
}
