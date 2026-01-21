import { describe, it, expect, vi, beforeEach } from 'vitest';
import { auth, links, groups, apiKeys, APIError } from './client';

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(localStorage.getItem).mockReturnValue(null);
  });

  describe('request helper', () => {
    it('includes auth token when available', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('test-token');
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, email: 'test@example.com' }),
      } as Response);

      await auth.me();

      expect(fetch).toHaveBeenCalledWith(
        '/api/auth/me',
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token',
          }),
        })
      );
    });

    it('does not include auth header when no token', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: 'abc', user: {} }),
      } as Response);

      await auth.login('test@example.com', 'password');

      const fetchCall = vi.mocked(fetch).mock.calls[0];
      const headers = fetchCall[1]?.headers as Record<string, string>;
      expect(headers.Authorization).toBeUndefined();
    });

    it('throws APIError on non-ok response', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({ error: 'Invalid credentials' }),
      } as Response);

      await expect(auth.login('test@example.com', 'wrong')).rejects.toThrow(APIError);

      // Reset mock for second assertion
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({ error: 'Invalid credentials' }),
      } as Response);

      await expect(auth.login('test@example.com', 'wrong')).rejects.toMatchObject({
        status: 401,
        message: 'Invalid credentials',
      });
    });

    it('handles json parse errors gracefully', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.reject(new Error('Invalid JSON')),
      } as Response);

      await expect(auth.login('test@example.com', 'password')).rejects.toThrow(APIError);
    });
  });

  describe('auth', () => {
    it('register sends correct request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: 'abc', user: { id: 1 } }),
      } as Response);

      await auth.register('test@example.com', 'password123', 'Test User');

      expect(fetch).toHaveBeenCalledWith(
        '/api/auth/register',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            email: 'test@example.com',
            password: 'password123',
            name: 'Test User',
          }),
        })
      );
    });

    it('login sends correct request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ token: 'abc', user: { id: 1 } }),
      } as Response);

      await auth.login('test@example.com', 'password123');

      expect(fetch).toHaveBeenCalledWith(
        '/api/auth/login',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            email: 'test@example.com',
            password: 'password123',
          }),
        })
      );
    });

    it('changePassword sends correct request', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('test-token');
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ message: 'Password changed successfully' }),
      } as Response);

      await auth.changePassword('oldpass', 'newpass123');

      expect(fetch).toHaveBeenCalledWith(
        '/api/auth/password',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({
            current_password: 'oldpass',
            new_password: 'newpass123',
          }),
        })
      );
    });

    it('me fetches current user', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('test-token');
      const mockUser = { id: 1, email: 'test@example.com', name: 'Test', system_role: 'user' };
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockUser),
      } as Response);

      const result = await auth.me();

      expect(result).toEqual(mockUser);
      expect(fetch).toHaveBeenCalledWith('/api/auth/me', expect.any(Object));
    });
  });

  describe('links', () => {
    beforeEach(() => {
      vi.mocked(localStorage.getItem).mockReturnValue('test-token');
    });

    it('search without params', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve([]),
      } as Response);

      await links.search();

      expect(fetch).toHaveBeenCalledWith('/api/links', expect.any(Object));
    });

    it('search with query param', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve([]),
      } as Response);

      await links.search({ q: 'test' });

      expect(fetch).toHaveBeenCalledWith('/api/links?q=test', expect.any(Object));
    });

    it('search with multiple params', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve([]),
      } as Response);

      await links.search({ q: 'test', tag: 'javascript', unread: true });

      expect(fetch).toHaveBeenCalledWith(
        '/api/links?q=test&tag=javascript&unread=true',
        expect.any(Object)
      );
    });

    it('get fetches link by slug', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, slug: 'my-link' }),
      } as Response);

      await links.get('my-link');

      expect(fetch).toHaveBeenCalledWith('/api/links/my-link', expect.any(Object));
    });

    it('create sends POST to group links endpoint', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, slug: 'new-link' }),
      } as Response);

      await links.create(5, { url: 'https://example.com', title: 'Example' });

      expect(fetch).toHaveBeenCalledWith(
        '/api/groups/5/links',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ url: 'https://example.com', title: 'Example' }),
        })
      );
    });

    it('update sends PUT request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, slug: 'my-link' }),
      } as Response);

      await links.update('my-link', { title: 'Updated Title' });

      expect(fetch).toHaveBeenCalledWith(
        '/api/links/my-link',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ title: 'Updated Title' }),
        })
      );
    });

    it('delete sends DELETE request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ message: 'Deleted' }),
      } as Response);

      await links.delete('my-link');

      expect(fetch).toHaveBeenCalledWith(
        '/api/links/my-link',
        expect.objectContaining({ method: 'DELETE' })
      );
    });
  });

  describe('groups', () => {
    beforeEach(() => {
      vi.mocked(localStorage.getItem).mockReturnValue('test-token');
    });

    it('list fetches all groups', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve([{ id: 1, name: 'Group 1' }]),
      } as Response);

      const result = await groups.list();

      expect(fetch).toHaveBeenCalledWith('/api/groups', expect.any(Object));
      expect(result).toEqual([{ id: 1, name: 'Group 1' }]);
    });

    it('create sends POST request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, name: 'New Group' }),
      } as Response);

      await groups.create('New Group');

      expect(fetch).toHaveBeenCalledWith(
        '/api/groups',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ name: 'New Group' }),
        })
      );
    });

    it('addMember sends correct request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, user_id: 2, group_id: 1, role: 'member' }),
      } as Response);

      await groups.addMember(1, 'user@example.com', 'member');

      expect(fetch).toHaveBeenCalledWith(
        '/api/groups/1/members',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ email: 'user@example.com', role: 'member' }),
        })
      );
    });
  });

  describe('apiKeys', () => {
    beforeEach(() => {
      vi.mocked(localStorage.getItem).mockReturnValue('test-token');
    });

    it('list fetches all API keys', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve([]),
      } as Response);

      await apiKeys.list();

      expect(fetch).toHaveBeenCalledWith('/api/api-keys', expect.any(Object));
    });

    it('create sends POST request with description', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ id: 1, key: 'sk_abc123' }),
      } as Response);

      await apiKeys.create('My API Key');

      expect(fetch).toHaveBeenCalledWith(
        '/api/api-keys',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ description: 'My API Key' }),
        })
      );
    });

    it('delete sends DELETE request', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ message: 'Deleted' }),
      } as Response);

      await apiKeys.delete(123);

      expect(fetch).toHaveBeenCalledWith(
        '/api/api-keys/123',
        expect.objectContaining({ method: 'DELETE' })
      );
    });
  });
});
