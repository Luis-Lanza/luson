import { useState, useEffect, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../stores/auth-store';

export function LoginPage() {
  const navigate = useNavigate();
  const { login, isLoading, error, clearError, getRedirectPath } = useAuthStore();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  // Clear error when inputs change
  useEffect(() => {
    if (error) {
      clearError();
    }
  }, [username, password, error, clearError]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    if (!username.trim() || !password.trim()) {
      return;
    }

    try {
      await login(username.trim(), password);
      const redirectPath = getRedirectPath();
      navigate(redirectPath);
    } catch {
      // Error is handled by the store
    }
  };

  const isFormValid = username.trim().length > 0 && password.trim().length > 0;

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
        <h1 className="mb-6 text-center text-2xl font-bold text-text">
          Battery POS
        </h1>
        <p className="mb-4 text-center text-text-muted">
          Inicia sesión para continuar
        </p>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <input
              type="text"
              placeholder="Usuario"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              disabled={isLoading}
              className="w-full rounded-md border border-border px-4 py-2 focus:border-primary focus:outline-none disabled:opacity-50"
              autoComplete="username"
            />
          </div>
          
          <div>
            <input
              type="password"
              placeholder="Contraseña"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={isLoading}
              className="w-full rounded-md border border-border px-4 py-2 focus:border-primary focus:outline-none disabled:opacity-50"
              autoComplete="current-password"
            />
          </div>
          
          {error && (
            <div 
              role="alert"
              className="rounded-md bg-danger/10 p-3 text-sm text-danger"
            >
              {error}
            </div>
          )}
          
          <button
            type="submit"
            disabled={!isFormValid || isLoading}
            className="w-full rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? 'Iniciando sesión...' : 'Iniciar sesión'}
          </button>
        </form>
      </div>
    </div>
  );
}
