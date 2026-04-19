import { ConnectionStatus } from '../components/ConnectionStatus';

export function POSPage() {
  return (
    <div className="flex h-full flex-col p-6">
      <h1 className="mb-6 text-3xl font-bold text-text">Punto de Venta</h1>
      <div className="flex flex-1 gap-6">
        <div className="flex-1 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold text-text">Productos</h2>
          <p className="text-text-muted">
            Selecciona productos para agregarlos al carrito.
          </p>
        </div>
        <div className="w-96 rounded-lg bg-white p-6 shadow">
          <h2 className="mb-4 text-lg font-semibold text-text">Carrito</h2>
          <p className="text-text-muted">El carrito está vacío.</p>
          <div className="mt-6 border-t border-border pt-4">
            <div className="flex justify-between text-lg font-semibold">
              <span>Total:</span>
              <span>$0.00</span>
            </div>
            <button className="mt-4 w-full rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark disabled:opacity-50">
              Procesar Venta
            </button>
          </div>
        </div>
      </div>
      <ConnectionStatus />
    </div>
  );
}
