/**
 * types.ts - TypeScript Type Definitions
 *
 * This file defines the "shape" of data we receive from and send to the API.
 * TypeScript uses these definitions to:
 * 1. Catch errors at compile time (before the code runs)
 * 2. Provide autocomplete in your editor
 * 3. Document what fields are available on each object
 *
 * Key concepts:
 * - interface: Defines the shape of an object (what properties it has)
 * - export: Makes the type available to import in other files
 * - ?: Optional property (might be undefined)
 * - |: Union type - can be one of multiple types (e.g., 'admin' | 'user')
 *
 * These types MUST match what the Go backend sends/expects.
 * If the backend changes, these need to be updated too.
 */

// ============================================================================
// Core User & Authentication Types
// ============================================================================

/**
 * Represents a user in the system.
 * This is returned by /api/auth/me and included in login/register responses.
 */
export interface User {
  id: number;
  email: string;
  name: string;
  system_role: 'admin' | 'user';  // Admins have extra permissions
  has_password: boolean;           // False for SSO-only users
  created_at: string;              // ISO 8601 date string
}

/**
 * Response from login/register endpoints.
 * Contains the JWT token and user info.
 */
export interface AuthResponse {
  token: string;  // JWT token - store this and send with future requests
  user: User;
}

// ============================================================================
// Groups & Membership Types
// ============================================================================

/**
 * A group is a collection of links that can be shared between users.
 * Users can be members of multiple groups.
 */
export interface Group {
  id: number;
  name: string;
  created_at: string;
}

/**
 * Represents a user's membership in a group.
 * The role determines what the user can do in that group.
 */
export interface GroupMembership {
  id: number;
  user_id: number;
  group_id: number;
  role: 'admin' | 'member';  // Admins can manage members, members can only view/add links
  user?: User;    // Sometimes included for display purposes
  group?: Group;  // Sometimes included for display purposes
}

// ============================================================================
// Links & Tags Types
// ============================================================================

/**
 * A tag is a label that can be applied to links for organization.
 */
export interface Tag {
  id: number;
  name: string;
  link_count?: number;  // How many links have this tag (only in some responses)
}

/**
 * A shortened link - the core entity of the app.
 */
export interface Link {
  id: number;
  group_id: number;       // Which group this link belongs to
  created_by_id: number;  // User who created it
  slug: string;           // The short code (e.g., "abc123" in shorty.io/abc123)
  url: string;            // The full URL this redirects to
  title: string;          // Display title
  description: string;    // Optional description/notes
  is_public: boolean;     // Can non-members access this link?
  is_unread: boolean;     // "Read later" flag
  click_count: number;    // How many times the link has been visited
  created_at: string;
  updated_at: string;
  tags?: Tag[];           // Tags applied to this link
  group?: Group;          // The group this belongs to (sometimes included)
}

// ============================================================================
// API Keys Types
// ============================================================================

/**
 * An API key allows programmatic access to the API without a password.
 * The full key is only shown once when created - we only store a prefix.
 */
export interface APIKey {
  id: number;
  key_prefix: string;      // First few characters for identification
  description: string;     // User-provided description
  last_used_at: string | null;  // When the key was last used (null if never)
  created_at: string;
}

/**
 * Response when creating a new API key.
 * The `key` field contains the full key - SAVE IT, it won't be shown again!
 */
export interface CreateAPIKeyResponse {
  id: number;
  key: string;           // Full API key - only shown once!
  key_prefix: string;
  description: string;
  created_at: string;
}

// ============================================================================
// Import/Export Types
// ============================================================================

/**
 * Result of importing bookmarks.
 */
export interface ImportResult {
  imported: number;   // Successfully imported
  skipped: number;    // Skipped (e.g., duplicates)
  errors?: string[];  // Any errors that occurred
}

/**
 * Pinboard bookmark format - used for import/export compatibility.
 * Pinboard is a popular bookmarking service with a standard format.
 */
export interface PinboardBookmark {
  href: string;        // URL
  description: string; // Title
  extended: string;    // Notes/description
  tags: string;        // Space-separated tags
  time: string;        // ISO 8601 timestamp
  shared: string;      // "yes" or "no"
  toread: string;      // "yes" or "no"
}

// ============================================================================
// Admin Types (for admin dashboard)
// ============================================================================

/**
 * Extended user info shown in admin user list.
 */
export interface AdminUser {
  id: number;
  email: string;
  name: string;
  system_role: 'admin' | 'user';
  created_at: string;
  link_count: number;   // Total links created by this user
  group_count: number;  // Groups this user belongs to
}

/**
 * System-wide statistics for the admin dashboard.
 */
export interface AdminStats {
  total_users: number;
  total_links: number;
  total_groups: number;
  total_tags: number;
  total_clicks: number;
  public_links: number;
  private_links: number;
  unread_links: number;
  admin_users: number;
  active_api_keys: number;
}

// ============================================================================
// OIDC (Single Sign-On) Types
// ============================================================================

/**
 * Public info about an OIDC provider (shown on login page).
 */
export interface OIDCProvider {
  id: number;
  name: string;     // Display name (e.g., "Google", "Okta")
  slug: string;     // URL-safe identifier
  enabled: boolean;
}

/**
 * Full OIDC provider config (admin only).
 */
export interface OIDCProviderAdmin {
  id: number;
  name: string;
  slug: string;
  issuer: string;         // OIDC issuer URL
  client_id: string;
  scopes: string;         // Space-separated OAuth scopes
  enabled: boolean;
  auto_provision: boolean; // Auto-create users on first login?
  created_at: string;
}

// ============================================================================
// SCIM Types (for identity provider integration)
// ============================================================================

/**
 * SCIM tokens allow identity providers to sync users/groups.
 */
export interface SCIMToken {
  id: number;
  token_prefix: string;
  description: string;
  last_used_at: string | null;
  created_at: string;
}

/**
 * Response when creating a SCIM token.
 */
export interface CreateSCIMTokenResponse {
  id: number;
  token: string;  // Full token - only shown once!
  token_prefix: string;
  description: string;
  created_at: string;
}
