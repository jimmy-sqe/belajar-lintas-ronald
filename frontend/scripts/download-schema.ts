import dotenv from 'dotenv';
import fs from 'fs';
import https from 'https';
import path from 'path';
import { fileURLToPath } from 'url';

dotenv.config();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const urlStr = process.env.NEXT_OPENAPI_URL;
const username = process.env.NEXT_OPENAPI_USERNAME;
const password = process.env.NEXT_OPENAPI_PASSWORD;

if (!urlStr) {
  console.error('Error: NEXT_OPENAPI_URL is missing in environment variables.');
  process.exit(1);
}

const url = new URL(urlStr);

interface RequestOptions extends https.RequestOptions {
  hostname: string;
  port: number;
  path: string;
  method: string;
  headers: Record<string, string>;
}

const options: RequestOptions = {
  hostname: url.hostname,
  port: url.port ? parseInt(url.port, 10) : url.protocol === 'https:' ? 443 : 80,
  path: url.pathname + url.search,
  method: 'GET',
  headers: {
    'User-Agent': 'Mozilla/5.0',
    Accept: 'application/json, text/plain, */*'
  }
};

if (username && password) {
  const auth = Buffer.from(`${username}:${password}`).toString('base64');
  options.headers['Authorization'] = `Basic ${auth}`;
}

const req = https.request(options, (res) => {
  if (!res.statusCode || res.statusCode < 200 || res.statusCode >= 300) {
    console.error(`Error: Failed to fetch API schema. Status Code: ${res.statusCode}`);
    if (res.statusCode === 401 || res.statusCode === 403) {
      console.error(
        'Authentication failed. Check your NEXT_OPENAPI_USERNAME and NEXT_OPENAPI_PASSWORD.'
      );
    }
    process.exit(1);
  }

  let data = '';
  res.on('data', (chunk: Buffer) => {
    data += chunk.toString();
  });

  res.on('end', () => {
    try {
      // Validate JSON
      interface OpenAPISchema {
        paths?: Record<
          string,
          Record<
            string,
            {
              parameters?: Array<{ name: string }>;
            }
          >
        >;
        [key: string]: unknown;
      }

      const schemaJson: OpenAPISchema = JSON.parse(data);

      // Remove x-utc-offset requirement from all endpoints
      if (schemaJson.paths) {
        for (const pathKey of Object.keys(schemaJson.paths)) {
          for (const methodKey of Object.keys(schemaJson.paths[pathKey])) {
            const methodObj = schemaJson.paths[pathKey][methodKey];
            if (methodObj.parameters && Array.isArray(methodObj.parameters)) {
              methodObj.parameters = methodObj.parameters.filter((p) => p.name !== 'x-utc-offset');
            }
          }
        }
      }

      const outDir = path.join(__dirname, '..', 'src', 'openapi');
      if (!fs.existsSync(outDir)) {
        fs.mkdirSync(outDir, { recursive: true });
      }

      const outPath = path.join(outDir, 'openapi.json');
      fs.writeFileSync(outPath, JSON.stringify(schemaJson, null, 2), 'utf8');

      console.log('Successfully downloaded openapi.json to src/openapi/openapi.json');
    } catch (e) {
      const errorMessage = e instanceof Error ? e.message : String(e);
      console.error('Error: Fetched data is not valid JSON.', errorMessage);
      process.exit(1);
    }
  });
});

req.on('error', (e: Error) => {
  console.error(`Error: Network failure. ${e.message}`);
  process.exit(1);
});

req.end();
