import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuthStore, ROLE } from '../stores/auth-store';
import type { UserRole } from '../types';

interface ProtectedRouteProps {
  roles?: UserRole[];
}

export function ProtectedRoute({ roles }: ProtectedRouteProps) {
  const location = useLocation();
  const { isAuthenticated, isAdmin, hasRole } = useAuthStore();

  // Not authenticated - redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Check role requirements
  if (roles && roles.length > 0) {
    const hasRequiredRole = roles.some((role) => {
      if (role === ROLE.ADMIN) {
        return isAdmin();
      }
      return hasRole(role);
    });

    if (!hasRequiredRole) {
      // User is authenticated but doesn't have required role
      // Redirect to dashboard (or could show 403 page)
      return <Navigate to="/dashboard" replace />;
    }
  }

  // All checks passed - render the protected content
  return <Outlet />;
}
