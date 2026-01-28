/**
 * client.ts - API Client
 *
 * This file provides a clean interface for calling the backend API.
 * Instead of using fetch() directly everywhere, components can call
 * functions like `auth.login(email, password)` or `links.create(...)`.
 *
 * Key concepts:
 * - fetch(): Browser API for making HTTP requests
 * - async/await: Modern way to handle asynchronous operations
 * - Generic types (<T>): Allow functions to work with different types
 * - Object shorthand: { auth, links } is shorthand for { auth: auth, links: links }
 *
 * How it works:
 * 1. The `request` helper handles common logic (auth headers, JSON parsing, errors)
 * 2. Exported objects (auth, links, etc.) group related API calls
 * 3. Each method returns a Promise that resolves to the typed response
 */

import type {
  User,
  Organization,
  OrganizationMember,
  Group,
  GroupMembership,
  Link,
  Tag,
  APIKey,
  CreateAPIKeyResponse,
  AuthResponse,
  ImportResult,
  PinboardBookmark,
  AdminUser,
  AdminStats,
  OIDCProvider,
  OIDCProviderAdmin,
  SCIMToken,
  CreateSCIMTokenResponse,
} from './types';

// Base URL for all API calls. In development, Vite proxies /api to the Go server.
const API_BASE = '/api';

// Key used to store the current organization ID in localStorage
const ORG_STORAGE_KEY = 'current_org_id';

/**
 * Get the currently selected organization ID from localStorage.
 * Returns null if no organization is selected or if the stored value is invalid.
 */
export function getCurrentOrgId(): number | null {
  const stored = localStorage.getItem(ORG_STORAGE_KEY);
  if (!stored) return null;
  const parsed = parseInt(stored, 10);
  // Handle NaN from corrupted/tampered localStorage
  return Number.isNaN(parsed) ? null : parsed;
}

/**
 * Set the current organization ID in localStorage.
 * This will be sent with all API requests via the X-Organization-ID header.
 */
export function setCurrentOrgId(orgId: number | null): void {
  if (orgId === null) {
    localStorage.removeItem(ORG_STORAGE_KEY);
  } else {
    localStorage.setItem(ORG_STORAGE_KEY, orgId.toString());
  }
}

/**
 * Custom error class for API errors.
 *
 * This extends the built-in Error class to include the HTTP status code.
 * Components can check error.status to handle different errors differently
 * (e.g., 401 = not authenticated, 404 = not found, 403 = forbidden).
 */
class APIError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);     // Call parent constructor with the message
    this.status = status;
    this.name = 'APIError';  // Helps identify this error type in catch blocks
  }
}

/**
 * Generic request helper that handles common API call logic.
 *
 * The <T> is a "generic type parameter" - it lets this function work with
 * any response type. When you call request<User>(...), T becomes User,
 * so the function returns Promise<User>.
 *
 * @param endpoint - The API path (e.g., '/auth/login')
 * @param options - fetch options (method, body, headers, etc.)
 * @returns Promise resolving to the parsed JSON response
 * @throws APIError if the response is not ok (status >= 400)
 */
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  // Get the JWT token from localStorage (set during login)
  const token = localStorage.getItem('token');

  // Get the current organization ID from localStorage
  const orgId = getCurrentOrgId();

  // Set up headers - always send JSON
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,  // Spread operator: merge in any custom headers
  };

  // If we have a token, add the Authorization header
  // The "Bearer" prefix is a standard convention for JWT tokens
  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
  }

  // If we have an organization selected, include it in the request
  // This tells the backend which organization context to use
  if (orgId) {
    (headers as Record<string, string>)['X-Organization-ID'] = orgId.toString();
  }

  // Make the actual HTTP request
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  // Check if the request was successful (status 200-299)
  if (!response.ok) {
    // Try to parse error message from response body
    const error = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new APIError(response.status, error.error || 'Request failed');
  }

  // Parse and return the JSON response
  return response.json();
}

// ============================================================================
// Authentication API
// ============================================================================

/**
 * Auth-related API calls.
 *
 * Usage:
 *   import { auth } from './api/client';
 *   const { token, user } = await auth.login('email@example.com', 'password');
 */
export const auth = {
  /**
   * Register a new user account.
   */
  register: (email: string, password: string, name: string) =>
    request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
    }),

  /**
   * Log in with email and password.
   */
  login: (email: string, password: string) =>
    request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  /**
   * Log out (invalidates the session on the server).
   */
  logout: () =>
    request<{ message: string }>('/auth/logout', { method: 'POST' }),

  /**
   * Get the current logged-in user's info.
   * Used to verify the token is still valid on app startup.
   */
  me: () => request<User>('/auth/me'),

  /**
   * Change the current user's password.
   * Requires the current password for verification.
   */
  changePassword: (currentPassword: string, newPassword: string) =>
    request<{ message: string }>('/auth/password', {
      method: 'PUT',
      body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }),
    }),
};

// ============================================================================
// Organizations API
// ============================================================================

/**
 * Organization management API calls.
 *
 * Organizations are the top-level container for multi-tenancy. They scope
 * SSO settings, SCIM provisioning, groups, and link slugs. Users can belong
 * to multiple organizations and switch between them.
 */
export const organizations = {
  /** Get all organizations the current user is a member of. */
  list: () => request<Organization[]>('/organizations'),

  /** Get a specific organization by ID. */
  get: (id: number) => request<Organization>(`/organizations/${id}`),

  /** Create a new organization. The creator becomes an admin. */
  create: (name: string, slug: string) =>
    request<Organization>('/organizations', {
      method: 'POST',
      body: JSON.stringify({ name, slug }),
    }),

  /** Update an organization's name (admin only). */
  update: (id: number, name: string) =>
    request<Organization>(`/organizations/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    }),

  /** Delete an organization (admin only). Cannot delete the global org. */
  delete: (id: number) =>
    request<{ message: string }>(`/organizations/${id}`, { method: 'DELETE' }),

  /** Get all members of an organization. */
  listMembers: (id: number) =>
    request<OrganizationMember[]>(`/organizations/${id}/members`),

  /** Add a user to an organization by their email (admin only). */
  addMember: (orgId: number, email: string, role: 'admin' | 'member') =>
    request<OrganizationMember>(`/organizations/${orgId}/members`, {
      method: 'POST',
      body: JSON.stringify({ email, role }),
    }),

  /** Change a member's role in an organization (admin only). */
  updateMember: (orgId: number, userId: number, role: 'admin' | 'member') =>
    request<OrganizationMember>(`/organizations/${orgId}/members/${userId}`, {
      method: 'PUT',
      body: JSON.stringify({ role }),
    }),

  /** Remove a member from an organization (admin only). */
  removeMember: (orgId: number, userId: number) =>
    request<{ message: string }>(`/organizations/${orgId}/members/${userId}`, {
      method: 'DELETE',
    }),
};

// ============================================================================
// Groups API
// ============================================================================

/**
 * Group management API calls.
 *
 * Groups are collections of links that can be shared between users.
 */
export const groups = {
  /** Get all groups the current user is a member of. */
  list: () => request<Group[]>('/groups'),

  /** Get a specific group by ID. */
  get: (id: number) => request<Group>(`/groups/${id}`),

  /** Create a new group. */
  create: (name: string) =>
    request<Group>('/groups', {
      method: 'POST',
      body: JSON.stringify({ name }),
    }),

  /** Update a group's name. */
  update: (id: number, name: string) =>
    request<Group>(`/groups/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    }),

  /** Delete a group. */
  delete: (id: number) =>
    request<{ message: string }>(`/groups/${id}`, { method: 'DELETE' }),

  /** Get all members of a group. */
  listMembers: (id: number) =>
    request<GroupMembership[]>(`/groups/${id}/members`),

  /** Add a user to a group by their email. */
  addMember: (groupId: number, email: string, role: 'admin' | 'member') =>
    request<GroupMembership>(`/groups/${groupId}/members`, {
      method: 'POST',
      body: JSON.stringify({ email, role }),
    }),

  /** Change a member's role in a group. */
  updateMember: (groupId: number, userId: number, role: 'admin' | 'member') =>
    request<GroupMembership>(`/groups/${groupId}/members/${userId}`, {
      method: 'PUT',
      body: JSON.stringify({ role }),
    }),

  /** Remove a member from a group. */
  removeMember: (groupId: number, userId: number) =>
    request<{ message: string }>(`/groups/${groupId}/members/${userId}`, {
      method: 'DELETE',
    }),
};

// ============================================================================
// Links API
// ============================================================================

/**
 * Link management API calls.
 *
 * Links are the core entity - shortened URLs with metadata.
 */
export const links = {
  /**
   * Search links with optional filters.
   *
   * URLSearchParams is a browser API that helps build query strings.
   * For example: { q: 'hello', tag: 'work' } becomes '?q=hello&tag=work'
   */
  search: (params?: { q?: string; tag?: string; unread?: boolean }) => {
    const searchParams = new URLSearchParams();
    if (params?.q) searchParams.set('q', params.q);
    if (params?.tag) searchParams.set('tag', params.tag);
    if (params?.unread) searchParams.set('unread', 'true');
    const query = searchParams.toString();
    return request<Link[]>(`/links${query ? `?${query}` : ''}`);
  },

  /** Get all links in a specific group. */
  listByGroup: (groupId: number) =>
    request<Link[]>(`/groups/${groupId}/links`),

  /** Get a single link by its slug. */
  get: (slug: string) => request<Link>(`/links/${slug}`),

  /**
   * Create a new link in a group.
   *
   * The data parameter uses TypeScript inline type definition.
   * Optional properties (title?, description?) may be omitted.
   */
  create: (groupId: number, data: {
    url: string;
    title?: string;
    description?: string;
    slug?: string;       // Custom slug (auto-generated if not provided)
    is_public?: boolean;
  }) =>
    request<Link>(`/groups/${groupId}/links`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  /** Update an existing link. */
  update: (slug: string, data: {
    url?: string;
    title?: string;
    description?: string;
    is_public?: boolean;
    is_unread?: boolean;
  }) =>
    request<Link>(`/links/${slug}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  /** Delete a link. */
  delete: (slug: string) =>
    request<{ message: string }>(`/links/${slug}`, { method: 'DELETE' }),
};

// ============================================================================
// Tags API
// ============================================================================

/**
 * Tag management API calls.
 *
 * Tags are labels for organizing links. They're created automatically
 * when you add them to a link.
 */
export const tags = {
  /** Get all tags across all groups the user can access. */
  list: () => request<Tag[]>('/tags'),

  /** Get tags used in a specific group. */
  listByGroup: (groupId: number) =>
    request<Tag[]>(`/groups/${groupId}/tags`),

  /** Get tags on a specific link. */
  getLinkTags: (slug: string) =>
    request<Tag[]>(`/links/${slug}/tags`),

  /** Replace all tags on a link. */
  setLinkTags: (slug: string, tagNames: string[]) =>
    request<Tag[]>(`/links/${slug}/tags`, {
      method: 'PUT',
      body: JSON.stringify({ tags: tagNames }),
    }),

  /** Add a single tag to a link. */
  addLinkTag: (slug: string, tagName: string) =>
    request<Tag>(`/links/${slug}/tags/${encodeURIComponent(tagName)}`, {
      method: 'POST',
    }),

  /** Remove a single tag from a link. */
  removeLinkTag: (slug: string, tagName: string) =>
    request<{ message: string }>(`/links/${slug}/tags/${encodeURIComponent(tagName)}`, {
      method: 'DELETE',
    }),
};

// ============================================================================
// API Keys
// ============================================================================

/**
 * API key management.
 *
 * API keys allow programmatic access without username/password.
 * Useful for scripts, integrations, and automation.
 */
export const apiKeys = {
  /** List all API keys for the current user. */
  list: () => request<APIKey[]>('/api-keys'),

  /** Create a new API key. Returns the full key (only shown once!). */
  create: (description?: string) =>
    request<CreateAPIKeyResponse>('/api-keys', {
      method: 'POST',
      body: JSON.stringify({ description }),
    }),

  /** Delete (revoke) an API key. */
  delete: (id: number) =>
    request<{ message: string }>(`/api-keys/${id}`, { method: 'DELETE' }),
};

// ============================================================================
// Import/Export
// ============================================================================

/**
 * Bulk import and export of bookmarks.
 *
 * Uses Pinboard format for compatibility with other bookmark managers.
 */
export const importExport = {
  /** Import bookmarks into a group. */
  import: (groupId: number, bookmarks: PinboardBookmark[]) =>
    request<ImportResult>('/import', {
      method: 'POST',
      body: JSON.stringify({ group_id: groupId, bookmarks }),
    }),

  /** Export bookmarks (optionally from a specific group). */
  export: (groupId?: number) => {
    const query = groupId ? `?group_id=${groupId}` : '';
    return request<PinboardBookmark[]>(`/export${query}`);
  },

  /** Export a single link. */
  exportSingle: (slug: string) =>
    request<PinboardBookmark>(`/export/${slug}`),
};

// ============================================================================
// OIDC (Single Sign-On)
// ============================================================================

/**
 * OIDC provider endpoints for SSO login.
 */
export const oidcProviders = {
  /** Get available SSO providers (shown on login page). */
  list: () => request<OIDCProvider[]>('/oidc/providers'),

  /** Get the OAuth authorization URL to redirect the user to. */
  getAuthURL: (slug: string, returnUrl?: string) =>
    request<{ auth_url: string }>(`/oidc/providers/${slug}/auth`, {
      method: 'POST',
      body: JSON.stringify({ return_url: returnUrl }),
    }),
};

// ============================================================================
// Admin API
// ============================================================================

/**
 * Admin-only endpoints for system management.
 *
 * These require the user to have system_role: 'admin'.
 */
export const admin = {
  /** Get system-wide statistics. */
  getStats: () => request<AdminStats>('/admin/stats'),

  /** List all users with optional search/filter. */
  listUsers: (params?: { q?: string; role?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.q) searchParams.set('q', params.q);
    if (params?.role) searchParams.set('role', params.role);
    const query = searchParams.toString();
    return request<AdminUser[]>(`/admin/users${query ? `?${query}` : ''}`);
  },

  /** Get a specific user by ID. */
  getUser: (id: number) => request<AdminUser>(`/admin/users/${id}`),

  /** Update a user (name, role). */
  updateUser: (id: number, data: { name?: string; system_role?: string }) =>
    request<AdminUser>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  /** Delete a user and all their data. */
  deleteUser: (id: number) =>
    request<{ message: string }>(`/admin/users/${id}`, { method: 'DELETE' }),

  // --- OIDC Provider Management ---

  /** List all OIDC providers (with full config). */
  listOIDCProviders: () => request<OIDCProviderAdmin[]>('/admin/oidc/providers'),

  /** Create a new OIDC provider. */
  createOIDCProvider: (data: {
    name: string;
    slug: string;
    issuer: string;
    client_id: string;
    client_secret: string;
    scopes?: string;
    enabled?: boolean;
    auto_provision?: boolean;
  }) =>
    request<OIDCProviderAdmin>('/admin/oidc/providers', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  /** Update an OIDC provider. */
  updateOIDCProvider: (id: number, data: {
    name?: string;
    issuer?: string;
    client_id?: string;
    client_secret?: string;
    scopes?: string;
    enabled?: boolean;
    auto_provision?: boolean;
  }) =>
    request<OIDCProviderAdmin>(`/admin/oidc/providers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  /** Delete an OIDC provider. */
  deleteOIDCProvider: (id: number) =>
    request<{ message: string }>(`/admin/oidc/providers/${id}`, { method: 'DELETE' }),

  // --- SCIM Token Management ---

  /** List all SCIM tokens. */
  listSCIMTokens: () => request<SCIMToken[]>('/admin/scim-tokens'),

  /** Create a new SCIM token. Returns the full token (only shown once!). */
  createSCIMToken: (description?: string) =>
    request<CreateSCIMTokenResponse>('/admin/scim-tokens', {
      method: 'POST',
      body: JSON.stringify({ description }),
    }),

  /** Delete (revoke) a SCIM token. */
  deleteSCIMToken: (id: number) =>
    request<{ message: string }>(`/admin/scim-tokens/${id}`, { method: 'DELETE' }),
};

// Export the APIError class so components can catch specific API errors
export { APIError };
