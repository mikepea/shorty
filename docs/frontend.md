# Frontend Development

This guide covers developing the React frontend for Shorty.

## Project Structure

```
web/
├── src/
│   ├── api/               # API client and types
│   │   ├── client.ts      # Fetch wrapper and API methods
│   │   ├── client.test.ts # API client tests
│   │   └── types.ts       # TypeScript interfaces
│   ├── components/        # Reusable UI components
│   ├── context/           # React context (auth state)
│   │   ├── AuthContext.tsx
│   │   └── AuthContext.test.tsx
│   ├── pages/             # Page components
│   └── test/              # Test setup
│       └── setup.ts       # Vitest configuration
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## Getting Started

```bash
cd web
npm install
npm run dev
```

The dev server runs on `http://localhost:3000` with hot reloading. API requests are proxied to the backend at `http://localhost:8080`.

## Available Scripts

| Script | Description |
|--------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run preview` | Preview production build |
| `npm run lint` | Run ESLint |
| `npm test` | Run tests once |
| `npm run test:watch` | Run tests in watch mode |
| `npm run test:coverage` | Run tests with coverage |

## Running Tests

The frontend uses [Vitest](https://vitest.dev/) with React Testing Library.

### Run All Tests

```bash
npm test
```

### Watch Mode (Development)

```bash
npm run test:watch
```

### With Coverage

```bash
npm run test:coverage
```

### Test Structure

Tests are co-located with source files:

- `src/api/client.test.ts` - API client tests
- `src/context/AuthContext.test.tsx` - Auth context tests

### Writing Tests

```tsx
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent />);
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });

  it('handles click', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();

    render(<MyComponent onClick={onClick} />);
    await user.click(screen.getByRole('button'));

    expect(onClick).toHaveBeenCalled();
  });
});
```

## Code Style

- Use TypeScript for all new code
- Follow React best practices
- Use functional components with hooks
- Prefer named exports

### Linting

```bash
npm run lint
```

## API Client

The API client in `src/api/client.ts` provides typed methods for all backend endpoints:

```typescript
import { auth, links, groups } from './api/client';

// Authentication
await auth.login(email, password);
await auth.register(email, password, name);
await auth.logout();
await auth.me();
await auth.changePassword(currentPassword, newPassword);

// Links
await links.search({ q: 'query', tag: 'javascript' });
await links.get('my-slug');
await links.create(groupId, { url, title });
await links.update('my-slug', { title: 'New Title' });
await links.delete('my-slug');

// Groups
await groups.list();
await groups.create('Group Name');
await groups.addMember(groupId, email, 'member');
```

### Error Handling

The client throws `APIError` for non-2xx responses:

```typescript
import { APIError } from './api/client';

try {
  await auth.login(email, password);
} catch (err) {
  if (err instanceof APIError) {
    console.log(err.status);  // HTTP status code
    console.log(err.message); // Error message from server
  }
}
```

## Auth Context

The `AuthContext` manages authentication state:

```tsx
import { useAuth } from './context/AuthContext';

function MyComponent() {
  const { user, isLoading, login, logout } = useAuth();

  if (isLoading) return <div>Loading...</div>;
  if (!user) return <div>Please log in</div>;

  return (
    <div>
      <p>Welcome, {user.name}!</p>
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

### User Object

```typescript
interface User {
  id: number;
  email: string;
  name: string;
  system_role: 'admin' | 'user';
  has_password: boolean;  // false for SSO-only users
  created_at: string;
}
```

## Adding a New Page

1. Create the page component in `src/pages/`:

```tsx
// src/pages/MyPage.tsx
export default function MyPage() {
  return (
    <div className="my-page">
      <h1>My Page</h1>
    </div>
  );
}
```

2. Add the route in `src/App.tsx`:

```tsx
import MyPage from './pages/MyPage';

// In the Routes:
<Route path="/my-page" element={<MyPage />} />
```

3. Add navigation link in `src/components/Sidebar.tsx` if needed.

## CI Checks

The CI pipeline runs these frontend checks:

1. **Tests** - `npm test`
2. **Build** - `npm run build`

Both must pass before merging.
