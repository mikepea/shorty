/**
 * OrgSwitcher.tsx - Organization Switcher Component
 *
 * This component allows users to switch between their organizations.
 * It displays the current organization and provides a dropdown to select another.
 *
 * Key concepts:
 * - Dropdown menus: Toggle visibility with state
 * - useRef: Access DOM elements directly (for click-outside detection)
 * - useEffect cleanup: Remove event listeners when component unmounts
 */

import { useState, useRef, useEffect } from 'react';
import { useOrganization } from '../context/OrganizationContext';

/**
 * OrgSwitcher - Dropdown for switching between organizations
 *
 * Displays the current org name with a dropdown arrow.
 * Clicking opens a menu of all available organizations.
 * Clicking an org switches to it and closes the menu.
 * Clicking outside the menu closes it.
 */
export default function OrgSwitcher() {
  // Get organization state from context
  const { organizations, currentOrg, isLoading, switchOrg } = useOrganization();

  // State to track if the dropdown is open
  // useState(false) initializes isOpen to false
  const [isOpen, setIsOpen] = useState(false);

  // useRef creates a reference to a DOM element
  // We use this to detect clicks outside the dropdown
  const dropdownRef = useRef<HTMLDivElement>(null);

  /**
   * Close dropdown when clicking outside.
   *
   * useEffect with cleanup function:
   * - The effect adds an event listener when component mounts
   * - The return function (cleanup) removes it when component unmounts
   * - This prevents memory leaks and stale event handlers
   */
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      // If the click was outside the dropdown, close it
      // dropdownRef.current is the DOM element
      // .contains() checks if the clicked element is inside it
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    // Only add listener when dropdown is open (optimization)
    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    // Cleanup function - runs when effect re-runs or component unmounts
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]); // Re-run effect when isOpen changes

  // Don't render anything while loading or if no orgs
  if (isLoading) {
    return (
      <div className="org-switcher org-switcher--loading">
        Loading...
      </div>
    );
  }

  if (organizations.length === 0) {
    return null;
  }

  // If only one org, just show the name (no dropdown needed)
  if (organizations.length === 1) {
    return (
      <div className="org-switcher org-switcher--single">
        <span className="org-switcher__name">{currentOrg?.name}</span>
      </div>
    );
  }

  return (
    <div className="org-switcher" ref={dropdownRef}>
      {/* Toggle button - clicking opens/closes the dropdown */}
      <button
        className="org-switcher__toggle"
        onClick={() => setIsOpen(!isOpen)}
        aria-expanded={isOpen}
        aria-haspopup="listbox"
      >
        <span className="org-switcher__name">{currentOrg?.name || 'Select Organization'}</span>
        <span className="org-switcher__arrow">{isOpen ? '▲' : '▼'}</span>
      </button>

      {/* Dropdown menu - only rendered when isOpen is true */}
      {isOpen && (
        <ul className="org-switcher__menu" role="listbox">
          {organizations.map((org) => (
            <li key={org.id}>
              <button
                className={`org-switcher__option ${org.id === currentOrg?.id ? 'org-switcher__option--selected' : ''}`}
                onClick={() => {
                  switchOrg(org);
                  setIsOpen(false);
                }}
                role="option"
                aria-selected={org.id === currentOrg?.id}
              >
                {org.name}
                {org.is_global && <span className="org-switcher__badge">Global</span>}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
