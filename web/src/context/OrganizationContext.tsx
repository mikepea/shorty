/**
 * OrganizationContext.tsx - Organization State Management
 *
 * This file implements the "Context" pattern for sharing organization state
 * across the entire application. It manages which organization the user is
 * currently working in.
 *
 * Key concepts:
 * - Multi-tenancy: Users can belong to multiple organizations
 * - Current organization: The org context for all API requests
 * - Organization switching: Users can switch between their orgs
 *
 * How it works:
 * 1. OrganizationProvider wraps the app and holds the organization state
 * 2. On mount, it fetches the user's organizations from the API
 * 3. It restores the last selected org from localStorage (if still valid)
 * 4. Any component can call useOrganization() to get/set the current org
 */

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import type { Organization } from '../api/types';
import { organizations as orgsApi, setCurrentOrgId, getCurrentOrgId } from '../api/client';
import { useAuth } from './AuthContext';

/**
 * TypeScript interface defining what's available in the organization context.
 */
interface OrganizationContextType {
  organizations: Organization[];       // All organizations the user belongs to
  currentOrg: Organization | null;     // Currently selected organization
  isLoading: boolean;                  // True while fetching organizations
  error: string | null;                // Error message if fetch failed
  switchOrg: (org: Organization) => void;  // Function to switch organizations
  refreshOrgs: () => Promise<void>;    // Function to refresh the org list
}

/**
 * Create the context with undefined as the default value.
 */
const OrganizationContext = createContext<OrganizationContextType | undefined>(undefined);

/**
 * OrganizationProvider - The component that provides org state to the app
 *
 * This component:
 * 1. Fetches the user's organizations when they log in
 * 2. Restores the last selected org from localStorage
 * 3. Provides functions to switch between organizations
 * 4. Syncs the current org with the API client
 */
export function OrganizationProvider({ children }: { children: ReactNode }) {
  // Get auth state to know when user logs in/out
  const { user, token } = useAuth();

  // State for organizations list
  const [organizations, setOrganizations] = useState<Organization[]>([]);

  // State for current organization
  const [currentOrg, setCurrentOrg] = useState<Organization | null>(null);

  // Loading and error states
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetch the user's organizations from the API.
   * useCallback memoizes the function so it doesn't change on every render.
   */
  const fetchOrganizations = useCallback(async () => {
    if (!token) {
      // Not logged in, clear everything
      setOrganizations([]);
      setCurrentOrg(null);
      setCurrentOrgId(null);
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const orgs = await orgsApi.list();
      setOrganizations(orgs);

      // Try to restore the previously selected org
      const savedOrgId = getCurrentOrgId();
      let orgToSelect: Organization | null = null;

      if (savedOrgId) {
        // Find the saved org in the list (it might have been removed)
        orgToSelect = orgs.find(o => o.id === savedOrgId) || null;
      }

      // If no saved org or it's no longer valid, select the global org or first org
      if (!orgToSelect && orgs.length > 0) {
        // Prefer the global org as default
        orgToSelect = orgs.find(o => o.is_global) || orgs[0];
      }

      if (orgToSelect) {
        setCurrentOrg(orgToSelect);
        setCurrentOrgId(orgToSelect.id);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load organizations');
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  /**
   * Load organizations when user logs in or token changes.
   */
  useEffect(() => {
    fetchOrganizations();
  }, [fetchOrganizations]);

  /**
   * Clear org state when user logs out.
   */
  useEffect(() => {
    if (!user) {
      setOrganizations([]);
      setCurrentOrg(null);
      setCurrentOrgId(null);
    }
  }, [user]);

  /**
   * Switch to a different organization.
   * Updates local state and syncs with the API client.
   */
  const switchOrg = useCallback((org: Organization) => {
    setCurrentOrg(org);
    setCurrentOrgId(org.id);
  }, []);

  /**
   * Refresh the organizations list (e.g., after creating a new org).
   */
  const refreshOrgs = useCallback(async () => {
    await fetchOrganizations();
  }, [fetchOrganizations]);

  return (
    <OrganizationContext.Provider value={{
      organizations,
      currentOrg,
      isLoading,
      error,
      switchOrg,
      refreshOrgs,
    }}>
      {children}
    </OrganizationContext.Provider>
  );
}

/**
 * useOrganization - Custom hook for accessing organization context
 *
 * Usage in any component:
 *   const { currentOrg, organizations, switchOrg } = useOrganization();
 */
export function useOrganization() {
  const context = useContext(OrganizationContext);

  if (context === undefined) {
    throw new Error('useOrganization must be used within an OrganizationProvider');
  }

  return context;
}
