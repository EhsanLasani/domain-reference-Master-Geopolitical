// Enterprise Business Central JavaScript
class BusinessCentral {
    constructor() {
        this.countries = [];
        this.selectedRecords = [];
        this.currentEntity = 'countries';
        this.editingRecord = null;
        this.sortField = null;
        this.sortDirection = 'asc';
        
        this.init();
    }

    init() {
        this.currentEntity = 'countries';
        this.currentDomain = 'geopolitical';
        this.setupEventListeners();
        this.loadData();
        this.updateUI();
    }

    setupEventListeners() {
        // Only add listeners for elements that exist
        const selectAll = document.getElementById('select-all');
        if (selectAll) {
            selectAll.addEventListener('change', (e) => this.selectAll(e.target.checked));
        }

        const countryForm = document.getElementById('country-form');
        if (countryForm) {
            countryForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.saveRecord();
            });
        }

        // Grid sorting
        document.querySelectorAll('.sortable').forEach(header => {
            header.addEventListener('click', (e) => this.sortData(e.currentTarget.dataset.field));
        });
    }

    switchRibbonTab(tabName) {
        // Update active tab
        document.querySelectorAll('.ribbon-tab').forEach(tab => tab.classList.remove('active'));
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');

        // Update active panel
        document.querySelectorAll('.ribbon-panel').forEach(panel => panel.classList.remove('active'));
        document.getElementById(`${tabName}-panel`).classList.add('active');
    }

    switchEntity(entityName) {
        if (!entityName) return;

        this.currentEntity = entityName;
        
        // Update entity name in toolbar immediately
        const entityNameEl = document.getElementById('current-entity');
        if (entityNameEl) {
            entityNameEl.textContent = this.getEntityDisplayName(entityName);
        }
        
        // Update breadcrumb
        const breadcrumb = document.querySelector('.breadcrumb .current');
        if (breadcrumb) {
            breadcrumb.textContent = this.getEntityDisplayName(entityName);
        }
        
        // Update grid headers based on entity
        this.updateGridHeaders(entityName);
        
        // Reset selected records
        this.selectedRecords = [];
        this.updateSelectedCount();

        // Load entity data
        this.loadData();
    }

    getEntityDisplayName(entityName) {
        const names = {
            'countries': 'Countries',
            'regions': 'Regions', 
            'languages': 'Languages',
            'timezones': 'Timezones',
            'subdivisions': 'Subdivisions',
            'locales': 'Locales',
            'currencies': 'Currencies',
            'fx_rates': 'FX Rates',
            'uoms': 'Units of Measure',
            'incoterms': 'INCOTERMS',
            'payment_terms': 'Payment Terms',
            'tenants': 'Tenants',
            'subscriptions': 'Subscriptions',
            'billing_accounts': 'Billing Accounts',
            'invoices': 'Invoices',
            'payments': 'Payments',
            'users': 'Users',
            'roles': 'Roles',
            'permissions': 'Permissions',
            'sessions': 'Sessions'
        };
        return names[entityName] || entityName;
    }

    async loadData() {
        console.log('Loading data for entity:', this.currentEntity);
        this.showLoading(true);
        
        try {
            let data = [];
            
            // Try API first for all entities
            try {
                const apiEndpoint = this.getApiEndpoint(this.currentEntity);
                console.log('Fetching from API:', apiEndpoint);
                const response = await fetch(apiEndpoint);
                
                if (response.ok) {
                    const result = await response.json();
                    console.log('API Response:', result);
                    data = result.countries || result[this.currentEntity] || result.data || result || [];
                } else {
                    throw new Error(`API ${response.status}: ${response.statusText}`);
                }
            } catch (apiError) {
                console.warn('API failed, using mock data:', apiError.message);
                data = this.getMockData(this.currentEntity);
            }
            
            this.countries = Array.isArray(data) ? data : [];
            console.log('Data loaded:', this.countries.length, 'records');
            
            this.renderGrid();
            this.updateRecordCount();
            this.showNotification(`Loaded ${this.countries.length} ${this.getEntityDisplayName(this.currentEntity)}`, 'success');
            
        } catch (error) {
            console.error('Error loading data:', error);
            this.showNotification(`Error loading data: ${error.message}`, 'error');
            this.countries = [];
            this.renderGrid();
            this.updateRecordCount();
        } finally {
            this.showLoading(false);
        }
    }

    getApiEndpoint(entityName) {
        const endpoints = {
            'countries': '/api/v1/countries',
            'languages': '/api/v1/languages',
            'timezones': '/api/v1/timezones',
            'subdivisions': '/api/v1/subdivisions',
            'locales': '/api/v1/locales',
            'regions': '/api/v1/regions',
            'currencies': '/api/v1/currencies',
            'fx_rates': '/api/v1/fx-rates',
            'uoms': '/api/v1/uoms',
            'incoterms': '/api/v1/incoterms',
            'payment_terms': '/api/v1/payment-terms'
        };
        return endpoints[entityName] || `/api/v1/${entityName}`;
    }
    
    getMockData(entityName) {
        const mockData = {
            'languages': [
                { language_code: 'en', language_name: 'English', native_name: 'English', direction: 'LTR', iso_code: 'en-US', is_active: true },
                { language_code: 'es', language_name: 'Spanish', native_name: 'Español', direction: 'LTR', iso_code: 'es-ES', is_active: true },
                { language_code: 'fr', language_name: 'French', native_name: 'Français', direction: 'LTR', iso_code: 'fr-FR', is_active: true },
                { language_code: 'de', language_name: 'German', native_name: 'Deutsch', direction: 'LTR', iso_code: 'de-DE', is_active: true },
                { language_code: 'zh', language_name: 'Chinese', native_name: '中文', direction: 'LTR', iso_code: 'zh-CN', is_active: true },
                { language_code: 'ar', language_name: 'Arabic', native_name: 'العربية', direction: 'RTL', iso_code: 'ar-SA', is_active: true }
            ],
            'timezones': [
                { timezone_code: 'UTC', timezone_name: 'Coordinated Universal Time', utc_offset: '+00:00', region: 'Global', is_active: true },
                { timezone_code: 'EST', timezone_name: 'Eastern Standard Time', utc_offset: '-05:00', region: 'North America', is_active: true },
                { timezone_code: 'PST', timezone_name: 'Pacific Standard Time', utc_offset: '-08:00', region: 'North America', is_active: true },
                { timezone_code: 'GMT', timezone_name: 'Greenwich Mean Time', utc_offset: '+00:00', region: 'Europe', is_active: true },
                { timezone_code: 'CET', timezone_name: 'Central European Time', utc_offset: '+01:00', region: 'Europe', is_active: true },
                { timezone_code: 'JST', timezone_name: 'Japan Standard Time', utc_offset: '+09:00', region: 'Asia', is_active: true }
            ],
            'subdivisions': [
                { subdivision_code: 'US-CA', subdivision_name: 'California', country_code: 'US', subdivision_type: 'State', is_active: true },
                { subdivision_code: 'US-NY', subdivision_name: 'New York', country_code: 'US', subdivision_type: 'State', is_active: true },
                { subdivision_code: 'GB-ENG', subdivision_name: 'England', country_code: 'GB', subdivision_type: 'Country', is_active: true },
                { subdivision_code: 'DE-BY', subdivision_name: 'Bavaria', country_code: 'DE', subdivision_type: 'State', is_active: true },
                { subdivision_code: 'IN-MH', subdivision_name: 'Maharashtra', country_code: 'IN', subdivision_type: 'State', is_active: true }
            ],
            'locales': [
                { locale_code: 'en-US', locale_name: 'English (United States)', language_code: 'en', country_code: 'US', is_active: true },
                { locale_code: 'en-GB', locale_name: 'English (United Kingdom)', language_code: 'en', country_code: 'GB', is_active: true },
                { locale_code: 'fr-FR', locale_name: 'French (France)', language_code: 'fr', country_code: 'FR', is_active: true },
                { locale_code: 'de-DE', locale_name: 'German (Germany)', language_code: 'de', country_code: 'DE', is_active: true },
                { locale_code: 'es-ES', locale_name: 'Spanish (Spain)', language_code: 'es', country_code: 'ES', is_active: true }
            ],
            'currencies': [
                { currency_code: 'USD', currency_name: 'US Dollar', symbol: '$', numeric_code: 840, decimal_places: 2, is_active: true },
                { currency_code: 'EUR', currency_name: 'Euro', symbol: '€', numeric_code: 978, decimal_places: 2, is_active: true },
                { currency_code: 'GBP', currency_name: 'British Pound', symbol: '£', numeric_code: 826, decimal_places: 2, is_active: true },
                { currency_code: 'JPY', currency_name: 'Japanese Yen', symbol: '¥', numeric_code: 392, decimal_places: 0, is_active: true },
                { currency_code: 'INR', currency_name: 'Indian Rupee', symbol: '₹', numeric_code: 356, decimal_places: 2, is_active: true }
            ],
            'fx_rates': [
                { from_currency: 'USD', to_currency: 'EUR', rate: 0.85, effective_date: '2024-01-15', is_active: true },
                { from_currency: 'USD', to_currency: 'GBP', rate: 0.79, effective_date: '2024-01-15', is_active: true },
                { from_currency: 'EUR', to_currency: 'GBP', rate: 0.93, effective_date: '2024-01-15', is_active: true },
                { from_currency: 'USD', to_currency: 'JPY', rate: 149.50, effective_date: '2024-01-15', is_active: true }
            ],
            'uoms': [
                { uom_code: 'KG', uom_name: 'Kilogram', category: 'Weight', base_unit: true, is_active: true },
                { uom_code: 'LB', uom_name: 'Pound', category: 'Weight', base_unit: false, is_active: true },
                { uom_code: 'M', uom_name: 'Meter', category: 'Length', base_unit: true, is_active: true },
                { uom_code: 'FT', uom_name: 'Foot', category: 'Length', base_unit: false, is_active: true }
            ]
        };
        return mockData[entityName] || [];
    }

    renderGrid() {
        console.log('Rendering grid with', this.countries.length, 'records');
        const tbody = document.getElementById('grid-body');
        if (!tbody) {
            console.error('Grid body element not found');
            return;
        }
        
        tbody.innerHTML = '';

        if (this.countries.length === 0) {
            const emptyRow = document.createElement('tr');
            emptyRow.innerHTML = '<td colspan="8" style="text-align: center; padding: 20px; color: #666;">No data available</td>';
            tbody.appendChild(emptyRow);
            return;
        }

        this.countries.forEach((record, index) => {
            const row = this.createGridRow(record, index);
            tbody.appendChild(row);
        });
        
        console.log('Grid rendered successfully');
    }

    createGridRow(record, index) {
        const row = document.createElement('tr');
        row.dataset.index = index;
        
        // Create appropriate columns based on entity type
        let columns = '';
        
        switch(this.currentEntity) {
            case 'countries':
                columns = `
                    <td class="editable" data-field="country_code">${record.country_code || ''}</td>
                    <td class="editable" data-field="country_name">${record.country_name || ''}</td>
                    <td class="editable" data-field="iso3_code">${record.iso3_code || ''}</td>
                    <td class="editable" data-field="continent_code">${record.continent_code || ''}</td>
                    <td class="editable" data-field="capital_city">${record.capital_city || ''}</td>
                `;
                break;
            case 'languages':
                columns = `
                    <td class="editable" data-field="language_code">${record.language_code || ''}</td>
                    <td class="editable" data-field="language_name">${record.language_name || ''}</td>
                    <td class="editable" data-field="native_name">${record.native_name || ''}</td>
                    <td class="editable" data-field="direction">${record.direction || ''}</td>
                    <td class="editable" data-field="iso_code">${record.iso_code || ''}</td>
                `;
                break;
            case 'timezones':
                columns = `
                    <td class="editable" data-field="timezone_code">${record.timezone_code || ''}</td>
                    <td class="editable" data-field="timezone_name">${record.timezone_name || ''}</td>
                    <td class="editable" data-field="utc_offset">${record.utc_offset || ''}</td>
                    <td class="editable" data-field="region">${record.region || ''}</td>
                    <td></td>
                `;
                break;
            case 'subdivisions':
                columns = `
                    <td class="editable" data-field="subdivision_code">${record.subdivision_code || ''}</td>
                    <td class="editable" data-field="subdivision_name">${record.subdivision_name || ''}</td>
                    <td class="editable" data-field="country_code">${record.country_code || ''}</td>
                    <td class="editable" data-field="subdivision_type">${record.subdivision_type || ''}</td>
                    <td></td>
                `;
                break;
            case 'locales':
                columns = `
                    <td class="editable" data-field="locale_code">${record.locale_code || ''}</td>
                    <td class="editable" data-field="locale_name">${record.locale_name || ''}</td>
                    <td class="editable" data-field="language_code">${record.language_code || ''}</td>
                    <td class="editable" data-field="country_code">${record.country_code || ''}</td>
                    <td></td>
                `;
                break;
            case 'regions':
                columns = `
                    <td class="editable" data-field="region_code">${record.region_code || ''}</td>
                    <td class="editable" data-field="region_name">${record.region_name || ''}</td>
                    <td class="editable" data-field="region_type">${record.region_type || ''}</td>
                    <td></td>
                    <td></td>
                `;
                break;
            case 'currencies':
                columns = `
                    <td class="editable" data-field="currency_code">${record.currency_code || ''}</td>
                    <td class="editable" data-field="currency_name">${record.currency_name || ''}</td>
                    <td class="editable" data-field="symbol">${record.symbol || ''}</td>
                    <td class="editable" data-field="numeric_code">${record.numeric_code || ''}</td>
                    <td class="editable" data-field="decimal_places">${record.decimal_places || ''}</td>
                `;
                break;
            default:
                // Generic columns for unknown entities
                const keys = Object.keys(record).filter(key => key !== 'is_active');
                const displayKeys = keys.slice(0, 5);
                columns = displayKeys.map(key => `<td class="editable" data-field="${key}">${record[key] || ''}</td>`).join('');
                // Pad with empty columns if needed
                while (displayKeys.length < 5) {
                    columns += '<td></td>';
                    displayKeys.push('');
                }
        }
        
        row.innerHTML = `
            <td class="select-col">
                <input type="checkbox" onchange="app.toggleSelection(${index}, this.checked)">
            </td>
            ${columns}
            <td>
                <span class="status-badge ${record.is_active ? 'status-active' : 'status-inactive'}">
                    ${record.is_active ? 'Active' : 'Inactive'}
                </span>
            </td>
            <td class="actions-col">
                <button class="action-btn" onclick="app.editRecord(${index})" title="Edit">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn" onclick="app.deleteRecord(${index})" title="Delete">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;

        // Add inline editing
        row.querySelectorAll('.editable').forEach(cell => {
            cell.addEventListener('dblclick', () => this.startInlineEdit(cell, record));
        });

        return row;
    }

    startInlineEdit(cell, record) {
        const field = cell.dataset.field;
        const currentValue = cell.textContent.trim();
        
        const input = document.createElement('input');
        input.type = 'text';
        input.value = currentValue;
        input.style.width = '100%';
        input.style.border = '1px solid #4a90e2';
        input.style.padding = '2px 4px';
        input.style.fontSize = '11px';
        
        input.addEventListener('blur', () => this.finishInlineEdit(cell, input, record, field));
        input.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.finishInlineEdit(cell, input, record, field);
            } else if (e.key === 'Escape') {
                cell.textContent = currentValue;
            }
        });
        
        cell.innerHTML = '';
        cell.appendChild(input);
        input.focus();
        input.select();
    }

    async finishInlineEdit(cell, input, record, field) {
        const newValue = input.value.trim();
        const oldValue = record[field];
        
        if (newValue === oldValue) {
            cell.textContent = newValue;
            return;
        }

        try {
            // Update record
            record[field] = newValue;
            
            if (this.currentEntity === 'countries') {
                const response = await fetch(`/api/v1/countries/${record.country_code}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(record)
                });
                
                if (!response.ok) {
                    throw new Error('Update failed');
                }
            }
            
            cell.textContent = newValue;
            this.showNotification('Record updated successfully', 'success');
        } catch (error) {
            record[field] = oldValue;
            cell.textContent = oldValue;
            this.showNotification('Error updating record', 'error');
        }
    }

    toggleSelection(index, selected) {
        if (selected) {
            this.selectedRecords.push(index);
        } else {
            this.selectedRecords = this.selectedRecords.filter(i => i !== index);
        }
        
        this.updateSelectedCount();
        this.updateGridRowSelection();
    }

    selectAll(selected) {
        const checkboxes = document.querySelectorAll('#grid-body input[type="checkbox"]');
        checkboxes.forEach((checkbox, index) => {
            checkbox.checked = selected;
            this.toggleSelection(index, selected);
        });
    }

    updateGridRowSelection() {
        const rows = document.querySelectorAll('#grid-body tr');
        rows.forEach((row, index) => {
            if (this.selectedRecords.includes(index)) {
                row.classList.add('selected');
            } else {
                row.classList.remove('selected');
            }
        });
    }

    sortData(field) {
        if (this.sortField === field) {
            this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            this.sortField = field;
            this.sortDirection = 'asc';
        }

        this.countries.sort((a, b) => {
            let aVal = a[field] || '';
            let bVal = b[field] || '';
            
            if (typeof aVal === 'string') {
                aVal = aVal.toLowerCase();
                bVal = bVal.toLowerCase();
            }
            
            if (this.sortDirection === 'asc') {
                return aVal < bVal ? -1 : aVal > bVal ? 1 : 0;
            } else {
                return aVal > bVal ? -1 : aVal < bVal ? 1 : 0;
            }
        });

        this.renderGrid();
        this.updateSortIndicators(field);
    }

    updateSortIndicators(field) {
        document.querySelectorAll('.sortable i').forEach(icon => {
            icon.className = 'fas fa-sort';
        });
        
        const activeHeader = document.querySelector(`[data-field="${field}"] i`);
        if (activeHeader) {
            activeHeader.className = `fas fa-sort-${this.sortDirection === 'asc' ? 'up' : 'down'}`;
        }
    }

    searchRecords(query) {
        if (!query) {
            this.renderGrid();
            return;
        }

        const filtered = this.countries.filter(record => {
            return Object.values(record).some(value => 
                String(value).toLowerCase().includes(query.toLowerCase())
            );
        });

        const tbody = document.getElementById('grid-body');
        tbody.innerHTML = '';
        
        filtered.forEach((record, index) => {
            const row = this.createGridRow(record, index);
            tbody.appendChild(row);
        });
    }

    showLoading(show) {
        const loadingOverlay = document.getElementById('loading-overlay');
        if (loadingOverlay) {
            loadingOverlay.style.display = show ? 'flex' : 'none';
        } else {
            console.warn('Loading overlay element not found');
        }
    }

    updateRecordCount() {
        const recordCountEl = document.getElementById('record-count');
        if (recordCountEl) {
            recordCountEl.textContent = this.countries.length;
            console.log('Updated record count to:', this.countries.length);
        } else {
            console.warn('Record count element not found');
        }
    }

    updateSelectedCount() {
        document.getElementById('selected-count').textContent = this.selectedRecords.length;
    }

    updateUI() {
        this.updateRecordCount();
        this.updateSelectedCount();
    }

    updateGridHeaders(entityName) {
        const thead = document.querySelector('.grid-table thead tr');
        if (!thead) return;
        
        let headers = '';
        
        switch(entityName) {
            case 'countries':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="country_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="country_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="iso3_code">ISO3 <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="continent_code">Continent <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="capital_city">Capital <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            case 'languages':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="language_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="language_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="native_name">Native Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="direction">Direction <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="iso_code">ISO Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            case 'timezones':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="timezone_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="timezone_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="utc_offset">UTC Offset <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="region">Region <i class="fas fa-sort"></i></th>
                    <th></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            case 'subdivisions':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="subdivision_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="subdivision_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="country_code">Country <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="subdivision_type">Type <i class="fas fa-sort"></i></th>
                    <th></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            case 'locales':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="locale_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="locale_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="language_code">Language <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="country_code">Country <i class="fas fa-sort"></i></th>
                    <th></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            case 'regions':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="region_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="region_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="region_type">Type <i class="fas fa-sort"></i></th>
                    <th></th>
                    <th></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            case 'currencies':
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable" data-field="currency_code">Code <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="currency_name">Name <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="symbol">Symbol <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="numeric_code">Numeric <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="decimal_places">Decimals <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
                break;
            default:
                headers = `
                    <th class="select-col"><input type="checkbox" id="select-all"></th>
                    <th class="sortable">Field 1 <i class="fas fa-sort"></i></th>
                    <th class="sortable">Field 2 <i class="fas fa-sort"></i></th>
                    <th class="sortable">Field 3 <i class="fas fa-sort"></i></th>
                    <th class="sortable">Field 4 <i class="fas fa-sort"></i></th>
                    <th class="sortable">Field 5 <i class="fas fa-sort"></i></th>
                    <th class="sortable" data-field="is_active">Status <i class="fas fa-sort"></i></th>
                    <th class="actions-col">Actions</th>
                `;
        }
        
        thead.innerHTML = headers;
        
        // Re-attach event listeners
        document.querySelectorAll('.sortable').forEach(header => {
            header.addEventListener('click', (e) => this.sortData(e.currentTarget.dataset.field));
        });
        
        document.getElementById('select-all').addEventListener('change', (e) => this.selectAll(e.target.checked));
    }

    // Navigation functions
    toggleNavNode(nodeItem) {
        const children = nodeItem.parentNode.querySelector('.nav-children');
        const toggle = nodeItem.querySelector('.toggle');
        
        if (children && children.style.display === 'none') {
            children.style.display = 'block';
            if (toggle) toggle.className = 'fas fa-chevron-down toggle';
            nodeItem.parentNode.classList.add('expanded');
        } else if (children) {
            children.style.display = 'none';
            if (toggle) toggle.className = 'fas fa-chevron-right toggle';
            nodeItem.parentNode.classList.remove('expanded');
        }
    }

    switchDomain(domainName) {
        // Update domain context
        console.log('Switching to domain:', domainName);
        this.currentDomain = domainName;
        
        // Update breadcrumb
        const breadcrumb = document.querySelector('.breadcrumb');
        if (breadcrumb) {
            breadcrumb.innerHTML = `
                <span>${this.getDomainDisplayName(domainName)}</span>
                <i class="fas fa-chevron-right"></i>
                <span class="current">${this.getEntityDisplayName(this.currentEntity)}</span>
            `;
        }
        
        // Update tenant info to show current domain
        const tenantInfo = document.querySelector('.bc-tenant-info');
        if (tenantInfo) {
            tenantInfo.innerHTML = `
                <i class="fas fa-home"></i>
                <span class="bc-tenant-name">Default Tenant</span>
                <span class="bc-tenant-separator">|</span>
                <span class="bc-dashboard-link">${this.getDomainDisplayName(domainName)}</span>
            `;
        }
        
        this.showNotification(`Switched to ${this.getDomainDisplayName(domainName)}`, 'success');
    }
    
    getDomainDisplayName(domainName) {
        const names = {
            'geopolitical': 'Geopolitical Reference',
            'commerce': 'Commerce Reference', 
            'tenant': 'Tenant Management',
            'iam': 'Identity & Access',
            'reports': 'Reports'
        };
        return names[domainName] || domainName;
    }

    // Filter functions
    applyQuickFilter(filterType) {
        // Update active quick filter
        document.querySelectorAll('.quick-filter').forEach(btn => btn.classList.remove('active'));
        event.target.classList.add('active');
        
        let filteredData = [...this.countries];
        
        switch(filterType) {
            case 'recent':
                const weekAgo = new Date();
                weekAgo.setDate(weekAgo.getDate() - 7);
                filteredData = this.countries.filter(c => new Date(c.created_at) > weekAgo);
                break;
            case 'modified':
                const today = new Date().toDateString();
                filteredData = this.countries.filter(c => new Date(c.updated_at).toDateString() === today);
                break;
            case 'inactive':
                filteredData = this.countries.filter(c => !c.is_active);
                break;
            default:
                filteredData = this.countries;
        }
        
        this.renderFilteredGrid(filteredData);
    }
    
    applyFilters() {
        let filteredData = [...this.countries];
        
        // Status filter
        const statusCheckboxes = document.querySelectorAll('.filter-section:first-child input[type="checkbox"]:checked');
        if (statusCheckboxes.length < 3) { // Not "All" selected
            const activeChecked = Array.from(statusCheckboxes).some(cb => cb.parentNode.textContent.includes('Active'));
            const inactiveChecked = Array.from(statusCheckboxes).some(cb => cb.parentNode.textContent.includes('Inactive'));
            
            filteredData = filteredData.filter(c => {
                if (activeChecked && c.is_active) return true;
                if (inactiveChecked && !c.is_active) return true;
                return false;
            });
        }
        
        // Continent filter
        const continentCheckboxes = document.querySelectorAll('.filter-section:nth-child(2) input[type="checkbox"]:checked');
        if (continentCheckboxes.length < 7) { // Not "All" selected
            const selectedContinents = Array.from(continentCheckboxes)
                .map(cb => cb.parentNode.textContent.match(/\(([^)]+)\)/)?.[1])
                .filter(Boolean);
            
            if (selectedContinents.length > 0) {
                filteredData = filteredData.filter(c => selectedContinents.includes(c.continent_code));
            }
        }
        
        this.renderFilteredGrid(filteredData);
        this.showNotification(`Applied filters: ${filteredData.length} records found`, 'success');
    }
    
    clearFilters() {
        // Reset all checkboxes
        document.querySelectorAll('.filter-option input[type="checkbox"]').forEach(cb => {
            cb.checked = cb.parentNode.textContent.includes('All');
        });
        
        // Clear date inputs
        document.querySelectorAll('.filter-date').forEach(input => input.value = '');
        
        // Reset quick filters
        document.querySelectorAll('.quick-filter').forEach(btn => btn.classList.remove('active'));
        document.querySelector('.quick-filter').classList.add('active');
        
        this.renderGrid();
        this.showNotification('Filters cleared', 'success');
    }
    
    renderFilteredGrid(data) {
        const tbody = document.getElementById('grid-body');
        tbody.innerHTML = '';
        
        data.forEach((record, index) => {
            const row = this.createGridRow(record, index);
            tbody.appendChild(row);
        });
        
        document.getElementById('record-count').textContent = data.length;
    }

    showNotification(message, type) {
        // Simple notification - could be enhanced
        const notification = document.createElement('div');
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 12px 16px;
            background: ${type === 'error' ? '#f8d7da' : '#d4edda'};
            color: ${type === 'error' ? '#721c24' : '#155724'};
            border: 1px solid ${type === 'error' ? '#f5c6cb' : '#c3e6cb'};
            border-radius: 4px;
            font-size: 12px;
            z-index: 10000;
            box-shadow: 0 2px 8px rgba(0,0,0,0.15);
        `;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 3000);
    }
}

// Global navigation functions
function toggleNavNode(nodeItem) {
    app.toggleNavNode(nodeItem);
}



function applyQuickFilter(filterType) {
    app.applyQuickFilter(filterType);
}

function applyFilters() {
    app.applyFilters();
}

function clearFilters() {
    app.clearFilters();
}

function toggleFilterPanel() {
    const filterPanel = document.querySelector('.nav-pane');
    const btn = event.target.closest('.bc-panel-btn');
    
    if (filterPanel.style.display === 'none') {
        filterPanel.style.display = 'flex';
        btn.classList.add('active');
    } else {
        filterPanel.style.display = 'none';
        btn.classList.remove('active');
    }
}

function toggleInfoPanel() {
    // Create info panel if it doesn't exist
    let infoPanel = document.getElementById('info-panel');
    if (!infoPanel) {
        infoPanel = document.createElement('div');
        infoPanel.id = 'info-panel';
        infoPanel.className = 'info-panel';
        infoPanel.innerHTML = `
            <div class="info-header">
                <h3>Information</h3>
                <button onclick="toggleInfoPanel()"><i class="fas fa-times"></i></button>
            </div>
            <div class="info-content">
                <p>Select a record to view details</p>
            </div>
        `;
        document.querySelector('.main-content').appendChild(infoPanel);
    }
    
    const btn = event.target.closest('.bc-panel-btn');
    
    if (infoPanel.style.display === 'none' || !infoPanel.style.display) {
        infoPanel.style.display = 'flex';
        btn.classList.add('active');
    } else {
        infoPanel.style.display = 'none';
        btn.classList.remove('active');
    }
}

function bookmarkPage() {
    const currentEntity = document.getElementById('current-entity').textContent;
    app.showNotification(`Bookmarked: ${currentEntity}`, 'success');
}

function activateDomain(element, domainName) {
    // Toggle open state for clicked item
    if (element.classList.contains('open')) {
        element.classList.remove('open');
    } else {
        // Close other submenus
        document.querySelectorAll('.bc-nav-item').forEach(item => {
            if (item !== element && !item.classList.contains('pinned')) {
                item.classList.remove('open');
            }
        });
        element.classList.add('open');
    }
    
    // Update active state
    document.querySelectorAll('.bc-nav-item').forEach(item => item.classList.remove('active'));
    element.classList.add('active');
    
    // Switch domain
    app.switchDomain(domainName);
    
    // Prevent event bubbling
    event.stopPropagation();
}

function switchEntity(entityName) {
    if (!entityName) {
        console.error('No entity name provided');
        return;
    }
    
    console.log('Switching to entity:', entityName);
    
    // Use the app's switchEntity method for consistency
    app.switchEntity(entityName);
    
    app.showNotification(`Switched to ${app.getEntityDisplayName(entityName)}`, 'success');
}

function pinSubmenu(element, domainName) {
    const navItem = element.closest('.bc-nav-item');
    const pinBtn = element;
    const submenu = navItem.querySelector('.bc-submenu');
    
    // First unpin any other pinned submenus
    document.querySelectorAll('.bc-nav-item.pinned').forEach(item => {
        if (item !== navItem) {
            item.classList.remove('pinned');
            const otherPin = item.querySelector('.bc-pin-btn');
            if (otherPin) otherPin.classList.remove('pinned');
        }
    });
    
    if (navItem.classList.contains('pinned')) {
        // Unpin
        navItem.classList.remove('pinned');
        pinBtn.classList.remove('pinned');
        document.body.classList.remove('submenu-pinned');
        document.documentElement.style.removeProperty('--submenu-height');
        app.showNotification('Submenu unpinned - toolbar restored', 'success');
    } else {
        // Pin
        navItem.classList.add('pinned');
        pinBtn.classList.add('pinned');
        document.body.classList.add('submenu-pinned');
        
        // Calculate actual submenu height
        const submenuHeight = submenu.offsetHeight;
        document.documentElement.style.setProperty('--submenu-height', submenuHeight + 'px');
        
        app.showNotification('Submenu pinned - toolbar moved down', 'success');
    }
    
    event.stopPropagation();
}

function pinNavBox(element, entityName) {
    const pinBtn = element;
    
    if (pinBtn.classList.contains('pinned')) {
        pinBtn.classList.remove('pinned');
        app.showNotification(`${entityName} unpinned`, 'success');
    } else {
        pinBtn.classList.add('pinned');
        app.showNotification(`${entityName} pinned`, 'success');
    }
}

// Close submenus when clicking outside
document.addEventListener('click', function(event) {
    if (!event.target.closest('.bc-nav-item')) {
        document.querySelectorAll('.bc-nav-item').forEach(item => {
            if (!item.classList.contains('pinned')) {
                item.classList.remove('open');
            }
        });
    }
});

// Keep submenu open when clicking inside it
document.addEventListener('click', function(event) {
    if (event.target.closest('.bc-submenu')) {
        event.stopPropagation();
    }
});

// Initialize submenu state on page load
document.addEventListener('DOMContentLoaded', function() {
    // Check if any submenu is pinned and apply body class
    if (document.querySelector('.bc-nav-item.pinned')) {
        document.body.classList.add('submenu-pinned');
    }
});

// Global functions for ribbon and dialog actions
function addRecord() {
    app.editingRecord = null;
    document.getElementById('dialog-title').textContent = `Add ${app.getEntityDisplayName(app.currentEntity).slice(0, -1)}`;
    document.getElementById('country-form').reset();
    document.getElementById('form-dialog').style.display = 'flex';
}

function editRecord(index) {
    app.editingRecord = app.countries[index];
    document.getElementById('dialog-title').textContent = `Edit ${app.getEntityDisplayName(app.currentEntity).slice(0, -1)}`;
    
    // Populate form
    Object.keys(app.editingRecord).forEach(key => {
        const input = document.getElementById(key);
        if (input) {
            if (input.type === 'checkbox') {
                input.checked = app.editingRecord[key];
            } else {
                input.value = app.editingRecord[key] || '';
            }
        }
    });
    
    document.getElementById('form-dialog').style.display = 'flex';
}

function deleteRecord(index) {
    if (confirm('Are you sure you want to delete this record?')) {
        app.countries.splice(index, 1);
        app.renderGrid();
        app.updateRecordCount();
        app.showNotification('Record deleted successfully', 'success');
    }
}

function refreshData() {
    console.log('Refresh button clicked');
    app.loadData();
}

function exportData() {
    const csv = app.countries.map(record => Object.values(record).join(',')).join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${app.currentEntity}.csv`;
    a.click();
}

function importData() {
    app.showNotification('Import functionality coming soon', 'info');
}

function toggleView(viewType) {
    document.querySelectorAll('.view-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelector(`[data-view="${viewType}"]`).classList.add('active');
    
    if (viewType === 'form' && app.countries.length > 0) {
        editRecord(0);
    }
}

function searchRecords() {
    const query = document.getElementById('search-input').value;
    app.searchRecords(query);
}

async function saveRecord() {
    const formData = new FormData(document.getElementById('country-form'));
    const record = {};
    
    for (let [key, value] of formData.entries()) {
        record[key] = value;
    }
    
    // Add checkbox values
    record.is_active = document.getElementById('is_active').checked;
    
    try {
        if (app.editingRecord) {
            // Update existing
            Object.assign(app.editingRecord, record);
            
            if (app.currentEntity === 'countries') {
                const response = await fetch(`/api/v1/countries/${app.editingRecord.country_code}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(app.editingRecord)
                });
                
                if (!response.ok) throw new Error('Update failed');
            }
        } else {
            // Create new
            if (app.currentEntity === 'countries') {
                const response = await fetch('/api/v1/countries', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(record)
                });
                
                if (!response.ok) throw new Error('Create failed');
            }
            
            app.countries.push(record);
        }
        
        closeDialog();
        app.renderGrid();
        app.updateRecordCount();
        app.showNotification('Record saved successfully', 'success');
    } catch (error) {
        app.showNotification('Error saving record', 'error');
    }
}

function closeDialog() {
    document.getElementById('form-dialog').style.display = 'none';
}

// Initialize application when DOM is ready
let app;
document.addEventListener('DOMContentLoaded', function() {
    app = new BusinessCentral();
});