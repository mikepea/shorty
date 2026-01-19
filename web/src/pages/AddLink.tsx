import { useState, useEffect, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { links as linksApi, groups as groupsApi } from '../api/client';
import type { Group } from '../api/types';

export default function AddLink() {
  const navigate = useNavigate();
  const [userGroups, setUserGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState('');

  const [url, setUrl] = useState('');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [slug, setSlug] = useState('');
  const [groupId, setGroupId] = useState<number>(0);
  const [isPublic, setIsPublic] = useState(false);
  const [tagsInput, setTagsInput] = useState('');

  useEffect(() => {
    groupsApi.list()
      .then((groups) => {
        setUserGroups(groups);
        if (groups.length > 0) {
          setGroupId(groups[0].id);
        }
      })
      .finally(() => setIsLoading(false));
  }, []);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setIsSaving(true);

    try {
      const link = await linksApi.create(groupId, {
        url,
        title: title || undefined,
        description: description || undefined,
        slug: slug || undefined,
        is_public: isPublic,
      });

      // Set tags if any
      if (tagsInput.trim()) {
        const tagNames = tagsInput.split(',').map((t) => t.trim()).filter(Boolean);
        if (tagNames.length > 0) {
          const { setLinkTags } = await import('../api/client').then(m => m.tags);
          await setLinkTags(link.slug, tagNames);
        }
      }

      navigate('/links');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create link');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="add-link-page">
      <header className="page-header">
        <h1>Add Link</h1>
      </header>

      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit} className="link-form">
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
          <label htmlFor="slug">Custom Slug</label>
          <input
            type="text"
            id="slug"
            value={slug}
            onChange={(e) => setSlug(e.target.value)}
            placeholder="my-custom-slug (optional, auto-generated if empty)"
            pattern="[a-z0-9-]+"
          />
          <small>Lowercase letters, numbers, and hyphens only</small>
        </div>

        <div className="form-group">
          <label htmlFor="group">Group *</label>
          <select
            id="group"
            value={groupId}
            onChange={(e) => setGroupId(Number(e.target.value))}
            required
          >
            {userGroups.map((group) => (
              <option key={group.id} value={group.id}>
                {group.name}
              </option>
            ))}
          </select>
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

        <div className="form-actions">
          <button type="submit" disabled={isSaving}>
            {isSaving ? 'Saving...' : 'Create Link'}
          </button>
          <Link to="/links" className="btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  );
}
