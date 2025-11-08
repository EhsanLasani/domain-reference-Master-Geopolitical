// ============================================================================
// FILE: country-mfe.config.js
// DOMAIN: Reference Master Geopolitical
// LAYER: Presentation Layer - Micro-Frontends
// PURPOSE: Module federation configuration for country MFE
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

const ModuleFederationPlugin = require('@module-federation/webpack');

module.exports = {
  mode: 'development',
  devServer: {
    port: 3001,
  },
  plugins: [
    new ModuleFederationPlugin({
      name: 'countryMfe',
      filename: 'remoteEntry.js',
      exposes: {
        './CountryList': './src/components/CountryList',
        './CountrySearch': './src/components/CountrySearch',
        './CountryDetails': './src/components/CountryDetails'
      },
      shared: {
        react: { singleton: true },
        'react-dom': { singleton: true },
        '@lasani/design-system': { singleton: true }
      }
    })
  ]
};