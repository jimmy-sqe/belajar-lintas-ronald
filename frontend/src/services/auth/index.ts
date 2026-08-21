/**
 * Template file: edited by the SDLC plugin at scaffold time to point
 * to the chosen auth-model option. In the monorepo, the default is
 * jwt-refresh (the "primary" demo model).
 *
 * Plugin scaffold replaces the literal "./jwt-refresh" below with the
 * chosen option and deletes the inactive option folders.
 *
 * App code should import from "@/services/auth" only — never from
 * "@/services/auth/<option>/...".
 */
import adapter from './jwt-refresh';
export default adapter;
export * from './jwt-refresh/types';
export * from './types';
