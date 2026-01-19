import type {
  User,
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
} from './types';

const API_BASE = '/api';

class APIError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = 'APIError';
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem('token');

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new APIError(response.status, error.error || 'Request failed');
  }

  return response.json();
}

// Auth
export const auth = {
  register: (email: string, password: string, name: string) =>
    request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
    }),

  login: (email: string, password: string) =>
    request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  logout: () =>
    request<{ message: string }>('/auth/logout', { method: 'POST' }),

  me: () => request<User>('/auth/me'),
};

// Groups
export const groups = {
  list: () => request<Group[]>('/groups'),

  get: (id: number) => request<Group>(`/groups/${id}`),

  create: (name: string) =>
    request<Group>('/groups', {
      method: 'POST',
      body: JSON.stringify({ name }),
    }),

  update: (id: number, name: string) =>
    request<Group>(`/groups/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    }),

  delete: (id: number) =>
    request<{ message: string }>(`/groups/${id}`, { method: 'DELETE' }),

  listMembers: (id: number) =>
    request<GroupMembership[]>(`/groups/${id}/members`),

  addMember: (groupId: number, email: string, role: 'admin' | 'member') =>
    request<GroupMembership>(`/groups/${groupId}/members`, {
      method: 'POST',
      body: JSON.stringify({ email, role }),
    }),

  updateMember: (groupId: number, userId: number, role: 'admin' | 'member') =>
    request<GroupMembership>(`/groups/${groupId}/members/${userId}`, {
      method: 'PUT',
      body: JSON.stringify({ role }),
    }),

  removeMember: (groupId: number, userId: number) =>
    request<{ message: string }>(`/groups/${groupId}/members/${userId}`, {
      method: 'DELETE',
    }),
};

// Links
export const links = {
  search: (params?: { q?: string; tag?: string; unread?: boolean }) => {
    const searchParams = new URLSearchParams();
    if (params?.q) searchParams.set('q', params.q);
    if (params?.tag) searchParams.set('tag', params.tag);
    if (params?.unread) searchParams.set('unread', 'true');
    const query = searchParams.toString();
    return request<Link[]>(`/links${query ? `?${query}` : ''}`);
  },

  listByGroup: (groupId: number) =>
    request<Link[]>(`/groups/${groupId}/links`),

  get: (slug: string) => request<Link>(`/links/${slug}`),

  create: (groupId: number, data: {
    url: string;
    title?: string;
    description?: string;
    slug?: string;
    is_public?: boolean;
  }) =>
    request<Link>(`/groups/${groupId}/links`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

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

  delete: (slug: string) =>
    request<{ message: string }>(`/links/${slug}`, { method: 'DELETE' }),
};

// Tags
export const tags = {
  list: () => request<Tag[]>('/tags'),

  listByGroup: (groupId: number) =>
    request<Tag[]>(`/groups/${groupId}/tags`),

  getLinkTags: (slug: string) =>
    request<Tag[]>(`/links/${slug}/tags`),

  setLinkTags: (slug: string, tagNames: string[]) =>
    request<Tag[]>(`/links/${slug}/tags`, {
      method: 'PUT',
      body: JSON.stringify({ tags: tagNames }),
    }),

  addLinkTag: (slug: string, tagName: string) =>
    request<Tag>(`/links/${slug}/tags/${encodeURIComponent(tagName)}`, {
      method: 'POST',
    }),

  removeLinkTag: (slug: string, tagName: string) =>
    request<{ message: string }>(`/links/${slug}/tags/${encodeURIComponent(tagName)}`, {
      method: 'DELETE',
    }),
};

// API Keys
export const apiKeys = {
  list: () => request<APIKey[]>('/api-keys'),

  create: (description?: string) =>
    request<CreateAPIKeyResponse>('/api-keys', {
      method: 'POST',
      body: JSON.stringify({ description }),
    }),

  delete: (id: number) =>
    request<{ message: string }>(`/api-keys/${id}`, { method: 'DELETE' }),
};

// Import/Export
export const importExport = {
  import: (groupId: number, bookmarks: PinboardBookmark[]) =>
    request<ImportResult>('/import', {
      method: 'POST',
      body: JSON.stringify({ group_id: groupId, bookmarks }),
    }),

  export: (groupId?: number) => {
    const query = groupId ? `?group_id=${groupId}` : '';
    return request<PinboardBookmark[]>(`/export${query}`);
  },

  exportSingle: (slug: string) =>
    request<PinboardBookmark>(`/export/${slug}`),
};

// Admin
export const admin = {
  getStats: () => request<AdminStats>('/admin/stats'),

  listUsers: (params?: { q?: string; role?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.q) searchParams.set('q', params.q);
    if (params?.role) searchParams.set('role', params.role);
    const query = searchParams.toString();
    return request<AdminUser[]>(`/admin/users${query ? `?${query}` : ''}`);
  },

  getUser: (id: number) => request<AdminUser>(`/admin/users/${id}`),

  updateUser: (id: number, data: { name?: string; system_role?: string }) =>
    request<AdminUser>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteUser: (id: number) =>
    request<{ message: string }>(`/admin/users/${id}`, { method: 'DELETE' }),
};

export { APIError };
