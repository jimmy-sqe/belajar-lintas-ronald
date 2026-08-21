import { registerOTel } from '@vercel/otel';

/**
 * OpenTelemetry server-side tracing. @vercel/otel reads
 * OTEL_EXPORTER_OTLP_ENDPOINT from the environment automatically.
 */
export async function register(): Promise<void> {
  registerOTel({ serviceName: process.env.NEXT_PUBLIC_OTEL_SERVICE_NAME || 'frontend-belajar-lintas-ronald' });
}
