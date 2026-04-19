import { useState, useEffect } from 'react';
import { purchaseBatchesApi, suppliersApi, productsApi } from '../lib/api';
import { useAuthStore } from '../stores/auth-store';
import type { PurchaseBatch, Supplier, Product } from '../types';

export function PurchaseBatchesPage() {
  const { user } = useAuthStore();
  const [batches, setBatches] = useState<PurchaseBatch[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const isAdmin = user?.role === 'admin';

  useEffect(() => {
    loadBatches();
  }, []);

  const loadBatches = async () => {
    try {
      setLoading(true);
      const data = await purchaseBatchesApi.list();
      setBatches(data);
      setError(null);
    } catch (err) {
      setError('Error al cargar lotes de compra');
    } finally {
      setLoading(false);
    }
  };

  const handleProcess = async (id: string) => {
    try {
      await purchaseBatchesApi.process(id);
      loadBatches();
    } catch (err) {
      setError('Error al procesar lote');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('es-BO', { style: 'currency', currency: 'BOB' }).format(amount);
  };

  return (
    <div className="p-6">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-text">Lotes de Compra</h1>
        {isAdmin && (
          <button
            onClick={() => setShowCreateForm(!showCreateForm)}
            className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark"
          >
            {showCreateForm ? 'Cancelar' : 'Nuevo Lote'}
          </button>
        )}
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-danger/10 p-4 text-danger">
          {error}
        </div>
      )}

      {showCreateForm && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Crear Lote de Compra</h2>
          <PurchaseBatchForm
            onSuccess={() => { setShowCreateForm(false); loadBatches(); }}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      )}

      {loading ? (
        <div className="text-center text-text-muted">Cargando...</div>
      ) : (
        <div className="overflow-hidden rounded-lg bg-white shadow">
          <table className="min-w-full">
            <thead className="bg-surface">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Proveedor</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Fecha</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Total</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Estado</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Items</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {batches.map((batch) => (
                <tr key={batch.id} className="hover:bg-surface/50">
                  <td className="px-6 py-4 text-sm text-text">{batch.supplier_id}</td>
                  <td className="px-6 py-4 text-sm text-text">{formatDate(batch.purchase_date)}</td>
                  <td className="px-6 py-4 text-sm text-text">{formatCurrency(batch.total_cost)}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                      batch.processed
                        ? 'bg-success/10 text-success'
                        : 'bg-warning/10 text-warning'
                    }`}>
                      {batch.processed ? 'Procesado' : 'Pendiente'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-text">
                    {batch.details?.length || 0} items
                  </td>
                  <td className="px-6 py-4 text-sm">
                    {!batch.processed && isAdmin && (
                      <button
                        onClick={() => handleProcess(batch.id)}
                        className="text-primary hover:text-primary-dark"
                      >
                        Procesar
                      </button>
                    )}
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

interface PurchaseBatchFormProps {
  onSuccess: () => void;
  onCancel: () => void;
}

function PurchaseBatchForm({ onSuccess, onCancel }: PurchaseBatchFormProps) {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [selectedSupplier, setSelectedSupplier] = useState('');
  const [purchaseDate, setPurchaseDate] = useState(new Date().toISOString().split('T')[0]);
  const [notes, setNotes] = useState('');
  const [items, setItems] = useState<{ product_id: string; product_type: string; quantity: string; unit_cost: string }[]>([
    { product_id: '', product_type: 'bateria', quantity: '', unit_cost: '' }
  ]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadSuppliers();
    loadProducts();
  }, []);

  const loadSuppliers = async () => {
    try {
      const response = await suppliersApi.list();
      setSuppliers(response.data);
    } catch (err) {
      console.error('Error loading suppliers:', err);
    }
  };

  const loadProducts = async () => {
    try {
      const data = await productsApi.list();
      setProducts(data);
    } catch (err) {
      console.error('Error loading products:', err);
    }
  };

  const addItem = () => {
    setItems([...items, { product_id: '', product_type: 'bateria', quantity: '', unit_cost: '' }]);
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const updateItem = (index: number, field: keyof typeof items[0], value: string) => {
    const newItems = [...items];
    newItems[index] = { ...newItems[index], [field]: value };
    setItems(newItems);
  };

  const calculateTotal = () => {
    return items.reduce((sum, item) => {
      const qty = parseFloat(item.quantity) || 0;
      const cost = parseFloat(item.unit_cost) || 0;
      return sum + (qty * cost);
    }, 0);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (items.some(item => !item.product_id || !item.quantity || !item.unit_cost)) {
      setError('Complete todos los campos de los items');
      return;
    }

    try {
      setLoading(true);
      await purchaseBatchesApi.create({
        supplier_id: selectedSupplier,
        purchase_date: purchaseDate,
        notes: notes || undefined,
        details: items.map(item => ({
          product_id: item.product_id,
          product_type: item.product_type,
          quantity: parseInt(item.quantity, 10),
          unit_cost: parseFloat(item.unit_cost),
        })),
      });
      onSuccess();
    } catch (err) {
      setError('Error al crear lote de compra');
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

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Proveedor</label>
          <select
            value={selectedSupplier}
            onChange={(e) => setSelectedSupplier(e.target.value)}
            required
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          >
            <option value="">Seleccionar proveedor</option>
            {suppliers.map((supplier) => (
              <option key={supplier.id} value={supplier.id}>{supplier.name}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Fecha de Compra</label>
          <input
            type="date"
            value={purchaseDate}
            onChange={(e) => setPurchaseDate(e.target.value)}
            required
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>
      </div>

      <div>
        <label className="mb-1 block text-sm font-medium text-text">Notas</label>
        <textarea
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          rows={2}
          className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
        />
      </div>

      <div className="border-t border-border pt-4">
        <div className="mb-2 flex items-center justify-between">
          <h3 className="font-medium text-text">Productos</h3>
          <button
            type="button"
            onClick={addItem}
            className="text-sm text-primary hover:text-primary-dark"
          >
            + Agregar Producto
          </button>
        </div>

        <div className="space-y-3">
          {items.map((item, index) => (
            <div key={index} className="grid grid-cols-1 gap-2 rounded-lg bg-surface p-3 md:grid-cols-5">
              <div className="md:col-span-2">
                <label className="mb-1 block text-xs font-medium text-text-muted">Producto</label>
                <select
                  value={item.product_id}
                  onChange={(e) => {
                    const product = products.find(p => p.id === e.target.value);
                    updateItem(index, 'product_id', e.target.value);
                    if (product) {
                      updateItem(index, 'product_type', product.product_type);
                    }
                  }}
                  required
                  className="w-full rounded-md border border-border px-2 py-1 text-sm focus:border-primary focus:outline-none"
                >
                  <option value="">Seleccionar producto</option>
                  {products.map((product) => (
                    <option key={product.id} value={product.id}>{product.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="mb-1 block text-xs font-medium text-text-muted">Cantidad</label>
                <input
                  type="number"
                  min="1"
                  value={item.quantity}
                  onChange={(e) => updateItem(index, 'quantity', e.target.value)}
                  required
                  className="w-full rounded-md border border-border px-2 py-1 text-sm focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="mb-1 block text-xs font-medium text-text-muted">Costo Unit.</label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={item.unit_cost}
                  onChange={(e) => updateItem(index, 'unit_cost', e.target.value)}
                  required
                  className="w-full rounded-md border border-border px-2 py-1 text-sm focus:border-primary focus:outline-none"
                />
              </div>
              <div className="flex items-end">
                {items.length > 1 && (
                  <button
                    type="button"
                    onClick={() => removeItem(index)}
                    className="rounded-md bg-danger/10 px-3 py-1 text-sm text-danger hover:bg-danger/20"
                  >
                    Eliminar
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="border-t border-border pt-4">
        <div className="text-right text-lg font-semibold text-text">
          Total: {new Intl.NumberFormat('es-BO', { style: 'currency', currency: 'BOB' }).format(calculateTotal())}
        </div>
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={loading || items.length === 0}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50"
        >
          {loading ? 'Creando...' : 'Crear Lote'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="rounded-md border border-border px-4 py-2 font-medium text-text hover:bg-surface"
        >
          Cancelar
        </button>
      </div>
    </form>
  );
}
