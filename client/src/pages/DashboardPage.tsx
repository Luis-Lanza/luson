import { ConnectionStatus } from '../components/ConnectionStatus';

export function DashboardPage() {
  return (
    <div className="p-6">
      <h1 className="mb-6 text-3xl font-bold text-text">Dashboard</h1>
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-lg bg-white p-6 shadow">
          <h3 className="text-sm font-medium text-text-muted">Ventas Hoy</h3>
          <p className="mt-2 text-2xl font-bold text-text">$0.00</p>
        </div>
        <div className="rounded-lg bg-white p-6 shadow">
          <h3 className="text-sm font-medium text-text-muted">Productos</h3>
          <p className="mt-2 text-2xl font-bold text-text">0</p>
        </div>
        <div className="rounded-lg bg-white p-6 shadow">
          <h3 className="text-sm font-medium text-text-muted">Clientes</h3>
          <p className="mt-2 text-2xl font-bold text-text">0</p>
        </div>
        <div className="rounded-lg bg-white p-6 shadow">
          <h3 className="text-sm font-medium text-text-muted">Stock Bajo</h3>
          <p className="mt-2 text-2xl font-bold text-danger">0</p>
        </div>
      </div>
      <ConnectionStatus />
    </div>
  );
}
