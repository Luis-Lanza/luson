import { createBrowserRouter, Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from './stores/auth-store';
import { ProtectedRoute } from './components/ProtectedRoute';
import { App } from './App';
import { LoginPage } from './pages/LoginPage';
import { DashboardPage } from './pages/DashboardPage';
import { POSPage } from './pages/POSPage';
import { UsersPage } from './pages/UsersPage';
import { BranchesPage } from './pages/BranchesPage';
import { SuppliersPage } from './pages/SuppliersPage';

// Public only route (redirects to dashboard if already authenticated)
function PublicRoute() {
  const { isAuthenticated, getRedirectPath } = useAuthStore();
  
  if (isAuthenticated) {
    return <Navigate to={getRedirectPath()} replace />;
  }
  
  return <Outlet />;
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        element: <PublicRoute />,
        children: [
          {
            path: 'login',
            element: <LoginPage />,
          },
        ],
      },
      {
        element: <ProtectedRoute />,
        children: [
          {
            index: true,
            element: <Navigate to="/dashboard" replace />,
          },
          {
            path: 'dashboard',
            element: <DashboardPage />,
          },
          {
            path: 'pos',
            element: <POSPage />,
          },
          {
            path: 'suppliers',
            element: <SuppliersPage />,
          },
        ],
      },
      {
        element: <ProtectedRoute roles={['admin']} />,
        children: [
          {
            path: 'users',
            element: <UsersPage />,
          },
          {
            path: 'branches',
            element: <BranchesPage />,
          },
        ],
      },
    ],
  },
]);
