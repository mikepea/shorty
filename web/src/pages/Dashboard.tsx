import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { links, groups, tags } from '../api/client';
import type { Link as LinkType, Group, Tag } from '../api/types';

export default function Dashboard() {
  const { user, logout } = useAuth();
  const [recentLinks, setRecentLinks] = useState<LinkType[]>([]);
  const [userGroups, setUserGroups] = useState<Group[]>([]);
  const [topTags, setTopTags] = useState<Tag[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      links.search(),
      groups.list(),
      tags.list(),
    ])
      .then(([linksData, groupsData, tagsData]) => {
        setRecentLinks(linksData.slice(0, 5));
        setUserGroups(groupsData);
        setTopTags(tagsData.slice(0, 10));
      })
      .finally(() => setIsLoading(false));
  }, []);

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Shorty</h1>
        <div className="user-menu">
          <span>Welcome, {user?.name}</span>
          <button onClick={logout}>Logout</button>
        </div>
      </header>

      <nav className="dashboard-nav">
        <Link to="/links">All Links</Link>
        <Link to="/links/new">Add Link</Link>
        <Link to="/groups">Groups</Link>
        <Link to="/settings">Settings</Link>
      </nav>

      <main className="dashboard-content">
        <section className="recent-links">
          <h2>Recent Links</h2>
          {recentLinks.length === 0 ? (
            <p>No links yet. <Link to="/links/new">Add your first link</Link></p>
          ) : (
            <ul className="link-list">
              {recentLinks.map((link) => (
                <li key={link.id} className="link-item">
                  <a href={`/${link.slug}`} target="_blank" rel="noopener noreferrer">
                    {link.title || link.url}
                  </a>
                  <span className="link-slug">/{link.slug}</span>
                  <span className="link-clicks">{link.click_count} clicks</span>
                </li>
              ))}
            </ul>
          )}
          <Link to="/links" className="view-all">View all links</Link>
        </section>

        <section className="groups-summary">
          <h2>Your Groups</h2>
          <ul className="group-list">
            {userGroups.map((group) => (
              <li key={group.id}>
                <Link to={`/groups/${group.id}`}>{group.name}</Link>
              </li>
            ))}
          </ul>
        </section>

        <section className="tags-cloud">
          <h2>Popular Tags</h2>
          <div className="tags">
            {topTags.map((tag) => (
              <Link
                key={tag.id}
                to={`/links?tag=${encodeURIComponent(tag.name)}`}
                className="tag"
              >
                {tag.name} ({tag.link_count})
              </Link>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
