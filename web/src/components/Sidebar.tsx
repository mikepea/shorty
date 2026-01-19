import { NavLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Sidebar() {
  const { user, logout } = useAuth();

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <h1>Shorty</h1>
      </div>

      <nav className="sidebar-nav">
        <NavLink to="/" end>Dashboard</NavLink>
        <NavLink to="/links">Links</NavLink>
        <NavLink to="/links/new">Add Link</NavLink>
        <NavLink to="/groups">Groups</NavLink>
        <NavLink to="/settings">Settings</NavLink>
        {user?.system_role === 'admin' && <NavLink to="/admin">Admin</NavLink>}
      </nav>

      <div className="sidebar-user">
        <span className="sidebar-user-name">{user?.name}</span>
        <span className="sidebar-user-email">{user?.email}</span>
        <button onClick={logout} className="sidebar-logout">Logout</button>
      </div>
    </aside>
  );
}
