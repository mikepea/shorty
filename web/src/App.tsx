/**
 * App.tsx - Main Application Component
 *
 * This file sets up:
 * 1. Routing - Which component shows for which URL (e.g., /login shows Login)
 * 2. Authentication - Protecting pages that require login
 * 3. Layout - Wrapping pages in a consistent layout with sidebar
 *
 * Key concepts:
 * - BrowserRouter: Enables URL-based navigation without page reloads
 * - Routes/Route: Define which component renders for each URL path
 * - Navigate: Redirects users to a different page programmatically
 * - Children props: Components can receive other components as "children"
 */

import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { OrganizationProvider } from './context/OrganizationContext';
import Layout from './components/Layout';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Links from './pages/Links';
import LinkDetail from './pages/LinkDetail';
import AddLink from './pages/AddLink';
import EditLink from './pages/EditLink';
import Groups from './pages/Groups';
import Organizations from './pages/Organizations';
import Settings from './pages/Settings';
import Admin from './pages/Admin';
import SSOCallback from './pages/SSOCallback';
import './App.css';

/**
 * ProtectedRoute - A wrapper that requires authentication
 *
 * This is a common pattern in React apps. It wraps pages that should only
 * be accessible to logged-in users. If not logged in, it redirects to /login.
 *
 * The { children } syntax is "destructuring" - extracting the children prop
 * from the props object. children is whatever you put between the opening
 * and closing tags: <ProtectedRoute>THIS IS CHILDREN</ProtectedRoute>
 */
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  // useAuth() is a "hook" that gives us access to authentication state.
  // Hooks are functions that let components use React features like state.
  const { user, isLoading } = useAuth();

  // While checking if the user is logged in, show a loading message.
  // This prevents a flash of the login page before we know the auth state.
  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  // If there's no user (not logged in), redirect to the login page.
  // The "replace" prop means this won't add to browser history,
  // so the back button won't return to this protected page.
  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // User is logged in! Wrap the page content in the Layout component
  // which provides the sidebar and consistent styling.
  return <Layout>{children}</Layout>;
}

/**
 * PublicRoute - A wrapper for pages that should NOT be accessible when logged in
 *
 * This redirects logged-in users away from login/register pages.
 * If you're already logged in, going to /login redirects you to the dashboard.
 */
function PublicRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  // If already logged in, redirect to home page (dashboard)
  if (user) {
    return <Navigate to="/" replace />;
  }

  // The <>{children}</> syntax is a "Fragment" - it lets us return
  // multiple elements without adding an extra div to the DOM.
  return <>{children}</>;
}

/**
 * AppRoutes - Defines all the URL routes in the application
 *
 * Each Route maps a URL path to a component. The path can include:
 * - Static segments: "/login" matches exactly "/login"
 * - Dynamic segments: "/links/:slug" matches "/links/abc", "/links/xyz", etc.
 *   The :slug part becomes a variable you can access in the component.
 */
function AppRoutes() {
  return (
    <Routes>
      {/* Public routes - accessible without login */}
      <Route
        path="/login"
        element={
          <PublicRoute>
            <Login />
          </PublicRoute>
        }
      />
      <Route
        path="/register"
        element={
          <PublicRoute>
            <Register />
          </PublicRoute>
        }
      />

      {/* Protected routes - require authentication */}
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/links"
        element={
          <ProtectedRoute>
            <Links />
          </ProtectedRoute>
        }
      />
      {/* Note: /links/new must come BEFORE /links/:slug
          otherwise "new" would be treated as a slug */}
      <Route
        path="/links/new"
        element={
          <ProtectedRoute>
            <AddLink />
          </ProtectedRoute>
        }
      />
      {/* :slug is a URL parameter - accessible via useParams() hook */}
      <Route
        path="/links/:slug"
        element={
          <ProtectedRoute>
            <LinkDetail />
          </ProtectedRoute>
        }
      />
      <Route
        path="/links/:slug/edit"
        element={
          <ProtectedRoute>
            <EditLink />
          </ProtectedRoute>
        }
      />
      <Route
        path="/groups"
        element={
          <ProtectedRoute>
            <Groups />
          </ProtectedRoute>
        }
      />
      <Route
        path="/organizations"
        element={
          <ProtectedRoute>
            <Organizations />
          </ProtectedRoute>
        }
      />
      <Route
        path="/settings"
        element={
          <ProtectedRoute>
            <Settings />
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin"
        element={
          <ProtectedRoute>
            <Admin />
          </ProtectedRoute>
        }
      />

      {/* SSO callback doesn't use ProtectedRoute because it handles
          its own authentication flow (receives token from OAuth provider) */}
      <Route path="/sso/callback" element={<SSOCallback />} />
    </Routes>
  );
}

/**
 * App - The root component of the application
 *
 * This sets up the "providers" that make features available throughout the app:
 * - BrowserRouter: Enables routing (URL-based navigation)
 * - AuthProvider: Makes authentication state available everywhere via useAuth()
 *
 * Providers use React's "Context" feature to pass data down the component tree
 * without having to pass props through every intermediate component.
 */
function App() {
  return (
    // BrowserRouter must wrap everything that uses routing
    <BrowserRouter>
      {/* AuthProvider must wrap everything that uses useAuth() */}
      <AuthProvider>
        {/* OrganizationProvider manages current org state, must be inside AuthProvider
            because it needs to know when the user logs in/out */}
        <OrganizationProvider>
          <AppRoutes />
        </OrganizationProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}

// "export default" makes this the main thing imported when you do:
// import App from './App'
export default App;
