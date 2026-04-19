import { useState } from 'react';
import type { Transfer, TransferStatus } from '../types';

// Mock data for now - will connect to API
const mockTransfers: Transfer[] = [];

const statusColors: Record<TransferStatus, string> = {
  pendiente: 'bg-yellow-100 text-yellow-800',
  aprobada: 'bg-blue-100 text-blue-800',
  rechazada: 'bg-red-100 text-red-800',
  enviada: 'bg-purple-100 text-purple-800',
  recibida: 'bg-green-100 text-green-800',
};

export function TransfersPage() {
  const [transfers] = useState<Transfer[]>(mockTransfers);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-slate-900 dark:text-white">
          Transferencias
        </h1>
        <button className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90">
          Nueva Transferencia
        </button>
      </div>

      {/* Status filters */}
      <div className="flex gap-2">
        {(['pendiente', 'aprobada', 'enviada', 'recibida', 'rechazada'] as TransferStatus[]).map((status) => (
          <button
            key={status}
            className={`rounded-full px-3 py-1 text-xs font-medium ${statusColors[status]}`}
          >
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </button>
        ))}
      </div>

      {/* Transfers table */}
      <div className="overflow-hidden rounded-lg border border-slate-200 dark:border-slate-700">
        <table className="min-w-full divide-y divide-slate-200 dark:divide-slate-700">
          <thead className="bg-slate-50 dark:bg-slate-800">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500">
                Origen
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500">
                Destino
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500">
                Tipo
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500">
                Estado
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500">
                Fecha
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-slate-500">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 bg-white dark:divide-slate-700 dark:bg-slate-900">
            {transfers.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-sm text-slate-500">
                  No hay transferencias registradas
                </td>
              </tr>
            ) : (
              transfers.map((transfer) => (
                <tr key={transfer.id}>
                  <td className="whitespace-nowrap px-4 py-3 text-sm">
                    {transfer.origin_type}
                  </td>
                  <td className="whitespace-nowrap px-4 py-3 text-sm">
                    {transfer.destination_type}
                  </td>
                  <td className="whitespace-nowrap px-4 py-3 text-sm">
                    {transfer.transfer_type}
                  </td>
                  <td className="whitespace-nowrap px-4 py-3 text-sm">
                    <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusColors[transfer.status]}`}>
                      {transfer.status}
                    </span>
                  </td>
                  <td className="whitespace-nowrap px-4 py-3 text-sm text-slate-500">
                    {new Date(transfer.created_at).toLocaleDateString()}
                  </td>
                  <td className="whitespace-nowrap px-4 py-3 text-right text-sm">
                    <button className="text-primary hover:underline">Ver</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
