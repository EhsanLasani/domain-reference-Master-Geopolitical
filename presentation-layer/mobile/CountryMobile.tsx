import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';

interface Country {
  country_code: string;
  country_name: string;
  iso3_code?: string;
  is_active: boolean;
}

interface CountryMobileProps {
  tenantId: string;
  onCountrySelect?: (country: Country) => void;
}

export const CountryMobile: React.FC<CountryMobileProps> = ({ tenantId, onCountrySelect }) => {
  const [countries, setCountries] = useState<Country[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchCountries();
  }, [tenantId]);

  const fetchCountries = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/v1/countries', {
        headers: {
          'X-Tenant-ID': tenantId,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch countries');
      }

      const data = await response.json();
      setCountries(data.countries || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const renderCountry = ({ item }: { item: Country }) => (
    <TouchableOpacity
      style={[styles.countryItem, !item.is_active && styles.inactiveItem]}
      onPress={() => onCountrySelect?.(item)}
      accessibilityLabel={`Country: ${item.country_name}`}
      accessibilityRole="button"
    >
      <View style={styles.countryInfo}>
        <Text style={styles.countryCode}>{item.country_code}</Text>
        <Text style={styles.countryName}>{item.country_name}</Text>
        {item.iso3_code && (
          <Text style={styles.iso3Code}>ISO3: {item.iso3_code}</Text>
        )}
      </View>
      <View style={[styles.statusIndicator, item.is_active ? styles.active : styles.inactive]} />
    </TouchableOpacity>
  );

  if (loading) {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color="#007AFF" />
        <Text style={styles.loadingText}>Loading countries...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.centerContainer}>
        <Text style={styles.errorText}>Error: {error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={fetchCountries}>
          <Text style={styles.retryButtonText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Countries ({countries.length})</Text>
      <FlatList
        data={countries}
        renderItem={renderCountry}
        keyExtractor={(item) => item.country_code}
        style={styles.list}
        showsVerticalScrollIndicator={false}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  header: {
    fontSize: 24,
    fontWeight: 'bold',
    padding: 16,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  list: {
    flex: 1,
  },
  countryItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    backgroundColor: '#fff',
    marginVertical: 1,
    borderLeftWidth: 4,
    borderLeftColor: '#007AFF',
  },
  inactiveItem: {
    opacity: 0.6,
    borderLeftColor: '#ccc',
  },
  countryInfo: {
    flex: 1,
  },
  countryCode: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
  },
  countryName: {
    fontSize: 16,
    color: '#666',
    marginTop: 2,
  },
  iso3Code: {
    fontSize: 12,
    color: '#999',
    marginTop: 2,
  },
  statusIndicator: {
    width: 12,
    height: 12,
    borderRadius: 6,
  },
  active: {
    backgroundColor: '#4CAF50',
  },
  inactive: {
    backgroundColor: '#F44336',
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  errorText: {
    fontSize: 16,
    color: '#F44336',
    textAlign: 'center',
    marginBottom: 20,
  },
  retryButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 8,
  },
  retryButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});