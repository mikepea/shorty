/**
 * Layout.tsx - Page Layout Component
 *
 * This component provides the consistent layout structure for all authenticated pages.
 * It renders the sidebar navigation alongside the main page content.
 *
 * Key concept: "Composition"
 * React components can contain other components. Layout "composes" Sidebar with
 * whatever page content is passed as children. This avoids duplicating the sidebar
 * code in every page component.
 *
 * Visual structure:
 * ┌─────────────────────────────────────────┐
 * │ ┌──────────┐ ┌────────────────────────┐ │
 * │ │          │ │                        │ │
 * │ │ Sidebar  │ │    Main Content        │ │
 * │ │          │ │    (children)          │ │
 * │ │          │ │                        │ │
 * │ └──────────┘ └────────────────────────┘ │
 * └─────────────────────────────────────────┘
 */

import Sidebar from './Sidebar';

/**
 * Props interface for Layout component.
 *
 * React.ReactNode is a type that accepts anything that can be rendered:
 * - JSX elements (<div>...</div>)
 * - Strings and numbers
 * - Arrays of the above
 * - null or undefined
 */
interface LayoutProps {
  children: React.ReactNode;
}

/**
 * Layout component - wraps page content with sidebar navigation.
 *
 * Usage in App.tsx:
 *   <Layout>
 *     <Dashboard />   ← This becomes "children"
 *   </Layout>
 *
 * The "export default" means this is the main export from this file.
 * Other files can import it as: import Layout from './components/Layout';
 */
export default function Layout({ children }: LayoutProps) {
  return (
    // CSS class "app-layout" uses flexbox to put sidebar and content side by side
    <div className="app-layout">
      {/* Sidebar component - navigation and user info */}
      <Sidebar />

      {/* Main content area - receives the page component passed as children */}
      <main className="main-content">
        {children}
      </main>
    </div>
  );
}
