/**
 * api/index.ts - API Module Entry Point
 *
 * This file re-exports everything from the api module, allowing other files
 * to import from './api' instead of './api/types' or './api/client'.
 *
 * This is a common pattern called a "barrel file" - it provides a single
 * import point for a module with multiple files.
 *
 * Usage:
 *   // Instead of:
 *   import { User } from './api/types';
 *   import { auth } from './api/client';
 *
 *   // You can do:
 *   import { User, auth } from './api';
 */

export * from './types';
export * from './client';
