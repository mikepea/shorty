import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { apiKeys, auth } from '../api/client';
import type { APIKey, CreateAPIKeyResponse } from '../api/types';

export default function Settings() {
  const { user } = useAuth();
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [newKeyDescription, setNewKeyDescription] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [newKey, setNewKey] = useState<CreateAPIKeyResponse | null>(null);
  const [error, setError] = useState('');

  // Password change state
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [passwordSuccess, setPasswordSuccess] = useState('');
  const [isChangingPassword, setIsChangingPassword] = useState(false);

  useEffect(() => {
    apiKeys.list()
      .then(setKeys)
      .finally(() => setIsLoading(false));
  }, []);

  const handleCreateKey = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsCreating(true);
    setError('');
    setNewKey(null);

    try {
      const key = await apiKeys.create(newKeyDescription || undefined);
      setNewKey(key);
      setKeys([...keys, {
        id: key.id,
        key_prefix: key.key_prefix,
        description: key.description,
        last_used_at: null,
        created_at: key.created_at,
      }]);
      setNewKeyDescription('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create API key');
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteKey = async (key: APIKey) => {
    if (!confirm('Are you sure you want to delete this API key?')) return;

    try {
      await apiKeys.delete(key.id);
      setKeys(keys.filter((k) => k.id !== key.id));
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete API key');
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    alert('API key copied to clipboard!');
  };

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordError('');
    setPasswordSuccess('');

    // Client-side validation
    if (newPassword.length < 8) {
      setPasswordError('New password must be at least 8 characters');
      return;
    }

    if (newPassword !== confirmPassword) {
      setPasswordError('New passwords do not match');
      return;
    }

    setIsChangingPassword(true);
    try {
      await auth.changePassword(currentPassword, newPassword);
      setPasswordSuccess('Password changed successfully');
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err) {
      setPasswordError(err instanceof Error ? err.message : 'Failed to change password');
    } finally {
      setIsChangingPassword(false);
    }
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="settings-page">
      <header className="page-header">
        <h1>Settings</h1>
      </header>

      <div className="settings-content">
        <section className="profile-section">
          <h2>Profile</h2>
          <div className="profile-info">
            <div className="profile-field">
              <label>Name</label>
              <span>{user?.name}</span>
            </div>
            <div className="profile-field">
              <label>Email</label>
              <span>{user?.email}</span>
            </div>
            <div className="profile-field">
              <label>Role</label>
              <span>{user?.system_role}</span>
            </div>
          </div>
        </section>

        {user?.has_password && (
          <section className="password-section">
            <h2>Change Password</h2>
            <p className="section-description">
              Update your account password.
            </p>

            {passwordError && <div className="error">{passwordError}</div>}
            {passwordSuccess && <div className="success">{passwordSuccess}</div>}

            <form onSubmit={handleChangePassword} className="password-form">
              <div className="form-group">
                <label htmlFor="current-password">Current Password</label>
                <input
                  id="current-password"
                  type="password"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="new-password">New Password</label>
                <input
                  id="new-password"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  required
                  minLength={8}
                />
                <small>Minimum 8 characters</small>
              </div>
              <div className="form-group">
                <label htmlFor="confirm-password">Confirm New Password</label>
                <input
                  id="confirm-password"
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
              </div>
              <button type="submit" disabled={isChangingPassword}>
                {isChangingPassword ? 'Changing...' : 'Change Password'}
              </button>
            </form>
          </section>
        )}

        <section className="api-keys-section">
          <h2>API Keys</h2>
          <p className="section-description">
            API keys allow programmatic access to your account. Keep them secret!
          </p>

          {error && <div className="error">{error}</div>}

          {newKey && (
            <div className="new-key-alert">
              <strong>New API Key Created!</strong>
              <p>Copy this key now - you won't be able to see it again:</p>
              <code onClick={() => copyToClipboard(newKey.key)}>{newKey.key}</code>
              <button onClick={() => setNewKey(null)}>Dismiss</button>
            </div>
          )}

          <form onSubmit={handleCreateKey} className="create-key-form">
            <input
              type="text"
              value={newKeyDescription}
              onChange={(e) => setNewKeyDescription(e.target.value)}
              placeholder="Description (optional)"
            />
            <button type="submit" disabled={isCreating}>
              {isCreating ? 'Creating...' : 'Create API Key'}
            </button>
          </form>

          {keys.length === 0 ? (
            <p className="no-keys">No API keys yet.</p>
          ) : (
            <ul className="api-keys-list">
              {keys.map((key) => (
                <li key={key.id} className="api-key-item">
                  <div className="key-info">
                    <code className="key-prefix">{key.key_prefix}...</code>
                    {key.description && (
                      <span className="key-description">{key.description}</span>
                    )}
                    <span className="key-meta">
                      Created {new Date(key.created_at).toLocaleDateString()}
                      {key.last_used_at && (
                        <> · Last used {new Date(key.last_used_at).toLocaleDateString()}</>
                      )}
                    </span>
                  </div>
                  <button
                    onClick={() => handleDeleteKey(key)}
                    className="btn-danger"
                  >
                    Delete
                  </button>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

    </div>
  );
}
