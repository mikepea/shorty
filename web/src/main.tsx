/**
 * main.tsx - Application Entry Point
 *
 * This is where React "boots up" and connects to the HTML page.
 * Think of it as the ignition key that starts the whole application.
 *
 * Key concepts:
 * - createRoot: Creates a React "root" that manages rendering into a DOM element
 * - StrictMode: A development tool that warns about potential problems
 * - The '!' after getElementById is TypeScript saying "trust me, this exists"
 */

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

// Find the <div id="root"> in index.html and render our App inside it.
// This is the bridge between React and the actual HTML page.
createRoot(document.getElementById('root')!).render(
  // StrictMode runs extra checks in development (not in production).
  // It helps catch common mistakes like forgotten cleanup in useEffect.
  <StrictMode>
    <App />
  </StrictMode>,
)
