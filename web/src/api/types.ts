// API Types matching the Go backend

export interface User {
  id: number;
  email: string;
  name: string;
  system_role: 'admin' | 'user';
  created_at: string;
}

export interface Group {
  id: number;
  name: string;
  created_at: string;
}

export interface GroupMembership {
  id: number;
  user_id: number;
  group_id: number;
  role: 'admin' | 'member';
  user?: User;
  group?: Group;
}

export interface Tag {
  id: number;
  name: string;
  link_count?: number;
}

export interface Link {
  id: number;
  group_id: number;
  created_by_id: number;
  slug: string;
  url: string;
  title: string;
  description: string;
  is_public: boolean;
  is_unread: boolean;
  click_count: number;
  created_at: string;
  updated_at: string;
  tags?: Tag[];
  group?: Group;
}

export interface APIKey {
  id: number;
  key_prefix: string;
  description: string;
  last_used_at: string | null;
  created_at: string;
}

export interface CreateAPIKeyResponse {
  id: number;
  key: string;
  key_prefix: string;
  description: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ImportResult {
  imported: number;
  skipped: number;
  errors?: string[];
}

export interface PinboardBookmark {
  href: string;
  description: string;
  extended: string;
  tags: string;
  time: string;
  shared: string;
  toread: string;
}

// Admin types
export interface AdminUser {
  id: number;
  email: string;
  name: string;
  system_role: 'admin' | 'user';
  created_at: string;
  link_count: number;
  group_count: number;
}

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

// OIDC types
export interface OIDCProvider {
  id: number;
  name: string;
  slug: string;
  enabled: boolean;
}

export interface OIDCProviderAdmin {
  id: number;
  name: string;
  slug: string;
  issuer: string;
  client_id: string;
  scopes: string;
  enabled: boolean;
  auto_provision: boolean;
  created_at: string;
}

// SCIM types
export interface SCIMToken {
  id: number;
  token_prefix: string;
  description: string;
  last_used_at: string | null;
  created_at: string;
}

export interface CreateSCIMTokenResponse {
  id: number;
  token: string;
  token_prefix: string;
  description: string;
  created_at: string;
}
