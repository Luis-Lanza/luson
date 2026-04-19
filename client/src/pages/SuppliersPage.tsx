import { useState, useEffect } from 'react';
import { suppliersApi } from '../lib/api';
import type { Supplier } from '../types';

export function SuppliersPage() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    loadSuppliers();
  }, []);

  const loadSuppliers = async () => {
    try {
      setLoading(true);
      const response = await suppliersApi.list();
      setSuppliers(response.data);
      setError(null);
    } catch (err) {
      setError('Error al cargar proveedores');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-text">Proveedores</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark"
        >
          {showCreateForm ? 'Cancelar' : 'Nuevo Proveedor'}
        </button>
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-danger/10 p-4 text-danger">
          {error}
        </div>
      )}

      {showCreateForm && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Crear Proveedor</h2>
          <CreateSupplierForm onSuccess={() => { setShowCreateForm(false); loadSuppliers(); }} />
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
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Contacto</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Dirección</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Estado</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Creado</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {suppliers.map((supplier) => (
                <tr key={supplier.id} className="hover:bg-surface/50">
                  <td className="px-6 py-4 text-sm text-text">{supplier.name}</td>
                  <td className="px-6 py-4 text-sm text-text">{supplier.contact || '-'}</td>
                  <td className="px-6 py-4 text-sm text-text">{supplier.address || '-'}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                      supplier.active 
                        ? 'bg-success/10 text-success' 
                        : 'bg-danger/10 text-danger'
                    }`}>
                      {supplier.active ? 'Activo' : 'Inactivo'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-text-muted">
                    {new Date(supplier.created_at).toLocaleDateString()}
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

function CreateSupplierForm({ onSuccess }: { onSuccess: () => void }) {
  const [name, setName] = useState('');
  const [contact, setContact] = useState('');
  const [address, setAddress] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setLoading(true);
      await suppliersApi.create({ 
        name, 
        contact: contact || undefined, 
        address: address || undefined 
      });
      onSuccess();
    } catch (err) {
      setError('Error al crear proveedor');
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
        <label className="mb-1 block text-sm font-medium text-text">Contacto</label>
        <input
          type="text"
          value={contact}
          onChange={(e) => setContact(e.target.value)}
          placeholder="Teléfono o email"
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
        {loading ? 'Creando...' : 'Crear Proveedor'}
      </button>
    </form>
  );
}
