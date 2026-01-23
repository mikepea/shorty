/**
 * Login.tsx - Login Page Component
 *
 * This page handles user authentication with email/password and SSO options.
 * It demonstrates several common React patterns:
 *
 * Key concepts:
 * - Form handling with controlled inputs
 * - useState for managing form state and UI state
 * - useEffect for loading data when component mounts
 * - Event handlers (onSubmit, onChange, onClick)
 * - Conditional rendering for loading states and errors
 * - Async operations with try/catch/finally
 */

import { useState, useEffect, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { oidcProviders } from '../api/client';
import type { OIDCProvider } from '../api/types';

export default function Login() {
  // ============================================================================
  // State Management
  // ============================================================================

  // Form field state - these are "controlled inputs" where React manages the value
  // useState returns [currentValue, setterFunction]
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  // UI state for feedback
  const [error, setError] = useState('');       // Error message to display
  const [isLoading, setIsLoading] = useState(false);  // Disable form while submitting

  // SSO providers loaded from the API
  const [providers, setProviders] = useState<OIDCProvider[]>([]);

  // Get the login function from our auth context
  const { login } = useAuth();

  // useNavigate returns a function to programmatically change the URL
  const navigate = useNavigate();

  // ============================================================================
  // Effects - Side effects that run when the component mounts
  // ============================================================================

  /**
   * Load available SSO providers when the component mounts.
   *
   * The empty dependency array [] means this runs once when the component
   * first renders, similar to componentDidMount in class components.
   */
  useEffect(() => {
    oidcProviders.list()
      .then(setProviders)  // If successful, store the providers
      .catch(() => {});     // Ignore errors - SSO is optional
  }, []);

  // ============================================================================
  // Event Handlers
  // ============================================================================

  /**
   * Handle form submission for email/password login.
   *
   * @param e - The form event. We call preventDefault() to stop the browser
   *            from doing a full page reload (default form behavior).
   */
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();  // Prevent default form submission (page reload)
    setError('');        // Clear any previous errors
    setIsLoading(true);  // Show loading state

    try {
      // Call the login function from AuthContext
      await login(email, password);
      // If successful, navigate to the dashboard
      navigate('/');
    } catch (err) {
      // If login fails, show the error message
      // The instanceof check safely extracts the message from Error objects
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      // finally runs whether the try succeeded or failed
      setIsLoading(false);
    }
  };

  /**
   * Handle SSO login - redirects to the identity provider.
   *
   * Unlike email/password login, SSO redirects the user to an external
   * authentication page, then back to our /sso/callback route.
   */
  const handleSSOLogin = async (provider: OIDCProvider) => {
    setError('');
    setIsLoading(true);

    try {
      // Build the callback URL where the SSO provider will redirect back to
      const returnUrl = window.location.origin + '/sso/callback';

      // Get the authorization URL from our backend
      const { auth_url } = await oidcProviders.getAuthURL(provider.slug, returnUrl);

      // Redirect to the SSO provider (leaves our app temporarily)
      window.location.href = auth_url;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to initiate SSO');
      setIsLoading(false);  // Only reset if we're not redirecting
    }
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <div className="auth-container">
      <h1>Login to Shorty</h1>

      {/* Conditional rendering: only show error div if there's an error message */}
      {error && <div className="error">{error}</div>}

      {/*
        SSO section - only rendered if there are providers available.
        This pattern (condition && JSX) is called "short-circuit evaluation".
        If providers.length is 0 (falsy), the JSX is never evaluated.
      */}
      {providers.length > 0 && (
        <div className="sso-section">
          <p className="sso-label">Sign in with</p>
          <div className="sso-buttons">
            {/*
              Map over providers to create a button for each.
              The "key" prop helps React track which items changed.
              Always use a unique identifier (like id) for keys, not array index.
            */}
            {providers.map((provider) => (
              <button
                key={provider.id}
                type="button"
                onClick={() => handleSSOLogin(provider)}
                disabled={isLoading}
                className="sso-button"
              >
                {provider.name}
              </button>
            ))}
          </div>
          <div className="divider">
            <span>or</span>
          </div>
        </div>
      )}

      {/*
        Login form with controlled inputs.
        "Controlled" means React state is the source of truth for the input value.
        The value prop sets what's displayed, onChange updates the state.
      */}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}                           // Controlled: displays state value
            onChange={(e) => setEmail(e.target.value)}  // Updates state on change
            required                                // HTML5 validation
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        {/* Button shows different text based on loading state */}
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Logging in...' : 'Login'}
        </button>
      </form>

      {/* Link component for client-side navigation (no page reload) */}
      <p>
        Don't have an account? <Link to="/register">Register</Link>
      </p>
    </div>
  );
}
