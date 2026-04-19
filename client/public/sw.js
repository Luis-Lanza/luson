/**
 * Service Worker Placeholder
 *
 * This service worker will be expanded in Phase 2 to handle:
 * - Static asset caching
 * - API request caching for offline support
 * - Background sync for pending changes
 * - Push notifications
 */

const CACHE_NAME = 'battery-pos-v1';
const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/manifest.json',
  // Assets will be added here
];

// Install event - cache static assets
self.addEventListener('install', (event) => {
  console.log('[SW] Service Worker installing...');

  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      console.log('[SW] Caching static assets');
      return cache.addAll(STATIC_ASSETS);
    })
  );

  self.skipWaiting();
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  console.log('[SW] Service Worker activating...');

  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => {
            console.log('[SW] Deleting old cache:', name);
            return caches.delete(name);
          })
      );
    })
  );

  self.clients.claim();
});

// Fetch event - serve from cache or network
self.addEventListener('fetch', (event) => {
  // Placeholder - will implement caching strategy in Phase 2
  console.log('[SW] Fetch:', event.request.url);
});

// Background sync event - handle pending changes
self.addEventListener('sync', (event) => {
  console.log('[SW] Background sync:', event.tag);
  // TODO: Implement background sync in Phase 2
});

// Push notification event
self.addEventListener('push', (event) => {
  console.log('[SW] Push received:', event);
  // TODO: Implement push notifications in Phase 2
});
