import { useSyncStore } from '../stores/sync-store';
import { cn } from '../lib/utils';

export function ConnectionStatus() {
  const isOnline = useSyncStore((state) => state.isOnline);

  return (
    <div
      className={cn(
        'fixed bottom-4 right-4 flex items-center gap-2 rounded-full px-3 py-1.5 text-sm font-medium shadow-lg',
        isOnline
          ? 'bg-success text-white'
          : 'bg-warning text-white'
      )}
      role="status"
      aria-live="polite"
    >
      <span
        className={cn(
          'h-2 w-2 rounded-full',
          isOnline ? 'bg-white' : 'bg-white animate-pulse'
        )}
      />
      <span>{isOnline ? 'En línea' : 'Sin conexión'}</span>
    </div>
  );
}
