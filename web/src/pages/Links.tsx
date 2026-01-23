/**
 * Links.tsx - Links List Page Component
 *
 * This page displays all links with search and tag filtering.
 * It demonstrates several important React patterns:
 *
 * Key concepts:
 * - Data fetching with useEffect
 * - URL-based state with useSearchParams (for shareable/bookmarkable URLs)
 * - Rendering lists with .map()
 * - Conditional rendering for empty states
 * - Event handling for search and navigation
 * - Promise.all for parallel API requests
 */

import { useState, useEffect } from 'react';
import { Link, useSearchParams, useNavigate } from 'react-router-dom';
import { links as linksApi, tags as tagsApi } from '../api/client';
import type { Link as LinkType, Tag } from '../api/types';

export default function Links() {
  // ============================================================================
  // Hooks
  // ============================================================================

  // useNavigate returns a function to programmatically change routes
  const navigate = useNavigate();

  /**
   * useSearchParams is like useState but syncs with the URL query string.
   * If the URL is /links?q=hello&tag=work, searchParams.get('q') returns 'hello'.
   * This makes the search/filter state shareable and bookmarkable!
   */
  const [searchParams, setSearchParams] = useSearchParams();

  // ============================================================================
  // State
  // ============================================================================

  // The list of links fetched from the API
  const [linksList, setLinksList] = useState<LinkType[]>([]);

  // All available tags for the filter dropdown
  const [allTags, setAllTags] = useState<Tag[]>([]);

  // Loading state - true while fetching data
  const [isLoading, setIsLoading] = useState(true);

  // Local state for the search input (synced to URL on submit)
  // Initialize with the current URL parameter value, or empty string
  const [searchQuery, setSearchQuery] = useState(searchParams.get('q') || '');

  // Get current tag filter from URL (computed, not state)
  const currentTag = searchParams.get('tag');

  // ============================================================================
  // Effects - Data Fetching
  // ============================================================================

  /**
   * Fetch links and tags whenever the URL search params change.
   *
   * The dependency array [searchParams, currentTag] means this effect runs:
   * 1. On initial mount
   * 2. Whenever searchParams changes (user searches or filters)
   *
   * Promise.all() runs multiple promises in parallel and waits for all to complete.
   * This is more efficient than awaiting each one sequentially.
   */
  useEffect(() => {
    setIsLoading(true);

    // Fetch links and tags in parallel
    Promise.all([
      linksApi.search({
        q: searchParams.get('q') || undefined,  // undefined omits the param
        tag: currentTag || undefined,
      }),
      tagsApi.list(),
    ])
      .then(([linksData, tagsData]) => {
        // Array destructuring: [first, second] = Promise.all result
        setLinksList(linksData);
        setAllTags(tagsData);
      })
      .finally(() => setIsLoading(false));
  }, [searchParams, currentTag]);

  // ============================================================================
  // Event Handlers
  // ============================================================================

  /**
   * Handle search form submission.
   * Updates the URL query string, which triggers the useEffect to refetch data.
   */
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();

    // Create a new URLSearchParams from current params
    const params = new URLSearchParams(searchParams);

    // Add or remove the 'q' parameter based on search input
    if (searchQuery) {
      params.set('q', searchQuery);
    } else {
      params.delete('q');  // Remove empty search params from URL
    }

    // Update the URL, which triggers the useEffect
    setSearchParams(params);
  };

  /**
   * Handle clicking on a table row - navigate to the link detail page.
   */
  const handleRowClick = (slug: string) => {
    navigate(`/links/${slug}`);
  };

  // ============================================================================
  // Render
  // ============================================================================

  // Early return for loading state - a common pattern
  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="links-page">
      {/* Page header with search and filter controls */}
      <header className="links-header">
        <h1>Links</h1>
        <div className="links-toolbar">
          {/* Search form */}
          <form onSubmit={handleSearch} className="search-form">
            <input
              type="search"
              placeholder="Search links..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button type="submit">Search</button>
          </form>

          {/* Tag filter dropdown - only shown if there are tags */}
          {allTags.length > 0 && (
            <div className="tag-filter">
              <select
                value={currentTag || ''}
                onChange={(e) => {
                  // Update URL params directly in the onChange handler
                  const params = new URLSearchParams(searchParams);
                  if (e.target.value) {
                    params.set('tag', e.target.value);
                  } else {
                    params.delete('tag');
                  }
                  setSearchParams(params);
                }}
              >
                <option value="">All tags</option>
                {/*
                  Map over tags to create an option for each.
                  The map() function transforms an array by applying a function
                  to each element and returning a new array.
                */}
                {allTags.map((tag) => (
                  <option key={tag.id} value={tag.name}>
                    {tag.name} ({tag.link_count})
                  </option>
                ))}
              </select>
            </div>
          )}

          {/* Link to add new link page */}
          <Link to="/links/new" className="btn-primary">Add Link</Link>
        </div>
      </header>

      {/* Main content area */}
      <main className="links-main">
        {/*
          Ternary operator for conditional rendering:
          condition ? (render if true) : (render if false)

          Show "no results" message if list is empty, otherwise show the table.
        */}
        {linksList.length === 0 ? (
          <p className="no-results">No links found.</p>
        ) : (
          <table className="links-table">
            <thead>
              <tr>
                <th>Title</th>
                <th>Slug</th>
                <th>Clicks</th>
                <th>Created</th>
                <th>Tags</th>
              </tr>
            </thead>
            <tbody>
              {/*
                Map over links to create a table row for each.
                Each row needs a unique "key" prop - React uses this to
                efficiently update the DOM when the list changes.
              */}
              {linksList.map((link) => (
                <tr
                  key={link.id}
                  onClick={() => handleRowClick(link.slug)}
                >
                  <td className="links-table-title">
                    {/* Use title if available, fallback to URL */}
                    <span className="link-title-text">{link.title || link.url}</span>
                  </td>
                  <td className="links-table-slug">/{link.slug}</td>
                  <td className="links-table-clicks">{link.click_count}</td>
                  <td className="links-table-date">
                    {/* Convert ISO date string to localized date format */}
                    {new Date(link.created_at).toLocaleDateString()}
                  </td>
                  <td className="links-table-tags">
                    {/*
                      Show first 3 tags, then "+N" if there are more.
                      The ?. (optional chaining) handles case where tags is undefined.
                      slice(0, 3) gets first 3 elements without modifying original array.
                    */}
                    {link.tags && link.tags.slice(0, 3).map((tag) => (
                      <span key={tag.id} className="tag">{tag.name}</span>
                    ))}
                    {link.tags && link.tags.length > 3 && (
                      <span className="tag-more">+{link.tags.length - 3}</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}

        {/* Summary count at the bottom */}
        <div className="links-count">
          {/* Template literal with conditional plural */}
          {linksList.length} link{linksList.length !== 1 ? 's' : ''}
        </div>
      </main>
    </div>
  );
}
