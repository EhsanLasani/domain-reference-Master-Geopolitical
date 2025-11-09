// Enhanced Data Grid Component for LASANI Platform
class LasaniDataGrid {
    constructor(containerId, options = {}) {
        this.container = document.getElementById(containerId);
        this.options = {
            apiEndpoint: '/api/v1',
            entity: 'countries',
            editable: true,
            selectable: true,
            sortable: true,
            filterable: true,
            ...options
        };
        
        this.data = [];
        this.selectedRows = new Set();
        this.sortColumn = null;
        this.sortDirection = 'asc';
        this.filters = {};
        
        this.init();
    }
    
    init() {
        this.render();
        this.bindEvents();
        this.loadData();
    }
    
    render() {
        this.container.innerHTML = `
            <div class="data-grid-container">
                <div class="grid-toolbar">
                    <div class="toolbar-left">
                        <button class="btn primary" data-action="add">
                            <i class="fas fa-plus"></i> Add
                        </button>
                        <button class="btn" data-action="edit" disabled>
                            <i class="fas fa-edit"></i> Edit
                        </button>
                        <button class="btn" data-action="delete" disabled>
                            <i class="fas fa-trash"></i> Delete
                        </button>
                    </div>
                    <div class="toolbar-right">
                        <button class="btn" data-action="refresh">
                            <i class="fas fa-sync"></i> Refresh
                        </button>
                        <button class="btn" data-action="export">
                            <i class="fas fa-download"></i> Export
                        </button>
                    </div>
                </div>
                
                <div class="grid-content">
                    <table class="enhanced-table">
                        <thead id="gridHeader"></thead>
                        <tbody id="gridBody"></tbody>
                    </table>
                </div>
                
                <div class="grid-footer">
                    <div class="record-count">
                        <span id="recordCount">0 records</span>
                    </div>
                    <div class="pagination" id="pagination"></div>
                </div>
            </div>
        `;
    }
    
    bindEvents() {
        // Toolbar events
        this.container.addEventListener('click', (e) => {
            const action = e.target.closest('[data-action]')?.dataset.action;
            if (action) {
                this.handleAction(action, e);
            }
        });
        
        // Table events
        this.container.addEventListener('change', (e) => {
            if (e.target.type === 'checkbox') {
                this.handleSelection(e);
            }
        });
        
        // Sort events
        this.container.addEventListener('click', (e) => {
            if (e.target.closest('.sortable')) {
                this.handleSort(e.target.closest('.sortable'));
            }
        });
    }
    
    async loadData() {
        try {
            this.showLoading();
            const response = await fetch(`${this.options.apiEndpoint}/${this.options.entity}`);
            const result = await response.json();
            
            this.data = result[this.options.entity] || result.countries || [];
            this.renderTable();
            this.updateRecordCount();
        } catch (error) {
            console.error('Failed to load data:', error);
            this.showError('Failed to load data');
        }
    }
    
    renderTable() {
        if (this.data.length === 0) {
            this.renderEmptyState();
            return;
        }
        
        const headers = this.getHeaders();
        this.renderHeaders(headers);
        this.renderRows();
    }
    
    getHeaders() {
        if (this.options.entity === 'countries') {
            return [
                { key: 'select', label: '', sortable: false, width: '40px' },
                { key: 'country_code', label: 'Code', sortable: true },
                { key: 'country_name', label: 'Name', sortable: true },
                { key: 'iso3_code', label: 'ISO3', sortable: true },
                { key: 'continent_code', label: 'Continent', sortable: true },
                { key: 'capital_city', label: 'Capital', sortable: true },
                { key: 'is_active', label: 'Status', sortable: true },
                { key: 'actions', label: 'Actions', sortable: false, width: '100px' }
            ];
        }
        return [];
    }
    
    renderHeaders(headers) {
        const headerRow = headers.map(header => {
            if (header.key === 'select') {
                return `<th style="width: ${header.width || 'auto'}">
                    <input type="checkbox" id="selectAll">
                </th>`;
            }
            
            const sortClass = header.sortable ? 'sortable' : '';
            const sortIcon = this.getSortIcon(header.key);
            
            return `<th class="${sortClass}" data-column="${header.key}" style="width: ${header.width || 'auto'}">
                ${header.label} ${sortIcon}
            </th>`;
        }).join('');
        
        document.getElementById('gridHeader').innerHTML = `<tr>${headerRow}</tr>`;
    }
    
    renderRows() {
        const tbody = document.getElementById('gridBody');
        const rows = this.data.map((item, index) => {
            const isSelected = this.selectedRows.has(index);
            return `<tr data-index="${index}" ${isSelected ? 'class="selected"' : ''}>
                ${this.renderRowCells(item, index)}
            </tr>`;
        }).join('');
        
        tbody.innerHTML = rows;
    }
    
    renderRowCells(item, index) {
        if (this.options.entity === 'countries') {
            return `
                <td><input type="checkbox" class="row-select" data-index="${index}" ${this.selectedRows.has(index) ? 'checked' : ''}></td>
                <td class="editable" data-field="country_code">${item.country_code || ''}</td>
                <td class="editable" data-field="country_name">${item.country_name || ''}</td>
                <td class="editable" data-field="iso3_code">${item.iso3_code || ''}</td>
                <td class="editable" data-field="continent_code">${item.continent_code || ''}</td>
                <td class="editable" data-field="capital_city">${item.capital_city || ''}</td>
                <td>
                    <span class="status-badge ${item.is_active ? 'active' : 'inactive'}">
                        ${item.is_active ? 'Active' : 'Inactive'}
                    </span>
                </td>
                <td>
                    <button class="action-btn" data-action="edit-row" data-index="${index}" title="Edit">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="action-btn danger" data-action="delete-row" data-index="${index}" title="Delete">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            `;
        }
        return '';
    }
    
    handleAction(action, event) {
        switch (action) {
            case 'add':
                this.showAddDialog();
                break;
            case 'edit':
                this.editSelected();
                break;
            case 'delete':
                this.deleteSelected();
                break;
            case 'refresh':
                this.refresh(event.target);
                break;
            case 'export':
                this.exportData();
                break;
            case 'edit-row':
                this.editRow(parseInt(event.target.closest('[data-index]').dataset.index));
                break;
            case 'delete-row':
                this.deleteRow(parseInt(event.target.closest('[data-index]').dataset.index));
                break;
        }
    }
    
    handleSelection(event) {
        if (event.target.id === 'selectAll') {
            this.selectAll(event.target.checked);
        } else if (event.target.classList.contains('row-select')) {
            this.selectRow(parseInt(event.target.dataset.index), event.target.checked);
        }
        this.updateToolbarState();
    }
    
    selectAll(checked) {
        this.selectedRows.clear();
        if (checked) {
            this.data.forEach((_, index) => this.selectedRows.add(index));
        }
        
        document.querySelectorAll('.row-select').forEach(cb => {
            cb.checked = checked;
        });
        
        this.renderRows();
    }
    
    selectRow(index, checked) {
        if (checked) {
            this.selectedRows.add(index);
        } else {
            this.selectedRows.delete(index);
        }
        
        // Update select all checkbox
        const selectAll = document.getElementById('selectAll');
        selectAll.indeterminate = this.selectedRows.size > 0 && this.selectedRows.size < this.data.length;
        selectAll.checked = this.selectedRows.size === this.data.length;
    }
    
    updateToolbarState() {
        const hasSelection = this.selectedRows.size > 0;
        document.querySelector('[data-action="edit"]').disabled = !hasSelection;
        document.querySelector('[data-action="delete"]').disabled = !hasSelection;
    }
    
    refresh(button) {
        const icon = button.querySelector('i');
        icon.classList.add('fa-spin');
        
        this.loadData().finally(() => {
            setTimeout(() => icon.classList.remove('fa-spin'), 500);
        });
    }
    
    updateRecordCount() {
        document.getElementById('recordCount').textContent = `${this.data.length} records`;
    }
    
    showLoading() {
        document.getElementById('gridBody').innerHTML = `
            <tr><td colspan="8" style="text-align: center; padding: 2rem;">
                <div class="loading-spinner"></div> Loading...
            </td></tr>
        `;
    }
    
    showError(message) {
        document.getElementById('gridBody').innerHTML = `
            <tr><td colspan="8" style="text-align: center; padding: 2rem; color: var(--danger-color);">
                <i class="fas fa-exclamation-triangle"></i> ${message}
            </td></tr>
        `;
    }
    
    renderEmptyState() {
        document.getElementById('gridBody').innerHTML = `
            <tr><td colspan="8" style="text-align: center; padding: 3rem;">
                <i class="fas fa-inbox" style="font-size: 2rem; color: var(--text-secondary); margin-bottom: 1rem;"></i>
                <p>No records found</p>
                <button class="btn primary" data-action="add">Add First Record</button>
            </td></tr>
        `;
    }
    
    getSortIcon(column) {
        if (this.sortColumn === column) {
            return this.sortDirection === 'asc' ? 
                '<i class="fas fa-sort-up"></i>' : 
                '<i class="fas fa-sort-down"></i>';
        }
        return '<i class="fas fa-sort"></i>';
    }
    
    handleSort(header) {
        const column = header.dataset.column;
        if (this.sortColumn === column) {
            this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            this.sortColumn = column;
            this.sortDirection = 'asc';
        }
        
        this.sortData();
        this.renderTable();
    }
    
    sortData() {
        this.data.sort((a, b) => {
            const aVal = a[this.sortColumn] || '';
            const bVal = b[this.sortColumn] || '';
            
            if (this.sortDirection === 'asc') {
                return aVal.toString().localeCompare(bVal.toString());
            } else {
                return bVal.toString().localeCompare(aVal.toString());
            }
        });
    }
    
    showAddDialog() {
        console.log('Add dialog - to be implemented');
    }
    
    editSelected() {
        console.log('Edit selected - to be implemented');
    }
    
    deleteSelected() {
        console.log('Delete selected - to be implemented');
    }
    
    editRow(index) {
        console.log('Edit row', index);
    }
    
    deleteRow(index) {
        console.log('Delete row', index);
    }
    
    exportData() {
        const csv = this.convertToCSV(this.data);
        this.downloadCSV(csv, `${this.options.entity}.csv`);
    }
    
    convertToCSV(data) {
        if (data.length === 0) return '';
        
        const headers = Object.keys(data[0]);
        const csvContent = [
            headers.join(','),
            ...data.map(row => headers.map(header => `"${row[header] || ''}"`).join(','))
        ].join('\n');
        
        return csvContent;
    }
    
    downloadCSV(csv, filename) {
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        a.click();
        window.URL.revokeObjectURL(url);
    }
}