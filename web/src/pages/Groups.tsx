import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { groups as groupsApi } from '../api/client';
import type { Group } from '../api/types';

export default function Groups() {
  const [groupsList, setGroupsList] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [newGroupName, setNewGroupName] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    groupsApi.list()
      .then(setGroupsList)
      .finally(() => setIsLoading(false));
  }, []);

  const handleCreateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newGroupName.trim()) return;

    setIsCreating(true);
    setError('');

    try {
      const group = await groupsApi.create(newGroupName.trim());
      setGroupsList([...groupsList, group]);
      setNewGroupName('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create group');
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteGroup = async (group: Group) => {
    if (!confirm(`Are you sure you want to delete "${group.name}"?`)) return;

    try {
      await groupsApi.delete(group.id);
      setGroupsList(groupsList.filter((g) => g.id !== group.id));
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete group');
    }
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="groups-page">
      <header className="page-header">
        <h1>Groups</h1>
      </header>

      {error && <div className="error">{error}</div>}

      <div className="groups-content">
        <section className="create-group">
          <h2>Create New Group</h2>
          <form onSubmit={handleCreateGroup} className="create-group-form">
            <input
              type="text"
              value={newGroupName}
              onChange={(e) => setNewGroupName(e.target.value)}
              placeholder="Group name"
              required
            />
            <button type="submit" disabled={isCreating}>
              {isCreating ? 'Creating...' : 'Create Group'}
            </button>
          </form>
        </section>

        <section className="groups-list-section">
          <h2>Your Groups</h2>
          {groupsList.length === 0 ? (
            <p>No groups yet. Create one above!</p>
          ) : (
            <ul className="groups-list">
              {groupsList.map((group) => (
                <li key={group.id} className="group-item">
                  <div className="group-info">
                    <span className="group-name">{group.name}</span>
                    <span className="group-date">
                      Created {new Date(group.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="group-actions">
                    <Link to={`/links?group=${group.id}`}>View Links</Link>
                    <button
                      onClick={() => handleDeleteGroup(group)}
                      className="btn-danger"
                    >
                      Delete
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      <Link to="/" className="back-link">Back to Dashboard</Link>
    </div>
  );
}
