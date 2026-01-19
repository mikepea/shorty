import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { links as linksApi } from '../api/client';
import type { Link as LinkType } from '../api/types';

export default function LinkDetail() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const [link, setLink] = useState<LinkType | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!slug) return;

    linksApi.get(slug)
      .then(setLink)
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load link'))
      .finally(() => setIsLoading(false));
  }, [slug]);

  const handleDelete = async () => {
    if (!slug || !confirm('Are you sure you want to delete this link?')) return;

    try {
      await linksApi.delete(slug);
      navigate('/links');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete link');
    }
  };

  const toggleUnread = async () => {
    if (!link || !slug) return;

    try {
      const updated = await linksApi.update(slug, {
        is_unread: !link.is_unread,
      });
      setLink(updated);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update link');
    }
  };

  const copyShortUrl = () => {
    if (!link) return;
    const url = `${window.location.origin}/${link.slug}`;
    navigator.clipboard.writeText(url);
    alert('Short URL copied to clipboard!');
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  if (error || !link) {
    return (
      <div className="link-detail-page">
        <header className="page-header">
          <h1>Link Not Found</h1>
          <Link to="/links" className="btn-secondary">Back to Links</Link>
        </header>
        <div className="error">{error || 'Link not found'}</div>
      </div>
    );
  }

  return (
    <div className="link-detail-page">
      <header className="page-header">
        <h1>{link.title || link.slug}</h1>
        <div className="page-header-actions">
          <Link to="/links" className="btn-secondary">Back to Links</Link>
        </div>
      </header>

      <div className="link-detail-content">
        <section className="link-detail-main">
          <div className="link-detail-url">
            <label>Target URL</label>
            <a href={link.url} target="_blank" rel="noopener noreferrer">
              {link.url}
            </a>
          </div>

          <div className="link-detail-short">
            <label>Short URL</label>
            <div className="short-url-row">
              <code onClick={copyShortUrl}>/{link.slug}</code>
              <button onClick={copyShortUrl} className="btn-secondary">Copy</button>
            </div>
          </div>

          {link.description && (
            <div className="link-detail-description">
              <label>Description</label>
              <p>{link.description}</p>
            </div>
          )}

        </section>

        <section className="link-detail-meta">
          <h2>Details</h2>
          <dl>
            <dt>Status</dt>
            <dd>
              <span className={`status-badge ${link.is_unread ? 'unread' : 'read'}`}>
                {link.is_unread ? 'Unread' : 'Read'}
              </span>
            </dd>
            <dt>Visibility</dt>
            <dd>{link.is_public ? 'Public' : 'Private'}</dd>
            <dt>Clicks</dt>
            <dd>{link.click_count}</dd>
            <dt>Created</dt>
            <dd>{new Date(link.created_at).toLocaleString()}</dd>
            <dt>Updated</dt>
            <dd>{new Date(link.updated_at).toLocaleString()}</dd>
          </dl>
        </section>

        <section className="link-detail-tags-section">
          <h2>Tags</h2>
          {link.tags && link.tags.length > 0 ? (
            <div className="tags">
              {link.tags.map((tag) => (
                <Link
                  key={tag.id}
                  to={`/links?tag=${encodeURIComponent(tag.name)}`}
                  className="tag"
                >
                  {tag.name}
                </Link>
              ))}
            </div>
          ) : (
            <p className="no-tags">No tags</p>
          )}
        </section>

        <section className="link-detail-actions">
          <h2>Actions</h2>
          <div className="action-buttons">
            <button onClick={toggleUnread}>
              {link.is_unread ? 'Mark as Read' : 'Mark as Unread'}
            </button>
            <Link to={`/links/${link.slug}/edit`} className="btn-primary">
              Edit Link
            </Link>
            <button onClick={handleDelete} className="btn-danger">
              Delete Link
            </button>
          </div>
        </section>
      </div>
    </div>
  );
}
