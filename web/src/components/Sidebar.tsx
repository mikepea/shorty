/**
 * Sidebar.tsx - Navigation Sidebar Component
 *
 * This component renders the left sidebar with navigation links and user info.
 * It appears on all authenticated pages (via the Layout component).
 *
 * Key concepts:
 * - NavLink: A special link component from react-router-dom that knows when
 *   it's "active" (the current URL matches its "to" prop)
 * - Conditional rendering: {condition && <Element />} only renders if condition is true
 * - Optional chaining: user?.name safely accesses name even if user is null
 */

import { NavLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import OrgSwitcher from './OrgSwitcher';

/**
 * Sidebar component - displays navigation and user info.
 *
 * Structure:
 * ┌─────────────────┐
 * │     Shorty      │  ← Header/logo
 * ├─────────────────┤
 * │ Dashboard       │
 * │ Links           │  ← Navigation links
 * │ Add Link        │     (NavLink adds "active" class
 * │ Groups          │      to current page's link)
 * │ Settings        │
 * │ Admin (if admin)│
 * ├─────────────────┤
 * │ User Name       │
 * │ user@email.com  │  ← User info section
 * │ [Logout]        │
 * └─────────────────┘
 */
export default function Sidebar() {
  // Get user data and logout function from auth context
  // Destructuring: { user, logout } extracts these from the context object
  const { user, logout } = useAuth();

  return (
    <aside className="sidebar">
      {/* App branding/header */}
      <div className="sidebar-header">
        <h1>Shorty</h1>
      </div>

      {/* Organization switcher - allows users to switch between their orgs */}
      <OrgSwitcher />

      {/* Navigation links */}
      <nav className="sidebar-nav">
        {/*
          NavLink is like <a> but for React Router.
          - "to" prop: the URL path to navigate to
          - "end" prop: only match exact path (so "/" doesn't match "/links")
          - Automatically adds "active" CSS class when the route matches
        */}
        <NavLink to="/" end>Dashboard</NavLink>
        <NavLink to="/links">Links</NavLink>
        <NavLink to="/links/new">Add Link</NavLink>
        <NavLink to="/groups">Groups</NavLink>
        <NavLink to="/organizations">Organizations</NavLink>
        <NavLink to="/settings">Settings</NavLink>

        {/*
          Conditional rendering with &&:
          If user?.system_role === 'admin' is true, render the Admin link.
          If false (or user is null), render nothing.

          The ?. is "optional chaining" - if user is null, it returns
          undefined instead of throwing an error.
        */}
        {user?.system_role === 'admin' && <NavLink to="/admin">Admin</NavLink>}
      </nav>

      {/* User info and logout */}
      <div className="sidebar-user">
        {/* Display user's name and email */}
        <span className="sidebar-user-name">{user?.name}</span>
        <span className="sidebar-user-email">{user?.email}</span>

        {/*
          Logout button calls the logout function from useAuth().
          onClick is an event handler - when the button is clicked,
          it calls the logout function.
        */}
        <button onClick={logout} className="sidebar-logout">Logout</button>
      </div>
    </aside>
  );
}
