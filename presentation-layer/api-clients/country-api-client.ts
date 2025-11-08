// ============================================================================
// FILE: country-api-client.ts
// DOMAIN: Reference Master Geopolitical
// LAYER: Presentation Layer - API Clients
// PURPOSE: Generated API client with error handling and retry logic
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

// Generated types from OpenAPI spec
export interface Country {
  country_id: string;
  country_code: string;
  country_name: string;
  iso3_code?: string;
  official_name?: string;
  is_active: boolean;
}

export interface CountryFilter {
  is_active?: boolean;
  continent_code?: string;
  limit?: number;
  offset?: number;
}

export interface CountriesResponse {
  data: Country[];
  total: number;
  has_more: boolean;
}

// API Error class
export class ApiError extends Error {
  constructor(
    public status: number,
    public body: string,
    public response?: Response
  ) {
    super(`API Error ${status}: ${body}`);
    this.name = 'ApiError';
  }
}

// Generated API client
export class CountryApiClient {
  private basePath: string;
  private retryPolicy = {
    maxAttempts: 3,
    backoffMs: 1000
  };

  constructor(basePath = '/api') {
    this.basePath = basePath;
  }

  async getCountries(filter: CountryFilter = {}): Promise<CountriesResponse> {
    return this.withRetry(async () => {
      const queryParams = new URLSearchParams();
      
      if (filter.is_active !== undefined) {
        queryParams.set('is_active', String(filter.is_active));
      }
      if (filter.continent_code) {
        queryParams.set('continent_code', filter.continent_code);
      }
      if (filter.limit) {
        queryParams.set('limit', String(filter.limit));
      }
      if (filter.offset) {
        queryParams.set('offset', String(filter.offset));
      }

      const response = await fetch(`${this.basePath}/countries?${queryParams}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new ApiError(response.status, await response.text(), response);
      }

      return await response.json();
    });
  }

  async getCountryByCode(code: string): Promise<Country> {
    return this.withRetry(async () => {
      const response = await fetch(`${this.basePath}/countries/${code}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new ApiError(response.status, await response.text(), response);
      }

      return await response.json();
    });
  }

  private async withRetry<T>(operation: () => Promise<T>): Promise<T> {
    let lastError: Error;

    for (let attempt = 1; attempt <= this.retryPolicy.maxAttempts; attempt++) {
      try {
        return await operation();
      } catch (error) {
        lastError = error as Error;

        if (attempt < this.retryPolicy.maxAttempts && this.isRetryable(error)) {
          await this.delay(this.retryPolicy.backoffMs * attempt);
          continue;
        }

        throw error;
      }
    }

    throw lastError!;
  }

  private isRetryable(error: any): boolean {
    return error instanceof ApiError && 
           (error.status >= 500 || error.status === 429);
  }

  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}