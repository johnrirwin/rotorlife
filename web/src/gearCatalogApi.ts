// Gear Catalog API client for crowd-sourced gear definitions

import type {
  GearCatalogItem,
  GearCatalogSearchParams,
  GearCatalogSearchResponse,
  GearCatalogCreateResponse,
  CreateGearCatalogParams,
  NearMatchParams,
  NearMatchResponse,
  GearType,
} from './gearCatalogTypes';

const API_BASE = import.meta.env.VITE_API_BASE_URL || '';

// Get access token from localStorage
function getAccessToken(): string | null {
  return localStorage.getItem('access_token');
}

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = getAccessToken();
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options?.headers,
  };
  
  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
  }
  
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }));
    throw new Error(error.message || error.error || `HTTP ${response.status}`);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}

/**
 * Search the gear catalog with optional filters
 */
export async function searchGearCatalog(params: GearCatalogSearchParams): Promise<GearCatalogSearchResponse> {
  const searchParams = new URLSearchParams();
  
  if (params.query) searchParams.set('q', params.query);
  if (params.gearType) searchParams.set('gearType', params.gearType);
  if (params.brand) searchParams.set('brand', params.brand);
  if (params.status) searchParams.set('status', params.status);
  if (params.limit) searchParams.set('limit', params.limit.toString());
  if (params.offset) searchParams.set('offset', params.offset.toString());

  const query = searchParams.toString();
  return fetchAPI<GearCatalogSearchResponse>(`/api/gear-catalog/search${query ? `?${query}` : ''}`);
}

/**
 * Get popular gear items, optionally filtered by type
 */
export async function getPopularGear(gearType?: GearType, limit?: number): Promise<{ items: GearCatalogItem[] }> {
  const searchParams = new URLSearchParams();
  
  if (gearType) searchParams.set('gearType', gearType);
  if (limit) searchParams.set('limit', limit.toString());

  const query = searchParams.toString();
  return fetchAPI<{ items: GearCatalogItem[] }>(`/api/gear-catalog/popular${query ? `?${query}` : ''}`);
}

/**
 * Get a single catalog item by ID
 */
export async function getGearCatalogItem(id: string): Promise<GearCatalogItem> {
  return fetchAPI<GearCatalogItem>(`/api/gear-catalog/${id}`);
}

/**
 * Create a new catalog item (or return existing if duplicate detected)
 * Returns { item, existing } where existing=true if we found a match
 */
export async function createGearCatalogItem(params: CreateGearCatalogParams): Promise<GearCatalogCreateResponse> {
  return fetchAPI<GearCatalogCreateResponse>('/api/gear-catalog', {
    method: 'POST',
    body: JSON.stringify(params),
  });
}

/**
 * Find near matches for potential duplicate detection before creating
 */
export async function findNearMatches(params: NearMatchParams): Promise<NearMatchResponse> {
  return fetchAPI<NearMatchResponse>('/api/gear-catalog/near-matches', {
    method: 'POST',
    body: JSON.stringify(params),
  });
}

/**
 * Flag a catalog item for review (e.g., incorrect data, duplicate, spam)
 */
export async function flagGearCatalogItem(id: string, reason: string): Promise<{ status: string }> {
  return fetchAPI<{ status: string }>(`/api/gear-catalog/${id}/flag`, {
    method: 'POST',
    body: JSON.stringify({ reason }),
  });
}

/**
 * Create or get a catalog item, then link it to inventory
 * This is a convenience function that:
 * 1. Creates/finds the catalog item
 * 2. Returns the catalog ID to use when adding to inventory
 */
export async function getOrCreateCatalogItem(params: CreateGearCatalogParams): Promise<GearCatalogItem> {
  const response = await createGearCatalogItem(params);
  return response.item;
}

/**
 * Typeahead search - debounced search for autocomplete UI
 * Returns a smaller subset of results for quick display
 */
export async function typeaheadSearch(
  query: string,
  gearType?: GearType,
  limit = 8
): Promise<GearCatalogItem[]> {
  if (!query || query.length < 2) {
    return [];
  }

  const response = await searchGearCatalog({
    query,
    gearType,
    limit,
    status: 'active',
  });

  return response.items;
}
