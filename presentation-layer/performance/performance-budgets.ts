// ============================================================================
// FILE: performance-budgets.ts
// DOMAIN: Reference Master Geopolitical
// LAYER: Presentation Layer - Performance
// PURPOSE: Performance and accessibility budgets
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

// Performance budget configuration
export const performanceBudgets = {
  // Core Web Vitals (in milliseconds)
  firstContentfulPaint: 1500,
  largestContentfulPaint: 2500,
  cumulativeLayoutShift: 0.1,
  firstInputDelay: 100,
  
  // Bundle size budgets (in bytes)
  initialBundle: 200 * 1024,       // 200KB
  chunkSize: 50 * 1024,           // 50KB per chunk
  totalAssets: 1024 * 1024,       // 1MB total
  
  // Network budgets (in milliseconds)
  apiResponseTime: 500,
  imageLoadTime: 1000,
  
  // Memory budgets
  heapSize: 50 * 1024 * 1024,     // 50MB
  domNodes: 1500
} as const;

// Accessibility budget configuration
export const a11yBudgets = {
  maxViolations: {
    critical: 0,
    serious: 0,
    moderate: 2,
    minor: 5
  },
  
  colorContrast: {
    normal: 4.5,    // WCAG AA normal text
    large: 3.0      // WCAG AA large text
  },
  
  requirements: {
    tabOrder: true,
    focusIndicators: true,
    skipLinks: true,
    semanticMarkup: true,
    altText: true,
    ariaLabels: true
  }
} as const;

// Performance monitoring utility
export class PerformanceMonitor {
  static measureWebVitals(): void {
    // Measure First Contentful Paint
    new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        if (entry.name === 'first-contentful-paint') {
          this.reportMetric('FCP', entry.startTime);
        }
      });
    }).observe({ entryTypes: ['paint'] });

    // Measure Largest Contentful Paint
    new PerformanceObserver((list) => {
      const entries = list.getEntries();
      const lastEntry = entries[entries.length - 1];
      this.reportMetric('LCP', lastEntry.startTime);
    }).observe({ entryTypes: ['largest-contentful-paint'] });
  }

  static reportMetric(name: string, value: number): void {
    const budgetKey = name.toLowerCase().replace(/\s+/g, '') as keyof typeof performanceBudgets;
    const budget = performanceBudgets[budgetKey];
    
    if (typeof budget === 'number' && value > budget) {
      console.warn(`Performance budget exceeded: ${name} = ${value}ms (budget: ${budget}ms)`);
      
      // Report to monitoring service
      this.sendToMonitoring(name, value, budget);
    }
  }

  private static sendToMonitoring(metric: string, value: number, budget: number): void {
    // Integration with monitoring service (e.g., DataDog, New Relic)
    if (typeof window !== 'undefined' && (window as any).analytics) {
      (window as any).analytics.track('Performance Budget Exceeded', {
        metric,
        value,
        budget,
        exceedBy: value - budget
      });
    }
  }
}

export type PerformanceBudgets = typeof performanceBudgets;
export type A11yBudgets = typeof a11yBudgets;