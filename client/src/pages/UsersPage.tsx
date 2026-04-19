import { useState, useEffect } from 'react';
import { usersApi } from '../lib/api';
import { useAuthStore } from '../stores/auth-store';
import type { User, UserRole } from '../types';

export function UsersPage() {
  const { isAdmin } = useAuthStore();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const response = await usersApi.list();
      setUsers(response.data);
      setError(null);
    } catch (err) {
      setError('Error al cargar usuarios');
    } finally {
      setLoading(false);
    }
  };

  if (!isAdmin()) {
    return (
      <div className="p-6">
        <div className="rounded-lg bg-danger/10 p-4 text-danger">
          No tienes permisos para ver esta página.
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-text">Usuarios</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark"
        >
          {showCreateForm ? 'Cancelar' : 'Nuevo Usuario'}
        </button>
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-danger/10 p-4 text-danger">
          {error}
        </div>
      )}

      {showCreateForm && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Crear Usuario</h2>
          <CreateUserForm onSuccess={() => { setShowCreateForm(false); loadUsers(); }} />
        </div>
      )}

      {loading ? (
        <div className="text-center text-text-muted">Cargando...</div>
      ) : (
        <div className="overflow-hidden rounded-lg bg-white shadow">
          <table className="min-w-full">
            <thead className="bg-surface">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Usuario</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Rol</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Estado</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Creado</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {users.map((user) => (
                <tr key={user.id} className="hover:bg-surface/50">
                  <td className="px-6 py-4 text-sm text-text">{user.username}</td>
                  <td className="px-6 py-4 text-sm text-text">{user.role}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                      user.active 
                        ? 'bg-success/10 text-success' 
                        : 'bg-danger/10 text-danger'
                    }`}>
                      {user.active ? 'Activo' : 'Inactivo'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-text-muted">
                    {new Date(user.created_at).toLocaleDateString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

function CreateUserForm({ onSuccess }: { onSuccess: () => void }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState<UserRole>('cajero');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      await usersApi.create({ username, password, role });
      onSuccess();
    } catch (err) {
      setError('Error al crear usuario');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="rounded-md bg-danger/10 p-3 text-sm text-danger">
          {error}
        </div>
      )}
      
      <div>
        <label className="mb-1 block text-sm font-medium text-text">Usuario</label>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-text">Contraseña</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-text">Rol</label>
        <select
          value={role}
          onChange={(e) => setRole(e.target.value as UserRole)}
          className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
        >
          <option value="cajero">Cajero</option>
          <option value="encargado_almacen">Encargado de Almacén</option>
          <option value="admin">Administrador</option>
        </select>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50"
      >
        {loading ? 'Creando...' : 'Crear Usuario'}
      </button>
    </form>
  );
}
