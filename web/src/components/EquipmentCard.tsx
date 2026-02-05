import type { EquipmentItem, SellerInfo } from '../equipmentTypes';

interface EquipmentCardProps {
  item: EquipmentItem;
  seller?: SellerInfo;
  onAddToInventory: (item: EquipmentItem) => void;
  onViewDetails: (item: EquipmentItem) => void;
}

export function EquipmentCard({ item, seller, onAddToInventory, onViewDetails }: EquipmentCardProps) {
  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency || 'USD',
    }).format(price);
  };

  const formatSpecs = (specs?: Record<string, unknown>) => {
    if (!specs || Object.keys(specs).length === 0) return null;
    return Object.entries(specs)
      .slice(0, 3)
      .map(([key, value]) => (
        <span key={key} className="px-2 py-0.5 bg-slate-700 text-slate-300 rounded text-xs">
          {key}: {String(value)}
        </span>
      ));
  };

  return (
    <div className="bg-slate-800 border border-slate-700 rounded-xl p-4 hover:border-slate-600 transition-all group">
      <div className="flex gap-4">
        {/* Image */}
        <div className="flex-shrink-0 w-24 h-24 rounded-lg overflow-hidden bg-slate-700">
          {item.imageUrl ? (
            <img
              src={item.imageUrl}
              alt={item.name}
              className="w-full h-full object-cover"
              loading="lazy"
              onError={(e) => {
                (e.target as HTMLImageElement).style.display = 'none';
              }}
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-slate-500">
              <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
          )}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          {/* Header */}
          <div className="flex items-start justify-between gap-2 mb-1">
            <div className="flex items-center gap-2">
              <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                item.inStock
                  ? 'bg-green-500/20 text-green-400'
                  : 'bg-red-500/20 text-red-400'
              }`}>
                {item.inStock ? 'In Stock' : 'Out of Stock'}
              </span>
              <span className="text-xs text-slate-500">
                {seller?.name || item.seller}
              </span>
            </div>
            <div className="text-lg font-semibold text-primary-400">
              {formatPrice(item.price, item.currency)}
            </div>
          </div>

          {/* Title */}
          <h3 
            className="text-white font-medium mb-1 line-clamp-2 cursor-pointer hover:text-primary-400 transition-colors"
            onClick={() => onViewDetails(item)}
          >
            {item.name}
          </h3>

          {/* Manufacturer */}
          {item.manufacturer && (
            <p className="text-slate-400 text-sm mb-2">
              {item.manufacturer}
            </p>
          )}

          {/* Specs */}
          {item.keySpecs && (
            <div className="flex flex-wrap gap-1 mb-3">
              {formatSpecs(item.keySpecs as Record<string, unknown>)}
            </div>
          )}

          {/* Actions */}
          <div className="flex items-center gap-2">
            <button
              onClick={() => onAddToInventory(item)}
              className="px-3 py-1.5 bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium rounded-lg transition-colors flex items-center gap-1"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              Add to Inventory
            </button>
            <a
              href={item.productUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-slate-300 text-sm font-medium rounded-lg transition-colors flex items-center gap-1"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              View
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

interface EquipmentListProps {
  items: EquipmentItem[];
  sellers: SellerInfo[];
  isLoading: boolean;
  error: string | null;
  onAddToInventory: (item: EquipmentItem) => void;
  onViewDetails: (item: EquipmentItem) => void;
}

export function EquipmentList({ items, sellers, isLoading, error, onAddToInventory, onViewDetails }: EquipmentListProps) {
  const sellerMap = new Map(sellers.map(s => [s.id, s]));

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="text-center max-w-md">
          <div className="w-16 h-16 bg-red-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-white mb-2">Failed to Load Equipment</h3>
          <p className="text-slate-400 text-sm">{error}</p>
        </div>
      </div>
    );
  }

  if (isLoading && items.length === 0) {
    return (
      <div className="flex-1 p-6 space-y-4">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="bg-slate-800 border border-slate-700 rounded-xl p-4 animate-pulse">
            <div className="flex gap-4">
              <div className="w-24 h-24 bg-slate-700 rounded-lg" />
              <div className="flex-1 space-y-3">
                <div className="flex justify-between">
                  <div className="w-20 h-5 bg-slate-700 rounded" />
                  <div className="w-16 h-6 bg-slate-700 rounded" />
                </div>
                <div className="h-5 bg-slate-700 rounded w-3/4" />
                <div className="h-4 bg-slate-700 rounded w-1/4" />
                <div className="flex gap-1">
                  <div className="w-16 h-5 bg-slate-700 rounded" />
                  <div className="w-16 h-5 bg-slate-700 rounded" />
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="text-center max-w-md">
          <div className="w-16 h-16 bg-slate-700 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-white mb-2">No Equipment Found</h3>
          <p className="text-slate-400 text-sm">
            Try adjusting your search or filters to find equipment.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-y-auto p-6">
      <div className="space-y-4 max-w-4xl mx-auto">
        {items.map(item => (
          <EquipmentCard
            key={item.id}
            item={item}
            seller={sellerMap.get(item.sellerId)}
            onAddToInventory={onAddToInventory}
            onViewDetails={onViewDetails}
          />
        ))}
      </div>
      {isLoading && (
        <div className="flex items-center justify-center py-8">
          <div className="w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full animate-spin" />
        </div>
      )}
    </div>
  );
}
