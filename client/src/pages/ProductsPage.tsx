import { useState, useEffect } from 'react';
import { productsApi } from '../lib/api';
import type { Product, ProductType, BatteryType, Polarity, VehicleType } from '../types';

export function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      setLoading(true);
      const data = await productsApi.list();
      setProducts(data);
      setError(null);
    } catch (err) {
      setError('Error al cargar productos');
    } finally {
      setLoading(false);
    }
  };

  const getProductTypeLabel = (type: ProductType) => {
    return type === 'bateria' ? 'Batería' : 'Accesorio';
  };

  return (
    <div className="p-6">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-text">Productos</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark"
        >
          {showCreateForm ? 'Cancelar' : 'Nuevo Producto'}
        </button>
      </div>

      {error && (
        <div className="mb-4 rounded-lg bg-danger/10 p-4 text-danger">
          {error}
        </div>
      )}

      {showCreateForm && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Crear Producto</h2>
          <ProductForm
            onSuccess={() => { setShowCreateForm(false); loadProducts(); }}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      )}

      {editingProduct && (
        <div className="mb-6 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold">Editar Producto</h2>
          <ProductForm
            product={editingProduct}
            onSuccess={() => { setEditingProduct(null); loadProducts(); }}
            onCancel={() => setEditingProduct(null)}
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
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Nombre</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Tipo</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Marca/Modelo</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Precio Mínimo</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Estado</th>
                <th className="px-6 py-3 text-left text-sm font-medium text-text-muted">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {products.map((product) => (
                <tr key={product.id} className="hover:bg-surface/50">
                  <td className="px-6 py-4 text-sm text-text">{product.name}</td>
                  <td className="px-6 py-4 text-sm text-text">{getProductTypeLabel(product.product_type)}</td>
                  <td className="px-6 py-4 text-sm text-text">
                    {product.brand && product.model ? `${product.brand} - ${product.model}` : '-'}
                  </td>
                  <td className="px-6 py-4 text-sm text-text">${product.min_sale_price.toFixed(2)}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                      product.active
                        ? 'bg-success/10 text-success'
                        : 'bg-danger/10 text-danger'
                    }`}>
                      {product.active ? 'Activo' : 'Inactivo'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <button
                      onClick={() => setEditingProduct(product)}
                      className="text-primary hover:text-primary-dark"
                    >
                      Editar
                    </button>
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

interface ProductFormProps {
  product?: Product;
  onSuccess: () => void;
  onCancel: () => void;
}

function ProductForm({ product, onSuccess, onCancel }: ProductFormProps) {
  const [productType, setProductType] = useState<ProductType>(product?.product_type || 'bateria');
  const [name, setName] = useState(product?.name || '');
  const [description, setDescription] = useState(product?.description || '');
  const [brand, setBrand] = useState(product?.brand || '');
  const [model, setModel] = useState(product?.model || '');
  const [voltage, setVoltage] = useState(product?.voltage?.toString() || '');
  const [amperage, setAmperage] = useState(product?.amperage?.toString() || '');
  const [batteryType, setBatteryType] = useState<BatteryType>(product?.battery_type || 'seca');
  const [polarity, setPolarity] = useState<Polarity>(product?.polarity || 'izquierda');
  const [acidLiters, setAcidLiters] = useState(product?.acid_liters?.toString() || '');
  const [vehicleType, setVehicleType] = useState<VehicleType>(product?.vehicle_type || 'auto');
  const [minSalePrice, setMinSalePrice] = useState(product?.min_sale_price.toString() || '');
  const [active, setActive] = useState(product?.active ?? true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setLoading(true);
      const baseData = {
        name,
        description: description || undefined,
        product_type: productType,
        min_sale_price: parseFloat(minSalePrice),
        active,
      };

      let data;
      if (productType === 'bateria') {
        data = {
          ...baseData,
          brand: brand || undefined,
          model: model || undefined,
          voltage: voltage ? parseFloat(voltage) : undefined,
          amperage: amperage ? parseFloat(amperage) : undefined,
          battery_type: batteryType,
          polarity: polarity,
          acid_liters: acidLiters ? parseFloat(acidLiters) : undefined,
          vehicle_type: vehicleType,
        };
      } else {
        data = baseData;
      }

      if (product) {
        await productsApi.update(product.id, data);
      } else {
        await productsApi.create(data as Omit<Product, 'id' | 'created_at'>);
      }
      onSuccess();
    } catch (err) {
      setError(product ? 'Error al actualizar producto' : 'Error al crear producto');
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
          <label className="mb-1 block text-sm font-medium text-text">Tipo de Producto</label>
          <select
            value={productType}
            onChange={(e) => setProductType(e.target.value as ProductType)}
            disabled={!!product}
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none disabled:bg-surface"
          >
            <option value="bateria">Batería</option>
            <option value="accesorio">Accesorio</option>
          </select>
        </div>

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
      </div>

      {productType === 'bateria' ? (
        <>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Marca</label>
              <input
                type="text"
                value={brand}
                onChange={(e) => setBrand(e.target.value)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Modelo</label>
              <input
                type="text"
                value={model}
                onChange={(e) => setModel(e.target.value)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Voltaje (V)</label>
              <input
                type="number"
                step="0.1"
                value={voltage}
                onChange={(e) => setVoltage(e.target.value)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Amperaje (Ah)</label>
              <input
                type="number"
                step="0.1"
                value={amperage}
                onChange={(e) => setAmperage(e.target.value)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Tipo de Batería</label>
              <select
                value={batteryType}
                onChange={(e) => setBatteryType(e.target.value as BatteryType)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              >
                <option value="seca">Seca</option>
                <option value="liquida">Líquida</option>
              </select>
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Polaridad</label>
              <select
                value={polarity}
                onChange={(e) => setPolarity(e.target.value as Polarity)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              >
                <option value="izquierda">Izquierda</option>
                <option value="derecha">Derecha</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Litros de Ácido</label>
              <input
                type="number"
                step="0.1"
                value={acidLiters}
                onChange={(e) => setAcidLiters(e.target.value)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-text">Tipo de Vehículo</label>
              <select
                value={vehicleType}
                onChange={(e) => setVehicleType(e.target.value as VehicleType)}
                className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
              >
                <option value="auto">Auto</option>
                <option value="moto">Moto</option>
                <option value="otro">Otro</option>
              </select>
            </div>
          </div>
        </>
      ) : (
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Descripción</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>
      )}

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <label className="mb-1 block text-sm font-medium text-text">Precio Mínimo de Venta</label>
          <input
            type="number"
            step="0.01"
            min="0"
            value={minSalePrice}
            onChange={(e) => setMinSalePrice(e.target.value)}
            required
            className="w-full rounded-md border border-border px-3 py-2 focus:border-primary focus:outline-none"
          />
        </div>
        <div className="flex items-center">
          <label className="flex items-center gap-2 text-sm font-medium text-text">
            <input
              type="checkbox"
              checked={active}
              onChange={(e) => setActive(e.target.checked)}
              className="rounded border-border"
            />
            Activo
          </label>
        </div>
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={loading}
          className="rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50"
        >
          {loading ? 'Guardando...' : (product ? 'Actualizar Producto' : 'Crear Producto')}
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
