export function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-surface">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
        <h1 className="mb-6 text-center text-2xl font-bold text-text">
          Battery POS
        </h1>
        <p className="mb-4 text-center text-text-muted">
          Inicia sesión para continuar
        </p>
        <div className="space-y-4">
          <input
            type="email"
            placeholder="Email"
            className="w-full rounded-md border border-border px-4 py-2 focus:border-primary focus:outline-none"
          />
          <input
            type="password"
            placeholder="Contraseña"
            className="w-full rounded-md border border-border px-4 py-2 focus:border-primary focus:outline-none"
          />
          <button className="w-full rounded-md bg-primary px-4 py-2 font-medium text-white hover:bg-primary-dark">
            Iniciar sesión
          </button>
        </div>
      </div>
    </div>
  );
}
