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

  useEffect(() => {
    setIsLoading(true);
    Promise.all([
      linksApi.search({
        q: searchParams.get('q') || undefined,
        tag: currentTag || undefined,
      }),
      tagsApi.list(),
    ])
      .then(([linksData, tagsData]) => {
        setLinksList(linksData);
        setAllTags(tagsData);
      })
      .finally(() => setIsLoading(false));
  }, [searchParams, currentTag]);

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
      <header className="links-header">
        <h1>Links</h1>
        <div className="links-toolbar">
          <form onSubmit={handleSearch} className="search-form">
            <input
              type="search"
              placeholder="Search links..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button type="submit">Search</button>
          </form>

          {allTags.length > 0 && (
            <div className="tag-filter">
              <select
                value={currentTag || ''}
                onChange={(e) => {
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
                {allTags.map((tag) => (
                  <option key={tag.id} value={tag.name}>
                    {tag.name} ({tag.link_count})
                  </option>
                ))}
              </select>
            </div>
          )}

          <Link to="/links/new" className="btn-primary">Add Link</Link>
        </div>
      </header>

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
                <th>Created</th>
                <th>Tags</th>
              </tr>
            </thead>
            <tbody>
              {linksList.map((link) => (
                <tr
                  key={link.id}
                  onClick={() => handleRowClick(link.slug)}
                >
                  <td className="links-table-title">
                    <span className="link-title-text">{link.title || link.url}</span>
                  </td>
                  <td className="links-table-slug">/{link.slug}</td>
                  <td className="links-table-clicks">{link.click_count}</td>
                  <td className="links-table-date">
                    {new Date(link.created_at).toLocaleDateString()}
                  </td>
                  <td className="links-table-tags">
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
        <div className="links-count">
          {linksList.length} link{linksList.length !== 1 ? 's' : ''}
        </div>
      </main>
    </div>
  );
}
