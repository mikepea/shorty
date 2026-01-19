import { useState, useEffect } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { links as linksApi, tags as tagsApi } from '../api/client';
import type { Link as LinkType, Tag } from '../api/types';

export default function Links() {
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

  const handleDelete = async (slug: string) => {
    if (!confirm('Are you sure you want to delete this link?')) return;

    try {
      await linksApi.delete(slug);
      setLinksList(linksList.filter((l) => l.slug !== slug));
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete link');
    }
  };

  const toggleUnread = async (link: LinkType) => {
    try {
      const updated = await linksApi.update(link.slug, {
        is_unread: !link.is_unread,
      });
      setLinksList(linksList.map((l) => (l.id === link.id ? updated : l)));
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update link');
    }
  };

  const copyShortUrl = (slug: string) => {
    const url = `${window.location.origin}/${slug}`;
    navigator.clipboard.writeText(url);
    alert('Short URL copied to clipboard!');
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
            <ul className="link-list">
              {linksList.map((link) => (
                <li key={link.id} className={`link-item ${link.is_unread ? 'unread' : ''}`}>
                  <div className="link-header">
                    <a
                      href={link.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="link-title"
                    >
                      {link.title || link.url}
                    </a>
                    <span className="link-visibility">
                      {link.is_public ? 'Public' : 'Private'}
                    </span>
                  </div>

                  <div className="link-meta">
                    <span className="link-slug" onClick={() => copyShortUrl(link.slug)}>
                      /{link.slug}
                    </span>
                    <span className="link-clicks">{link.click_count} clicks</span>
                    <span className="link-date">
                      {new Date(link.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  {link.description && (
                    <p className="link-description">{link.description}</p>
                  )}

                  {link.tags && link.tags.length > 0 && (
                    <div className="link-tags">
                      {link.tags.map((tag) => (
                        <span key={tag.id} className="tag">
                          {tag.name}
                        </span>
                      ))}
                    </div>
                  )}

                  <div className="link-actions">
                    <button onClick={() => toggleUnread(link)}>
                      {link.is_unread ? 'Mark Read' : 'Mark Unread'}
                    </button>
                    <Link to={`/links/${link.slug}/edit`}>Edit</Link>
                    <button onClick={() => handleDelete(link.slug)} className="btn-danger">
                      Delete
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </main>
      </div>

      <Link to="/" className="back-link">Back to Dashboard</Link>
    </div>
  );
}
