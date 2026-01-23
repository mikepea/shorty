/**
 * AuthContext.tsx - Authentication State Management
 *
 * This file implements the "Context" pattern for sharing authentication state
 * across the entire application. Without Context, you'd have to pass user info
 * through every component as props - that's called "prop drilling" and it's tedious.
 *
 * Key concepts:
 * - Context: A way to share values between components without passing props
 * - Provider: A component that "provides" the context value to its children
 * - Custom Hook: A function (useAuth) that makes using the context easier
 * - useState: React hook for storing data that can change over time
 * - useEffect: React hook for running code when the component mounts or updates
 *
 * How it works:
 * 1. AuthProvider wraps the app and holds the authentication state
 * 2. Any component can call useAuth() to get the current user and auth functions
 * 3. When auth state changes, all components using useAuth() automatically re-render
 */

import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import type { User } from '../api/types';
import { auth } from '../api/client';

/**
 * TypeScript interface defining what's available in the auth context.
 * This helps catch errors at compile time instead of runtime.
 */
interface AuthContextType {
  user: User | null;           // Current logged-in user, or null if not logged in
  token: string | null;        // JWT token for API authentication
  isLoading: boolean;          // True while checking if user is logged in on app start
  login: (email: string, password: string) => Promise<void>;    // Function to log in
  register: (email: string, password: string, name: string) => Promise<void>;  // Function to register
  logout: () => Promise<void>; // Function to log out
  setToken: (token: string) => void;  // Function to set token (used by SSO)
}

/**
 * Create the context with undefined as the default value.
 * The actual value is provided by AuthProvider below.
 */
const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * AuthProvider - The component that provides auth state to the app
 *
 * This component:
 * 1. Stores the current user and token in state
 * 2. Checks localStorage on startup to restore the session
 * 3. Provides login/register/logout functions
 * 4. Wraps children with the context provider
 *
 * @param children - The child components that will have access to auth context
 */
export function AuthProvider({ children }: { children: ReactNode }) {
  // useState creates a piece of state that persists across re-renders.
  // The [value, setValue] pattern is called "array destructuring".
  // When you call setValue(newValue), React re-renders the component.
  const [user, setUser] = useState<User | null>(null);

  // This useState has an "initializer function" - the arrow function runs
  // only once when the component first mounts, to get the initial value.
  // This is more efficient than reading localStorage on every render.
  const [token, setTokenState] = useState<string | null>(() =>
    localStorage.getItem('token')
  );

  // Track whether we're still checking if the user is logged in
  const [isLoading, setIsLoading] = useState(true);

  /**
   * useEffect runs code after the component renders.
   *
   * The empty array [] at the end means "run once when component mounts".
   * But we have [token] which means "run when token changes".
   *
   * This effect:
   * 1. If there's a token, try to fetch the current user from the API
   * 2. If that fails (token expired/invalid), clear the token
   * 3. Set isLoading to false when done
   */
  useEffect(() => {
    if (token) {
      // Call the API to get current user info
      auth.me()
        .then(setUser)  // If successful, set the user
        .catch(() => {
          // If failed (token invalid), clear everything
          localStorage.removeItem('token');
          setTokenState(null);
        })
        .finally(() => setIsLoading(false));  // Always stop loading when done
    } else {
      // No token means no need to check, just stop loading
      setIsLoading(false);
    }
  }, [token]);  // Re-run this effect whenever token changes

  /**
   * Log in with email and password.
   *
   * "async" functions can use "await" to pause until a Promise resolves.
   * This makes asynchronous code look like synchronous code.
   */
  const login = async (email: string, password: string) => {
    const response = await auth.login(email, password);

    // Store token in localStorage so it survives page refreshes
    localStorage.setItem('token', response.token);

    // Update state to trigger re-renders in components using useAuth()
    setTokenState(response.token);
    setUser(response.user);
  };

  /**
   * Register a new account.
   * Similar to login - stores the token and user after successful registration.
   */
  const register = async (email: string, password: string, name: string) => {
    const response = await auth.register(email, password, name);
    localStorage.setItem('token', response.token);
    setTokenState(response.token);
    setUser(response.user);
  };

  /**
   * Log out the current user.
   *
   * The try/finally ensures we clear local state even if the API call fails.
   * For example, if the server is down, we still want to log the user out locally.
   */
  const logout = async () => {
    try {
      await auth.logout();  // Tell the server to invalidate the session
    } finally {
      // Always clear local state, even if the API call failed
      localStorage.removeItem('token');
      setTokenState(null);
      setUser(null);
    }
  };

  /**
   * Set a token directly (used by SSO callback).
   * This is separate from login because SSO provides the token from the OAuth flow.
   */
  const handleSetToken = (newToken: string) => {
    localStorage.setItem('token', newToken);
    setTokenState(newToken);
    // Note: user will be fetched automatically by the useEffect above
    // because we're changing the token state
  };

  // Render the Provider component that makes the context available.
  // All children (and their children, etc.) can now call useAuth().
  return (
    <AuthContext.Provider value={{
      user,
      token,
      isLoading,
      login,
      register,
      logout,
      setToken: handleSetToken
    }}>
      {children}
    </AuthContext.Provider>
  );
}

/**
 * useAuth - Custom hook for accessing authentication context
 *
 * Custom hooks are a pattern for reusing stateful logic. By convention,
 * they start with "use". This hook:
 * 1. Gets the context value
 * 2. Throws an error if used outside of AuthProvider (helps catch bugs)
 * 3. Returns the context value for the component to use
 *
 * Usage in any component:
 *   const { user, login, logout } = useAuth();
 */
export function useAuth() {
  const context = useContext(AuthContext);

  // This error helps catch a common mistake: using useAuth() in a component
  // that isn't wrapped by AuthProvider. Without this check, context would
  // be undefined and you'd get confusing errors later.
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  return context;
}
