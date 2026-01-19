import { useState, useEffect, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { oidcProviders } from '../api/client';
import type { OIDCProvider } from '../api/types';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [providers, setProviders] = useState<OIDCProvider[]>([]);
  const { login } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    oidcProviders.list()
      .then(setProviders)
      .catch(() => {}); // Ignore errors - SSO is optional
  }, []);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(email, password);
      navigate('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSSOLogin = async (provider: OIDCProvider) => {
    setError('');
    setIsLoading(true);

    try {
      const returnUrl = window.location.origin + '/sso/callback';
      const { auth_url } = await oidcProviders.getAuthURL(provider.slug, returnUrl);
      window.location.href = auth_url;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to initiate SSO');
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <h1>Login to Shorty</h1>
      {error && <div className="error">{error}</div>}

      {providers.length > 0 && (
        <div className="sso-section">
          <p className="sso-label">Sign in with</p>
          <div className="sso-buttons">
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

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
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
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Logging in...' : 'Login'}
        </button>
      </form>
      <p>
        Don't have an account? <Link to="/register">Register</Link>
      </p>
    </div>
  );
}
