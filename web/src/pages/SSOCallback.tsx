import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function SSOCallback() {
  const [searchParams] = useSearchParams();
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { setToken } = useAuth();

  useEffect(() => {
    const token = searchParams.get('token');
    const errorParam = searchParams.get('error');

    if (errorParam) {
      setError(errorParam);
      return;
    }

    if (token) {
      // Store token and redirect to home
      setToken(token);
      navigate('/', { replace: true });
    } else {
      setError('No token received from SSO provider');
    }
  }, [searchParams, navigate, setToken]);

  if (error) {
    return (
      <div className="auth-container">
        <h1>SSO Login Failed</h1>
        <div className="error">{error}</div>
        <p>
          <a href="/login">Back to Login</a>
        </p>
      </div>
    );
  }

  return (
    <div className="auth-container">
      <h1>Completing SSO Login...</h1>
      <div className="loading">Please wait...</div>
    </div>
  );
}
