import { useState, useEffect, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { links as linksApi, tags as tagsApi } from '../api/client';
import type { Link as LinkType } from '../api/types';

export default function EditLink() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState('');

  const [url, setUrl] = useState('');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isPublic, setIsPublic] = useState(false);
  const [isUnread, setIsUnread] = useState(false);
  const [tagsInput, setTagsInput] = useState('');
  const [originalLink, setOriginalLink] = useState<LinkType | null>(null);

  useEffect(() => {
    if (!slug) return;

    linksApi.get(slug)
      .then((link) => {
        setOriginalLink(link);
        setUrl(link.url);
        setTitle(link.title || '');
        setDescription(link.description || '');
        setIsPublic(link.is_public);
        setIsUnread(link.is_unread);
        setTagsInput(link.tags?.map(t => t.name).join(', ') || '');
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load link'))
      .finally(() => setIsLoading(false));
  }, [slug]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!slug) return;

    setError('');
    setIsSaving(true);

    try {
      await linksApi.update(slug, {
        url,
        title: title || undefined,
        description: description || undefined,
        is_public: isPublic,
        is_unread: isUnread,
      });

      // Update tags
      const tagNames = tagsInput.split(',').map((t) => t.trim()).filter(Boolean);
      await tagsApi.setLinkTags(slug, tagNames);

      navigate(`/links/${slug}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update link');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  if (!originalLink) {
    return (
      <div className="add-link-page">
        <header className="page-header">
          <h1>Link Not Found</h1>
        </header>
        <div className="error">{error || 'Link not found'}</div>
        <Link to="/links" className="btn-secondary">Back to Links</Link>
      </div>
    );
  }

  return (
    <div className="add-link-page">
      <header className="page-header">
        <h1>Edit Link</h1>
      </header>

      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit} className="link-form">
        <div className="form-group">
          <label htmlFor="slug">Slug</label>
          <input
            type="text"
            id="slug"
            value={`/${originalLink.slug}`}
            disabled
          />
          <small>Slug cannot be changed</small>
        </div>

        <div className="form-group">
          <label htmlFor="url">URL *</label>
          <input
            type="url"
            id="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
            placeholder="https://example.com"
          />
        </div>

        <div className="form-group">
          <label htmlFor="title">Title</label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Page title (optional)"
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Notes about this link (optional)"
            rows={3}
          />
        </div>

        <div className="form-group">
          <label htmlFor="tags">Tags</label>
          <input
            type="text"
            id="tags"
            value={tagsInput}
            onChange={(e) => setTagsInput(e.target.value)}
            placeholder="tag1, tag2, tag3 (comma separated)"
          />
        </div>

        <div className="form-group checkbox">
          <label>
            <input
              type="checkbox"
              checked={isPublic}
              onChange={(e) => setIsPublic(e.target.checked)}
            />
            Make this link public
          </label>
        </div>

        <div className="form-group checkbox">
          <label>
            <input
              type="checkbox"
              checked={isUnread}
              onChange={(e) => setIsUnread(e.target.checked)}
            />
            Mark as unread
          </label>
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSaving}>
            {isSaving ? 'Saving...' : 'Save Changes'}
          </button>
          <Link to={`/links/${slug}`} className="btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  );
}
