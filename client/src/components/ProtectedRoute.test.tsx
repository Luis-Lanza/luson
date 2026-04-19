import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { ProtectedRoute } from './ProtectedRoute';

// Mock auth store
let mockAuthState = {
  isAuthenticated: false,
  isAdmin: () => false,
  hasRole: (_role: string) => false,
};

vi.mock('../stores/auth-store', () => ({
  useAuthStore: () => mockAuthState,
  ROLE: {
    ADMIN: 'admin',
    ENCARGADO_ALMACEN: 'encargado_almacen',
    CAJERO: 'cajero',
  },
}));

describe('ProtectedRoute', () => {
  beforeEach(() => {
    mockAuthState = {
      isAuthenticated: false,
      isAdmin: () => false,
      hasRole: () => false,
    };
  });

  it('should redirect to login when not authenticated', () => {
    mockAuthState.isAuthenticated = false;

    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/dashboard" element={<div data-testid="dashboard">Dashboard</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('login-page')).toBeInTheDocument();
    expect(screen.queryByTestId('dashboard')).not.toBeInTheDocument();
  });

  it('should render children when authenticated', () => {
    mockAuthState.isAuthenticated = true;

    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/dashboard" element={<div data-testid="dashboard">Dashboard</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('dashboard')).toBeInTheDocument();
    expect(screen.queryByTestId('login-page')).not.toBeInTheDocument();
  });

  it('should redirect to login when admin required but user is not admin', () => {
    mockAuthState.isAuthenticated = true;
    mockAuthState.isAdmin = () => false;

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route path="/dashboard" element={<div data-testid="dashboard">Dashboard</div>} />
          <Route element={<ProtectedRoute roles={['admin']} />}>
            <Route path="/admin" element={<div data-testid="admin-page">Admin</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('dashboard')).toBeInTheDocument();
    expect(screen.queryByTestId('admin-page')).not.toBeInTheDocument();
  });

  it('should render admin page when admin required and user is admin', () => {
    mockAuthState.isAuthenticated = true;
    mockAuthState.isAdmin = () => true;
    mockAuthState.hasRole = (role: string) => role === 'admin';

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route path="/dashboard" element={<div data-testid="dashboard">Dashboard</div>} />
          <Route element={<ProtectedRoute roles={['admin']} />}>
            <Route path="/admin" element={<div data-testid="admin-page">Admin</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('admin-page')).toBeInTheDocument();
    expect(screen.queryByTestId('dashboard')).not.toBeInTheDocument();
  });

  it('should allow access with specific role', () => {
    mockAuthState.isAuthenticated = true;
    mockAuthState.hasRole = (role: string) => role === 'cajero';

    render(
      <MemoryRouter initialEntries={['/pos']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route element={<ProtectedRoute roles={['cajero', 'admin']} />}>
            <Route path="/pos" element={<div data-testid="pos-page">POS</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('pos-page')).toBeInTheDocument();
  });

  it('should redirect to login when specific role required but user does not have it', () => {
    mockAuthState.isAuthenticated = true;
    mockAuthState.hasRole = () => false;

    render(
      <MemoryRouter initialEntries={['/pos']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route path="/dashboard" element={<div data-testid="dashboard">Dashboard</div>} />
          <Route element={<ProtectedRoute roles={['cajero']} />}>
            <Route path="/pos" element={<div data-testid="pos-page">POS</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('dashboard')).toBeInTheDocument();
    expect(screen.queryByTestId('pos-page')).not.toBeInTheDocument();
  });

  it('should render outlet content for any authenticated user when no roles specified', () => {
    mockAuthState.isAuthenticated = true;

    render(
      <MemoryRouter initialEntries={['/suppliers']}>
        <Routes>
          <Route path="/login" element={<div data-testid="login-page">Login</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/suppliers" element={<div data-testid="suppliers-page">Suppliers</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByTestId('suppliers-page')).toBeInTheDocument();
  });
});
