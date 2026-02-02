// Social/Pilot Directory types

import type { AircraftType } from './aircraftTypes';

// Profile visibility settings
export type ProfileVisibility = 'public' | 'private';

// Social settings for a user
export interface SocialSettings {
  profileVisibility: ProfileVisibility;
  showAircraft: boolean;
  allowSearch: boolean;
}

// Pilot search result (from /api/pilots/search)
export interface PilotSearchResult {
  id: string;
  callSign?: string;
  displayName?: string;
  googleName?: string;
  effectiveAvatarUrl: string;
}

// Pilot search response
export interface PilotSearchResponse {
  pilots: PilotSearchResult[];
  total: number;
}

// Sanitized ELRS settings (safe for public display)
// CRITICAL: This type intentionally OMITS sensitive fields like bindPhrase, modelMatch, uid
export interface ELRSSanitizedSettings {
  receiverModel?: string;    // e.g., "EP1", "RP1", "RP3"
  packetRate?: string;       // e.g., "250Hz", "500Hz"
  telemetryRatio?: string;   // e.g., "1:128", "1:64"
  switchMode?: string;       // e.g., "Hybrid", "Wide"
  outputPower?: string;      // e.g., "250mW", "500mW", "Dynamic"
  regulatoryDomain?: string; // e.g., "FCC", "LBT"
  firmwareVersion?: string;  // e.g., "3.4.0"
  rxProtocol?: string;       // Protocol type if applicable
}

// Component category types
export type ComponentCategory = 
  | 'fc' 
  | 'esc' 
  | 'elrs_module' 
  | 'vtx' 
  | 'motors' 
  | 'camera' 
  | 'frame' 
  | 'propellers' 
  | 'antenna';

// Public component info (no purchase details)
// NOTE: Purchase price, seller, and notes are intentionally omitted for privacy
export interface AircraftComponentPublic {
  category: ComponentCategory;
  name?: string;
  manufacturer?: string;
  imageUrl?: string;
}

// Public aircraft info shown on pilot profiles
export interface AircraftPublic {
  id: string;
  name: string;
  nickname?: string;
  type?: AircraftType;
  hasImage: boolean;
  description?: string;
  createdAt: string;
  components?: AircraftComponentPublic[];
  elrsSettings?: ELRSSanitizedSettings; // Sanitized ELRS data
}

// Full pilot profile (from /api/pilots/:id)
export interface PilotProfile {
  id: string;
  callSign?: string;
  displayName?: string;
  googleName?: string;
  effectiveAvatarUrl: string;
  createdAt: string;
  aircraft: AircraftPublic[];
  isFollowing: boolean;
  followerCount: number;
  followingCount: number;
}

// Pilot summary (for follower/following lists)
export interface PilotSummary {
  id: string;
  callSign?: string;
  displayName?: string;
  effectiveAvatarUrl: string;
}

// Follow list response
export interface FollowListResponse {
  pilots: PilotSummary[];
  totalCount: number;
}

// Follow/unfollow response
export interface FollowResponse {
  success: boolean;
  following: boolean;
  followId?: string;
}

// Avatar upload response
export interface AvatarUploadResponse {
  effectiveAvatarUrl: string;
  avatarType: 'google' | 'custom';
  customAvatarUrl: string;
}

