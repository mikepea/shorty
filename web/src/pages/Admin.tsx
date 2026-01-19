import { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { admin } from '../api/client';
import type { AdminStats, AdminUser } from '../api/types';

export default function Admin() {
  const { user } = useAuth();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState('');
  const [error, setError] = useState('');

  // Redirect non-admins
  if (user && user.system_role !== 'admin') {
    return <Navigate to="/" replace />;
  }

  useEffect(() => {
    setIsLoading(true);
    Promise.all([
      admin.getStats(),
      admin.listUsers(),
    ])
      .then(([statsData, usersData]) => {
        setStats(statsData);
        setUsers(usersData);
      })
      .catch(() => setError('Failed to load admin data'))
      .finally(() => setIsLoading(false));
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    admin.listUsers({ q: searchQuery || undefined, role: roleFilter || undefined })
      .then(setUsers)
      .catch(() => setError('Failed to search users'))
      .finally(() => setIsLoading(false));
  };

  const handleRoleChange = async (targetUser: AdminUser, newRole: 'admin' | 'user') => {
    if (targetUser.id === user?.id) {
      alert('Cannot change your own role');
      return;
    }
    try {
      const updated = await admin.updateUser(targetUser.id, { system_role: newRole });
      setUsers(users.map(u => u.id === updated.id ? updated : u));
      // Refresh stats if role changed
      const newStats = await admin.getStats();
      setStats(newStats);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update user');
    }
  };

  const handleDeleteUser = async (targetUser: AdminUser) => {
    if (targetUser.id === user?.id) {
      alert('Cannot delete yourself');
      return;
    }
    if (!confirm(`Are you sure you want to delete user "${targetUser.name}"? This will also delete all their links and data.`)) {
      return;
    }
    try {
      await admin.deleteUser(targetUser.id);
      setUsers(users.filter(u => u.id !== targetUser.id));
      // Refresh stats
      const newStats = await admin.getStats();
      setStats(newStats);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete user');
    }
  };

  if (isLoading && !stats) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="admin-page">
      <header className="page-header">
        <h1>Admin Dashboard</h1>
      </header>

      {error && <div className="error">{error}</div>}

      {stats && (
        <section className="stats-section">
          <h2>System Statistics</h2>
          <div className="stats-grid">
            <div className="stat-card">
              <span className="stat-value">{stats.total_users}</span>
              <span className="stat-label">Total Users</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.admin_users}</span>
              <span className="stat-label">Admins</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.total_links}</span>
              <span className="stat-label">Total Links</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.total_clicks}</span>
              <span className="stat-label">Total Clicks</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.total_groups}</span>
              <span className="stat-label">Groups</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.total_tags}</span>
              <span className="stat-label">Tags</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.public_links}</span>
              <span className="stat-label">Public Links</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.private_links}</span>
              <span className="stat-label">Private Links</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.unread_links}</span>
              <span className="stat-label">Unread Links</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">{stats.active_api_keys}</span>
              <span className="stat-label">API Keys</span>
            </div>
          </div>
        </section>
      )}

      <section className="users-section">
        <h2>User Management</h2>

        <form onSubmit={handleSearch} className="search-form">
          <input
            type="search"
            placeholder="Search users..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          <select
            value={roleFilter}
            onChange={(e) => setRoleFilter(e.target.value)}
          >
            <option value="">All Roles</option>
            <option value="admin">Admins</option>
            <option value="user">Users</option>
          </select>
          <button type="submit">Search</button>
        </form>

        {users.length === 0 ? (
          <p className="no-results">No users found.</p>
        ) : (
          <table className="users-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Links</th>
                <th>Groups</th>
                <th>Joined</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id} className={u.id === user?.id ? 'current-user' : ''}>
                  <td>{u.name}</td>
                  <td>{u.email}</td>
                  <td>
                    <select
                      value={u.system_role}
                      onChange={(e) => handleRoleChange(u, e.target.value as 'admin' | 'user')}
                      disabled={u.id === user?.id}
                    >
                      <option value="user">User</option>
                      <option value="admin">Admin</option>
                    </select>
                  </td>
                  <td>{u.link_count}</td>
                  <td>{u.group_count}</td>
                  <td>{new Date(u.created_at).toLocaleDateString()}</td>
                  <td>
                    <button
                      onClick={() => handleDeleteUser(u)}
                      className="btn-danger"
                      disabled={u.id === user?.id}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

    </div>
  );
}
