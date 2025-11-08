// ============================================================================
// FILE: country.e2e.spec.ts
// DOMAIN: Reference Master Geopolitical
// LAYER: Presentation Layer - E2E Tests
// PURPOSE: End-to-end testing with Page Object Model
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

import { test, expect, Page } from '@playwright/test';

class CountryPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto('/countries');
  }

  async searchCountry(code: string) {
    await this.page.fill('[data-testid="country-search"]', code);
    await this.page.click('[data-testid="search-button"]');
  }

  async getCountryName() {
    return await this.page.textContent('[data-testid="country-name"]');
  }

  async isCountryActive() {
    const badge = await this.page.locator('[data-testid="country-status"]');
    return (await badge.textContent()) === 'Active';
  }
}

test.describe('Country Management', () => {
  let countryPage: CountryPage;

  test.beforeEach(async ({ page }) => {
    countryPage = new CountryPage(page);
    await countryPage.goto();
  });

  test('should display country information', async () => {
    await countryPage.searchCountry('US');
    
    const countryName = await countryPage.getCountryName();
    expect(countryName).toBe('United States');
    
    const isActive = await countryPage.isCountryActive();
    expect(isActive).toBe(true);
  });

  test('should meet performance budget', async ({ page }) => {
    const startTime = Date.now();
    await countryPage.goto();
    const loadTime = Date.now() - startTime;
    
    expect(loadTime).toBeLessThan(2500); // LCP budget
  });

  test('should be accessible', async ({ page }) => {
    await countryPage.goto();
    
    // Check for proper heading structure
    const h1 = await page.locator('h1').count();
    expect(h1).toBeGreaterThan(0);
    
    // Check for skip links
    const skipLink = await page.locator('[href="#main-content"]').count();
    expect(skipLink).toBeGreaterThan(0);
  });
});