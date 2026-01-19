import { useState, useEffect } from 'react';
import { Link, useSearchParams, useNavigate } from 'react-router-dom';
import { links as linksApi, tags as tagsApi } from '../api/client';
import type { Link as LinkType, Tag } from '../api/types';

export default function Links() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [linksList, setLinksList] = useState<LinkType[]>([]);
  const [allTags, setAllTags] = useState<Tag[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState(searchParams.get('q') || '');

  const currentTag = searchParams.get('tag');
  const showUnread = searchParams.get('unread') === 'true';

  useEffect(() => {
    setIsLoading(true);
    Promise.all([
      linksApi.search({
        q: searchParams.get('q') || undefined,
        tag: currentTag || undefined,
        unread: showUnread,
      }),
      tagsApi.list(),
    ])
      .then(([linksData, tagsData]) => {
        setLinksList(linksData);
        setAllTags(tagsData);
      })
      .finally(() => setIsLoading(false));
  }, [searchParams, currentTag, showUnread]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    const params = new URLSearchParams(searchParams);
    if (searchQuery) {
      params.set('q', searchQuery);
    } else {
      params.delete('q');
    }
    setSearchParams(params);
  };

  const handleRowClick = (slug: string) => {
    navigate(`/links/${slug}`);
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="links-page">
      <header className="page-header">
        <h1>Links</h1>
        <Link to="/links/new" className="btn-primary">Add Link</Link>
      </header>

      <div className="links-layout">
        <aside className="links-sidebar">
          <form onSubmit={handleSearch} className="search-form">
            <input
              type="search"
              placeholder="Search links..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button type="submit">Search</button>
          </form>

          <div className="filter-section">
            <h3>Filter</h3>
            <label>
              <input
                type="checkbox"
                checked={showUnread}
                onChange={(e) => {
                  const params = new URLSearchParams(searchParams);
                  if (e.target.checked) {
                    params.set('unread', 'true');
                  } else {
                    params.delete('unread');
                  }
                  setSearchParams(params);
                }}
              />
              Unread only
            </label>
          </div>

          <div className="tags-section">
            <h3>Tags</h3>
            {currentTag && (
              <button
                onClick={() => {
                  const params = new URLSearchParams(searchParams);
                  params.delete('tag');
                  setSearchParams(params);
                }}
                className="clear-tag"
              >
                Clear: {currentTag}
              </button>
            )}
            <ul className="tag-list">
              {allTags.map((tag) => (
                <li key={tag.id}>
                  <button
                    onClick={() => {
                      const params = new URLSearchParams(searchParams);
                      params.set('tag', tag.name);
                      setSearchParams(params);
                    }}
                    className={currentTag === tag.name ? 'active' : ''}
                  >
                    {tag.name} ({tag.link_count})
                  </button>
                </li>
              ))}
            </ul>
          </div>
        </aside>

        <main className="links-main">
          {linksList.length === 0 ? (
            <p className="no-results">No links found.</p>
          ) : (
            <table className="links-table">
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Slug</th>
                  <th>Clicks</th>
                  <th>Tags</th>
                  <th>Created</th>
                </tr>
              </thead>
              <tbody>
                {linksList.map((link) => (
                  <tr
                    key={link.id}
                    onClick={() => handleRowClick(link.slug)}
                    className={link.is_unread ? 'unread' : ''}
                  >
                    <td className="links-table-title">
                      <span className="link-title-text">{link.title || link.url}</span>
                      {link.is_unread && <span className="unread-badge">Unread</span>}
                    </td>
                    <td className="links-table-slug">/{link.slug}</td>
                    <td className="links-table-clicks">{link.click_count}</td>
                    <td className="links-table-tags">
                      {link.tags && link.tags.slice(0, 3).map((tag) => (
                        <span key={tag.id} className="tag">{tag.name}</span>
                      ))}
                      {link.tags && link.tags.length > 3 && (
                        <span className="tag-more">+{link.tags.length - 3}</span>
                      )}
                    </td>
                    <td className="links-table-date">
                      {new Date(link.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          <div className="links-count">
            {linksList.length} link{linksList.length !== 1 ? 's' : ''}
          </div>
        </main>
      </div>
    </div>
  );
}
