// ============================================================================
// FILE: CountryList.tsx
// DOMAIN: Reference Master Geopolitical
// LAYER: Presentation Layer - Micro-Frontends
// PURPOSE: Country list micro-frontend component
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

import React, { useState, useEffect } from 'react';
import { CountryApiClient, Country } from '../api-clients/country-api-client';
import { designTokens } from '../design-system/design-tokens';

interface CountryListProps {
  onCountrySelect?: (country: Country) => void;
  showActiveOnly?: boolean;
}

export const CountryList: React.FC<CountryListProps> = ({ 
  onCountrySelect, 
  showActiveOnly = true 
}) => {
  const [countries, setCountries] = useState<Country[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const apiClient = new CountryApiClient();

  useEffect(() => {
    const loadCountries = async () => {
      try {
        setLoading(true);
        const response = await apiClient.getCountries({ 
          is_active: showActiveOnly,
          limit: 50 
        });
        setCountries(response.data);
      } catch (err) {
        setError('Failed to load countries');
      } finally {
        setLoading(false);
      }
    };

    loadCountries();
  }, [showActiveOnly]);

  if (loading) return <div>Loading countries...</div>;
  if (error) return <div style={{ color: designTokens.colors.semantic.error }}>{error}</div>;

  return (
    <div style={{ 
      fontFamily: designTokens.typography.fontFamily.sans.join(', '),
      padding: designTokens.spacing.md 
    }}>
      <h2 style={{ 
        fontSize: designTokens.typography.fontSize.xl,
        marginBottom: designTokens.spacing.lg 
      }}>
        Countries
      </h2>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {countries.map((country) => (
          <li 
            key={country.country_id}
            onClick={() => onCountrySelect?.(country)}
            style={{
              padding: designTokens.spacing.sm,
              marginBottom: designTokens.spacing.xs,
              backgroundColor: designTokens.colors.neutral[50],
              borderRadius: designTokens.borderRadius.md,
              cursor: 'pointer',
              border: `1px solid ${designTokens.colors.neutral[100]}`
            }}
          >
            <strong>{country.country_code}</strong> - {country.country_name}
          </li>
        ))}
      </ul>
    </div>
  );
};