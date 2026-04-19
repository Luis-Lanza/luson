import { useState, useEffect } from 'react';
import { branchesApi } from '../lib/api';
import { useAuthStore } from '../stores/auth-store';
import type { Branch } from '../types';

export function BranchesPage() {
  const { isAdmin } = useAuthStore();
  const [branches, setBranches] = useState<Branch[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    loadBranches();
  }, []);

  const loadBranches = async () => {
    try {
      setLoading(true);
      const response = await branchesApi.list();
      setBranches(response.data);
      setError(null);
    } catch (err) {
      setError('Error al cargar sucursales');
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
        <h1 className="text-3xl font-bold text-text">Sucursales</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark"
        >
          {showCreateForm ? 'Cancelar' : 'Nueva Sucursal'}
        </button>
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-danger/10 p-4 text-danger">
          {error}
        </div>
      )}

      {showCreateForm && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Crear Sucursal</h2>
          <CreateBranchForm onSuccess={() => { setShowCreateForm(false); loadBranches(); }} />
        </div>
      )}

      {loading ? (
        <div className="text-center text-text-muted">Cargando...</div>
      ) : (
        <div className="overflow-hidden rounded-lg bg-white shadow">
          <table className="min-w-full">
            <thead className="bg-surface">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Nombre</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Dirección</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Caja Chica</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Estado</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Creado</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {branches.map((branch) => (
                <tr key={branch.id} className="hover:bg-surface/50">
                  <td className="px-6 py-4 text-sm text-text">{branch.name}</td>
                  <td className="px-6 py-4 text-sm text-text">{branch.address || '-'}</td>
                  <td className="px-6 py-4 text-sm text-text">
                    ${branch.petty_cash_balance.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                      branch.active 
                        ? 'bg-success/10 text-success' 
                        : 'bg-danger/10 text-danger'
                    }`}>
                      {branch.active ? 'Activa' : 'Inactiva'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-text-muted">
                    {new Date(branch.created_at).toLocaleDateString()}
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

function CreateBranchForm({ onSuccess }: { onSuccess: () => void }) {
  const [name, setName] = useState('');
  const [address, setAddress] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      await branchesApi.create({ name, address: address || undefined });
      onSuccess();
    } catch (err) {
      setError('Error al crear sucursal');
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
        <label className="mb-1 block text-sm font-medium text-text">Nombre</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
        />
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-text">Dirección</label>
        <input
          type="text"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
          className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50"
      >
        {loading ? 'Creando...' : 'Crear Sucursal'}
      </button>
    </form>
  );
}
