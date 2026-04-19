import { useState, useEffect } from 'react';
import { stockApi, branchesApi } from '../lib/api';
import { useAuthStore } from '../stores/auth-store';
import type { Stock, Branch } from '../types';

export function StockPage() {
  const { user } = useAuthStore();
  const [stock, setStock] = useState<Stock[]>([]);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [lowAlerts, setLowAlerts] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedLocationType, setSelectedLocationType] = useState<string>('');
  const [selectedLocationId, setSelectedLocationId] = useState<string>('');
  const [showLowAlerts, setShowLowAlerts] = useState(false);
  const [showAdjustForm, setShowAdjustForm] = useState(false);

  const isWarehouseManager = user?.role === 'encargado_almacen' || user?.role === 'admin';

  useEffect(() => {
    loadBranches();
    loadLowAlerts();
  }, []);

  useEffect(() => {
    if (selectedLocationType && selectedLocationId) {
      loadStock();
    }
  }, [selectedLocationType, selectedLocationId]);

  const loadBranches = async () => {
    try {
      const response = await branchesApi.list();
      setBranches(response.data);
    } catch (err) {
      console.error('Error loading branches:', err);
    }
  };

  const loadStock = async () => {
    try {
      setLoading(true);
      const params: Record<string, string> = {};
      if (selectedLocationType) params.location_type = selectedLocationType;
      if (selectedLocationId) params.location_id = selectedLocationId;

      const data = await stockApi.list(Object.keys(params).length > 0 ? params : undefined);
      setStock(data);
      setError(null);
    } catch (err) {
      setError('Error al cargar stock');
    } finally {
      setLoading(false);
    }
  };

  const loadLowAlerts = async () => {
    try {
      const data = await stockApi.getLowAlerts();
      setLowAlerts(data);
    } catch (err) {
      console.error('Error loading low alerts:', err);
    }
  };

  const getLocationName = (locationType: string, locationId: string) => {
    if (locationType === 'sucursal') {
      const branch = branches.find(b => b.id === locationId);
      return branch?.name || `Sucursal ${locationId}`;
    }
    if (locationType === 'almacen_general') return 'Almacén General';
    if (locationType === 'camion') return `Camión ${locationId}`;
    return `${locationType} ${locationId}`;
  };

  return (
    <div className="p-6">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-text">Stock</h1>
        {isWarehouseManager && (
          <button
            onClick={() => setShowAdjustForm(!showAdjustForm)}
            className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark"
          >
            {showAdjustForm ? 'Cancelar' : 'Ajustar Stock'}
          </button>
        )}
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-danger/10 p-4 text-danger">
          {error}
        </div>
      )}

      {/* Low Stock Alerts */}
      <div className="mb-6">
        <button
          onClick={() => setShowLowAlerts(!showLowAlerts)}
          className="mb-2 flex items-center gap-2 text-lg font-semibold text-text hover:text-primary"
        >
          <span>Alertas de Stock Bajo</span>
          {lowAlerts.length > 0 && (
            <span className="rounded-full bg-danger px-2 py-0.5 text-xs text-white">
              {lowAlerts.length}
            </span>
          )}
          <span className="transform text-sm transition-transform" style={{ transform: showLowAlerts ? 'rotate(180deg)' : 'rotate(0deg)' }}>▼</span>
        </button>

        {showLowAlerts && (
          <div className="overflow-hidden rounded-lg border border-warning/30 bg-warning/5 shadow">
            {lowAlerts.length === 0 ? (
              <div className="p-4 text-text-muted">No hay alertas de stock bajo</div>
            ) : (
              <table className="min-w-full">
                <thead className="bg-warning/10">
                  <tr>
                    <th className="px-6 py-3 text-left text-sm font-medium text-text">Producto</th>
                    <th className="px-6 py-3 text-left text-sm font-medium text-text">Ubicación</th>
                    <th className="px-6 py-3 text-left text-sm font-medium text-text">Cantidad</th>
                    <th className="px-6 py-3 text-left text-sm font-medium text-text">Mínimo</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-warning/20">
                  {lowAlerts.map((item) => (
                    <tr key={item.id}>
                      <td className="px-6 py-3 text-sm text-text">{item.product_id}</td>
                      <td className="px-6 py-3 text-sm text-text">
                        {getLocationName(item.location_type, item.location_id)}
                      </td>
                      <td className="px-6 py-3 text-sm font-medium text-danger">{item.quantity}</td>
                      <td className="px-6 py-3 text-sm text-text-muted">{item.min_stock_alert}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}
      </div>

      {/* Stock Adjustment Form */}
      {showAdjustForm && isWarehouseManager && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Ajustar Stock</h2>
          <StockAdjustForm
            onSuccess={() => { setShowAdjustForm(false); loadStock(); loadLowAlerts(); }}
            onCancel={() => setShowAdjustForm(false)}
            branches={branches}
          />
        </div>
      )}

      {/* Location Filter */}
      <div className="mb-6 rounded-lg bg-white p-4 shadow">
        <h2 className="mb-4 text-lg font-semibold">Filtrar por Ubicación</h2>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label className="mb-1 block text-sm font-medium text-text">Tipo de Ubicación</label>
            <select
              value={selectedLocationType}
              onChange={(e) => {
                setSelectedLocationType(e.target.value);
                setSelectedLocationId('');
              }}
              className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
            >
              <option value="">Todas</option>
              <option value="almacen_general">Almacén General</option>
              <option value="sucursal">Sucursal</option>
              <option value="camion">Camión</option>
            </select>
          </div>
          {selectedLocationType === 'sucursal' && (
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Sucursal</label>
              <select
                value={selectedLocationId}
                onChange={(e) => setSelectedLocationId(e.target.value)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              >
                <option value="">Seleccionar sucursal</option>
                {branches.map((branch) => (
                  <option key={branch.id} value={branch.id}>{branch.name}</option>
                ))}
              </select>
            </div>
          )}
        </div>
      </div>

      {/* Stock List */}
      {loading ? (
        <div className="text-center text-text-muted">Cargando...</div>
      ) : (
        <div className="overflow-hidden rounded-lg bg-white shadow">
          <table className="min-w-full">
            <thead className="bg-surface">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Producto ID</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Tipo</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Ubicación</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Cantidad</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Mínimo Alerta</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Actualizado</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {stock.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-4 text-center text-text-muted">
                    No hay registros de stock para los filtros seleccionados
                  </td>
                </tr>
              ) : (
                stock.map((item) => (
                  <tr key={item.id} className="hover:bg-surface/50">
                    <td className="px-6 py-4 text-sm text-text">{item.product_id}</td>
                    <td className="px-6 py-4 text-sm text-text">{item.product_type}</td>
                    <td className="px-6 py-4 text-sm text-text">
                      {getLocationName(item.location_type, item.location_id)}
                    </td>
                    <td className={`px-6 py-4 text-sm font-medium ${
                      item.quantity <= item.min_stock_alert ? 'text-danger' : 'text-text'
                    }`}>
                      {item.quantity}
                    </td>
                    <td className="px-6 py-4 text-sm text-text-muted">{item.min_stock_alert}</td>
                    <td className="px-6 py-4 text-sm text-text-muted">
                      {new Date(item.updated_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

interface StockAdjustFormProps {
  onSuccess: () => void;
  onCancel: () => void;
  branches: Branch[];
}

function StockAdjustForm({ onSuccess, onCancel, branches }: StockAdjustFormProps) {
  const [productId, setProductId] = useState('');
  const [productType, setProductType] = useState('bateria');
  const [locationType, setLocationType] = useState('almacen_general');
  const [locationId, setLocationId] = useState('');
  const [quantity, setQuantity] = useState('');
  const [reason, setReason] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setLoading(true);
      await stockApi.adjust({
        product_id: productId,
        product_type: productType,
        location_type: locationType,
        location_id: locationId,
        quantity: parseInt(quantity, 10),
        reason,
      });
      onSuccess();
    } catch (err) {
      setError('Error al ajustar stock');
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
          <label className="mb-1 block text-sm font-medium text-text">ID del Producto</label>
          <input
            type="text"
            value={productId}
            onChange={(e) => setProductId(e.target.value)}
            required
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Tipo de Producto</label>
          <select
            value={productType}
            onChange={(e) => setProductType(e.target.value)}
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          >
            <option value="bateria">Batería</option>
            <option value="accesorio">Accesorio</option>
          </select>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Tipo de Ubicación</label>
          <select
            value={locationType}
            onChange={(e) => setLocationType(e.target.value)}
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          >
            <option value="almacen_general">Almacén General</option>
            <option value="sucursal">Sucursal</option>
            <option value="camion">Camión</option>
          </select>
        </div>
        {locationType === 'sucursal' && (
          <div>
            <label className="mb-1 block text-sm font-medium text-text">Sucursal</label>
            <select
              value={locationId}
              onChange={(e) => setLocationId(e.target.value)}
              required
              className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
            >
              <option value="">Seleccionar sucursal</option>
              {branches.map((branch) => (
                <option key={branch.id} value={branch.id}>{branch.name}</option>
              ))}
            </select>
          </div>
        )}
        {locationType === 'camion' && (
          <div>
            <label className="mb-1 block text-sm font-medium text-text">ID del Camión</label>
            <input
              type="text"
              value={locationId}
              onChange={(e) => setLocationId(e.target.value)}
              required
              className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
            />
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <label className="mb-1 block text-sm font-medium text-text">
            Cantidad (positiva para ingreso, negativa para salida)
          </label>
          <input
            type="number"
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            required
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Motivo</label>
          <input
            type="text"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            required
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={loading}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50"
        >
          {loading ? 'Ajustando...' : 'Ajustar Stock'}
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
