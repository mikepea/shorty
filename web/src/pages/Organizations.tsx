/**
 * Organizations.tsx - Organizations Management Page
 *
 * This page allows users to:
 * 1. View all organizations they belong to
 * 2. Create new organizations
 * 3. Manage organization members (if admin)
 * 4. Leave organizations
 *
 * Key concepts:
 * - Fetching data on component mount with useEffect
 * - Form submission with controlled inputs
 * - Conditional rendering based on user's role in each org
 */

import { useState, useEffect } from 'react';
import { organizations as orgsApi, type Organization, type OrganizationMember } from '../api';
import { useOrganization } from '../context/OrganizationContext';
import { useAuth } from '../context/AuthContext';

export default function Organizations() {
  // ============================================================================
  // State Management
  // ============================================================================

  // Get organization context for refreshing the list after changes
  const { organizations, refreshOrgs, switchOrg, currentOrg } = useOrganization();
  const { user } = useAuth();

  // Form state for creating a new organization
  const [newOrgName, setNewOrgName] = useState('');
  const [newOrgSlug, setNewOrgSlug] = useState('');

  // State for tracking which org's members are being viewed
  const [selectedOrgId, setSelectedOrgId] = useState<number | null>(null);
  const [members, setMembers] = useState<OrganizationMember[]>([]);
  const [loadingMembers, setLoadingMembers] = useState(false);

  // Form state for adding a member
  const [newMemberEmail, setNewMemberEmail] = useState('');
  const [newMemberRole, setNewMemberRole] = useState<'admin' | 'member'>('member');

  // Loading and error states
  const [isCreating, setIsCreating] = useState(false);
  const [isAddingMember, setIsAddingMember] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // ============================================================================
  // Effects
  // ============================================================================

  /**
   * Load members when an org is selected for viewing.
   */
  useEffect(() => {
    if (selectedOrgId) {
      loadMembers(selectedOrgId);
    }
  }, [selectedOrgId]);

  // ============================================================================
  // Event Handlers
  // ============================================================================

  /**
   * Auto-generate slug from name.
   * Converts to lowercase, replaces spaces with hyphens, removes special chars.
   */
  const handleNameChange = (name: string) => {
    setNewOrgName(name);
    // Auto-generate slug from name
    const slug = name
      .toLowerCase()
      .replace(/\s+/g, '-')      // Replace spaces with hyphens
      .replace(/[^a-z0-9-]/g, '') // Remove non-alphanumeric chars (except hyphens)
      .replace(/-+/g, '-')        // Replace multiple hyphens with single
      .replace(/^-|-$/g, '');     // Remove leading/trailing hyphens
    setNewOrgSlug(slug);
  };

  /**
   * Create a new organization.
   */
  const handleCreateOrg = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newOrgName.trim() || !newOrgSlug.trim()) return;

    setIsCreating(true);
    setError('');
    setSuccess('');

    try {
      const org = await orgsApi.create(newOrgName.trim(), newOrgSlug.trim());
      setNewOrgName('');
      setNewOrgSlug('');
      setSuccess(`Organization "${org.name}" created successfully!`);
      await refreshOrgs();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create organization');
    } finally {
      setIsCreating(false);
    }
  };

  /**
   * Load members for an organization.
   */
  const loadMembers = async (orgId: number) => {
    setLoadingMembers(true);
    try {
      const memberList = await orgsApi.listMembers(orgId);
      setMembers(memberList);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load members');
    } finally {
      setLoadingMembers(false);
    }
  };

  /**
   * Add a member to an organization.
   */
  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedOrgId || !newMemberEmail.trim()) return;

    setIsAddingMember(true);
    setError('');

    try {
      await orgsApi.addMember(selectedOrgId, newMemberEmail.trim(), newMemberRole);
      setNewMemberEmail('');
      setNewMemberRole('member');
      await loadMembers(selectedOrgId);
      setSuccess('Member added successfully!');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add member');
    } finally {
      setIsAddingMember(false);
    }
  };

  /**
   * Remove a member from an organization.
   */
  const handleRemoveMember = async (userId: number) => {
    if (!selectedOrgId) return;
    if (!confirm('Are you sure you want to remove this member?')) return;

    try {
      await orgsApi.removeMember(selectedOrgId, userId);
      await loadMembers(selectedOrgId);
      setSuccess('Member removed successfully!');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to remove member');
    }
  };

  /**
   * Update a member's role.
   */
  const handleUpdateRole = async (userId: number, role: 'admin' | 'member') => {
    if (!selectedOrgId) return;

    try {
      await orgsApi.updateMember(selectedOrgId, userId, role);
      await loadMembers(selectedOrgId);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update role');
    }
  };

  /**
   * Switch to an organization.
   */
  const handleSwitchOrg = (org: Organization) => {
    switchOrg(org);
    setSuccess(`Switched to "${org.name}"`);
  };

  /**
   * Get the user's role in an organization.
   */
  const getUserRole = (org: Organization): 'admin' | 'member' | undefined => {
    return org.role;
  };

  /**
   * Check if user is admin of the selected org.
   */
  const isAdminOfSelected = (): boolean => {
    if (!selectedOrgId) return false;
    const org = organizations.find(o => o.id === selectedOrgId);
    return org?.role === 'admin';
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <div className="orgs-page">
      <div className="page-header">
        <h1>Organizations</h1>
      </div>

      {/* Display error/success messages */}
      {error && <div className="error">{error}</div>}
      {success && <div className="success">{success}</div>}

      <div className="orgs-content">
        {/* Create New Organization */}
        <section>
          <h2>Create Organization</h2>
          <form onSubmit={handleCreateOrg} className="create-org-form">
            <input
              type="text"
              placeholder="Organization name"
              value={newOrgName}
              onChange={(e) => handleNameChange(e.target.value)}
              required
            />
            <input
              type="text"
              placeholder="Slug (URL identifier)"
              value={newOrgSlug}
              onChange={(e) => setNewOrgSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
              pattern="^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$"
              title="Lowercase letters, numbers, and hyphens only. No leading/trailing hyphens."
              required
            />
            <button type="submit" disabled={isCreating}>
              {isCreating ? 'Creating...' : 'Create'}
            </button>
          </form>
        </section>

        {/* Organizations List */}
        <section>
          <h2>Your Organizations</h2>
          {organizations.length === 0 ? (
            <p className="no-results">You don't belong to any organizations yet.</p>
          ) : (
            <ul className="orgs-list">
              {organizations.map((org) => (
                <li key={org.id} className="org-item">
                  <div className="org-info">
                    <span className="org-name">
                      {org.name}
                      {org.is_global && <span className="org-global-badge">Global</span>}
                      {getUserRole(org) === 'admin' && <span className="org-role-badge">Admin</span>}
                    </span>
                    <span className="org-slug">/{org.slug}</span>
                    <span className="org-meta">
                      <span>{org.member_count} member{org.member_count !== 1 ? 's' : ''}</span>
                    </span>
                  </div>
                  <div className="org-actions">
                    {currentOrg?.id !== org.id && (
                      <button
                        className="btn-secondary"
                        onClick={() => handleSwitchOrg(org)}
                      >
                        Switch
                      </button>
                    )}
                    {currentOrg?.id === org.id && (
                      <span className="org-role-badge">Current</span>
                    )}
                    <button
                      className="btn-secondary"
                      onClick={() => setSelectedOrgId(selectedOrgId === org.id ? null : org.id)}
                    >
                      {selectedOrgId === org.id ? 'Hide Members' : 'Members'}
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>

        {/* Members Section (shown when an org is selected) */}
        {selectedOrgId && (
          <section>
            <h2>
              Members of {organizations.find(o => o.id === selectedOrgId)?.name}
            </h2>

            {/* Add member form (admin only) */}
            {isAdminOfSelected() && (
              <form onSubmit={handleAddMember} className="create-org-form" style={{ marginBottom: '1rem' }}>
                <input
                  type="email"
                  placeholder="Email address"
                  value={newMemberEmail}
                  onChange={(e) => setNewMemberEmail(e.target.value)}
                  required
                />
                <select
                  value={newMemberRole}
                  onChange={(e) => setNewMemberRole(e.target.value as 'admin' | 'member')}
                >
                  <option value="member">Member</option>
                  <option value="admin">Admin</option>
                </select>
                <button type="submit" disabled={isAddingMember}>
                  {isAddingMember ? 'Adding...' : 'Add Member'}
                </button>
              </form>
            )}

            {loadingMembers ? (
              <p>Loading members...</p>
            ) : members.length === 0 ? (
              <p className="no-results">No members found.</p>
            ) : (
              <ul className="orgs-list">
                {members.map((member) => (
                  <li key={member.id} className="org-item">
                    <div className="org-info">
                      <span className="org-name">{member.name}</span>
                      <span className="org-slug">{member.email}</span>
                      <span className="org-meta">
                        <span className="org-role-badge">{member.role}</span>
                        {member.user_id === user?.id && <span> (You)</span>}
                      </span>
                    </div>
                    {isAdminOfSelected() && member.user_id !== user?.id && (
                      <div className="org-actions">
                        <select
                          value={member.role}
                          onChange={(e) => handleUpdateRole(member.user_id, e.target.value as 'admin' | 'member')}
                        >
                          <option value="member">Member</option>
                          <option value="admin">Admin</option>
                        </select>
                        <button
                          className="btn-danger"
                          onClick={() => handleRemoveMember(member.user_id)}
                        >
                          Remove
                        </button>
                      </div>
                    )}
                  </li>
                ))}
              </ul>
            )}
          </section>
        )}
      </div>
    </div>
  );
}
