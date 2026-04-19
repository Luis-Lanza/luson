import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { App } from './App';

const testRouter = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        index: true,
        element: <div data-testid="test-content">Test Content</div>,
      },
    ],
  },
]);

describe('App', () => {
  it('renders without crashing', () => {
    render(<RouterProvider router={testRouter} />);
    expect(screen.getByTestId('test-content')).toBeInTheDocument();
  });
});
