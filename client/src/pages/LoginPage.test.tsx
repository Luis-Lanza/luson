import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import { LoginPage } from './LoginPage';

// Mock auth store
const mockLogin = vi.fn();
const mockClearError = vi.fn();
const mockGetRedirectPath = vi.fn();

let mockStoreState = {
  login: mockLogin,
  isLoading: false,
  error: null as string | null,
  clearError: mockClearError,
  isAuthenticated: false,
  getRedirectPath: mockGetRedirectPath,
};

vi.mock('../stores/auth-store', () => ({
  useAuthStore: () => mockStoreState,
}));

// Mock navigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('LoginPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockStoreState = {
      login: mockLogin,
      isLoading: false,
      error: null,
      clearError: mockClearError,
      isAuthenticated: false,
      getRedirectPath: mockGetRedirectPath,
    };
    mockGetRedirectPath.mockReturnValue('/dashboard');
  });

  it('should render login form', () => {
    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    expect(screen.getByText('Battery POS')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Usuario')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Contraseña')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /iniciar sesión/i })).toBeInTheDocument();
  });

  it('should update input values on change', async () => {
    const user = userEvent.setup();
    
    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    const usernameInput = screen.getByPlaceholderText('Usuario');
    const passwordInput = screen.getByPlaceholderText('Contraseña');

    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123');

    expect(usernameInput).toHaveValue('testuser');
    expect(passwordInput).toHaveValue('password123');
  });

  it('should call login with credentials on submit', async () => {
    const user = userEvent.setup();
    mockLogin.mockResolvedValueOnce(undefined);

    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    await user.type(screen.getByPlaceholderText('Usuario'), 'testuser');
    await user.type(screen.getByPlaceholderText('Contraseña'), 'password123');
    await user.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith('testuser', 'password123');
    });
  });

  it('should show loading state while logging in', async () => {
    // Start with loading state to simulate during-login state
    mockStoreState.isLoading = true;

    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    expect(screen.getByText('Iniciando sesión...')).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('should display error message when login fails', async () => {
    mockStoreState.error = 'Credenciales inválidas';

    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    // Check that error is displayed
    const errorElement = screen.getByRole('alert');
    expect(errorElement).toHaveTextContent('Credenciales inválidas');
  });

  it('should navigate after successful login', async () => {
    const user = userEvent.setup();
    mockLogin.mockResolvedValueOnce(undefined);
    mockGetRedirectPath.mockReturnValueOnce('/pos');

    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    await user.type(screen.getByPlaceholderText('Usuario'), 'cashier');
    await user.type(screen.getByPlaceholderText('Contraseña'), 'pass');
    await user.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/pos');
    });
  });

  it('should navigate to dashboard for admin users', async () => {
    const user = userEvent.setup();
    mockLogin.mockResolvedValueOnce(undefined);
    mockGetRedirectPath.mockReturnValueOnce('/dashboard');

    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    await user.type(screen.getByPlaceholderText('Usuario'), 'admin');
    await user.type(screen.getByPlaceholderText('Contraseña'), 'admin123');
    await user.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
    });
  });

  it('should disable submit button when inputs are empty', () => {
    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    const button = screen.getByRole('button', { name: /iniciar sesión/i });
    expect(button).toBeDisabled();
  });

  it('should enable submit button when both inputs have values', async () => {
    const user = userEvent.setup();
    
    render(
      <BrowserRouter>
        <LoginPage />
      </BrowserRouter>
    );

    const usernameInput = screen.getByPlaceholderText('Usuario');
    const passwordInput = screen.getByPlaceholderText('Contraseña');
    const button = screen.getByRole('button', { name: /iniciar sesión/i });

    expect(button).toBeDisabled();

    await user.type(usernameInput, 'test');
    expect(button).toBeDisabled();

    await user.type(passwordInput, 'pass');
    expect(button).not.toBeDisabled();
  });
});
