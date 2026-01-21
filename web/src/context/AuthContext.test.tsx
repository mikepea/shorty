import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AuthProvider, useAuth } from './AuthContext';
import { auth } from '../api/client';

// Mock the API client
vi.mock('../api/client', () => ({
  auth: {
    me: vi.fn(),
    login: vi.fn(),
    register: vi.fn(),
    logout: vi.fn(),
  },
}));

// Test component that uses the auth context
function TestComponent() {
  const { user, token, isLoading, login, register, logout } = useAuth();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <div data-testid="user">{user ? user.email : 'No user'}</div>
      <div data-testid="token">{token || 'No token'}</div>
      <button onClick={() => login('test@example.com', 'password')}>Login</button>
      <button onClick={() => register('new@example.com', 'password', 'New User')}>Register</button>
      <button onClick={() => logout()}>Logout</button>
    </div>
  );
}

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(localStorage.getItem).mockReturnValue(null);
    vi.mocked(localStorage.setItem).mockImplementation(() => {});
    vi.mocked(localStorage.removeItem).mockImplementation(() => {});
  });

  describe('useAuth hook', () => {
    it('throws error when used outside AuthProvider', () => {
      // Suppress console.error for this test
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      expect(() => {
        render(<TestComponent />);
      }).toThrow('useAuth must be used within an AuthProvider');

      consoleSpy.mockRestore();
    });
  });

  describe('AuthProvider', () => {
    it('shows loading state initially when token exists', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('existing-token');
      vi.mocked(auth.me).mockImplementation(() => new Promise(() => {})); // Never resolves

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('loads user from token on mount', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('existing-token');
      vi.mocked(auth.me).mockResolvedValue({
        id: 1,
        email: 'existing@example.com',
        name: 'Existing User',
        system_role: 'user',
        has_password: true,
        created_at: '2024-01-01',
      });

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('existing@example.com');
      });
      expect(screen.getByTestId('token')).toHaveTextContent('existing-token');
    });

    it('clears token when auth.me fails', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue('invalid-token');
      vi.mocked(auth.me).mockRejectedValue(new Error('Unauthorized'));

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('No user');
      });
      expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    });

    it('shows no user when no token exists', async () => {
      vi.mocked(localStorage.getItem).mockReturnValue(null);

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('No user');
      });
      expect(screen.getByTestId('token')).toHaveTextContent('No token');
      expect(auth.me).not.toHaveBeenCalled();
    });
  });

  describe('login', () => {
    it('stores token and sets user on successful login', async () => {
      const user = userEvent.setup();
      vi.mocked(auth.login).mockResolvedValue({
        token: 'new-token',
        user: {
          id: 1,
          email: 'test@example.com',
          name: 'Test User',
          system_role: 'user',
          has_password: true,
          created_at: '2024-01-01',
        },
      });

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
      });

      await user.click(screen.getByText('Login'));

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('test@example.com');
      });
      expect(localStorage.setItem).toHaveBeenCalledWith('token', 'new-token');
      expect(auth.login).toHaveBeenCalledWith('test@example.com', 'password');
    });

    it('propagates error on failed login', async () => {
      const user = userEvent.setup();
      vi.mocked(auth.login).mockRejectedValue(new Error('Invalid credentials'));

      // Create a component that catches the error
      function TestWithErrorHandling() {
        const { login } = useAuth();
        const handleLogin = async () => {
          try {
            await login('test@example.com', 'password');
          } catch {
            // Error expected
          }
        };
        return <button onClick={handleLogin}>Login</button>;
      }

      render(
        <AuthProvider>
          <TestWithErrorHandling />
        </AuthProvider>
      );

      await user.click(screen.getByText('Login'));

      expect(localStorage.setItem).not.toHaveBeenCalled();
    });
  });

  describe('register', () => {
    it('stores token and sets user on successful registration', async () => {
      const user = userEvent.setup();
      vi.mocked(auth.register).mockResolvedValue({
        token: 'new-token',
        user: {
          id: 2,
          email: 'new@example.com',
          name: 'New User',
          system_role: 'user',
          has_password: true,
          created_at: '2024-01-01',
        },
      });

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
      });

      await user.click(screen.getByText('Register'));

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('new@example.com');
      });
      expect(localStorage.setItem).toHaveBeenCalledWith('token', 'new-token');
      expect(auth.register).toHaveBeenCalledWith('new@example.com', 'password', 'New User');
    });
  });

  describe('logout', () => {
    it('clears token and user on logout', async () => {
      const user = userEvent.setup();
      vi.mocked(localStorage.getItem).mockReturnValue('existing-token');
      vi.mocked(auth.me).mockResolvedValue({
        id: 1,
        email: 'test@example.com',
        name: 'Test',
        system_role: 'user',
        has_password: true,
        created_at: '2024-01-01',
      });
      vi.mocked(auth.logout).mockResolvedValue({ message: 'Logged out' });

      render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('test@example.com');
      });

      await user.click(screen.getByText('Logout'));

      await waitFor(() => {
        expect(screen.getByTestId('user')).toHaveTextContent('No user');
      });
      expect(screen.getByTestId('token')).toHaveTextContent('No token');
      expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    });

    it('clears local state even if API logout fails', async () => {
      const user = userEvent.setup();
      vi.mocked(localStorage.getItem).mockReturnValue('existing-token');
      vi.mocked(auth.me).mockResolvedValue({
        id: 1,
        email: 'test@example.com',
        name: 'Test',
        system_role: 'user',
        has_password: true,
        created_at: '2024-01-01',
      });
      vi.mocked(auth.logout).mockRejectedValue(new Error('Network error'));

      // Create a component that catches the logout error
      function TestWithErrorHandling() {
        const { user: authUser, logout } = useAuth();
        const handleLogout = async () => {
          try {
            await logout();
          } catch {
            // Error expected but state should still be cleared
          }
        };
        return (
          <div>
            <div data-testid="user-email">{authUser ? authUser.email : 'No user'}</div>
            <button onClick={handleLogout}>Logout</button>
          </div>
        );
      }

      render(
        <AuthProvider>
          <TestWithErrorHandling />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(screen.getByTestId('user-email')).toHaveTextContent('test@example.com');
      });

      await user.click(screen.getByText('Logout'));

      await waitFor(() => {
        expect(screen.getByTestId('user-email')).toHaveTextContent('No user');
      });
      expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    });
  });
});
