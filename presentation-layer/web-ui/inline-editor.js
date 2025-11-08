/**
 * Generic Inline Editor for All Entities
 * Supports countries, regions, languages, timezones, subdivisions, locales
 */

class InlineEntityEditor {
    constructor(entityConfig) {
        this.baseUrl = 'http://localhost:8081/api/v1';
        this.tenantId = 'default-tenant';
        this.entityConfig = entityConfig;
        this.entities = [];
        this.editingRow = null;
        this.drafts = this.loadDrafts();
        this.validationRules = this.initValidationRules();
        this.init();
    }

    init() {
        this.loadEntities();
        this.setupEventListeners();
        this.updateStatusInfo();
        this.updatePageTitle();
    }

    initValidationRules() {
        const rules = {
            countries: {
                country_code: { required: true, pattern: /^[A-Z]{2}$/, message: 'Must be 2 uppercase letters' },
                country_name: { required: true, minLength: 2, message: 'Must be at least 2 characters' },
                iso3_code: { pattern: /^[A-Z]{3}$/, message: 'Must be 3 uppercase letters' },
                continent_code: { enum: ['AF', 'AS', 'EU', 'NA', 'SA', 'OC', 'AN'], message: 'Invalid continent' }
            },
            regions: {
                region_code: { required: true, pattern: /^[A-Z0-9_-]{1,10}$/, message: 'Must be 1-10 characters' },
                region_name: { required: true, minLength: 2, message: 'Must be at least 2 characters' },
                region_type: { enum: ['CONTINENT', 'SUBCONTINENT', 'REGION', 'SUBREGION'], message: 'Invalid type' }
            },
            languages: {
                language_code: { required: true, pattern: /^[a-z]{2,3}$/, message: 'Must be 2-3 lowercase letters' },
                language_name: { required: true, minLength: 2, message: 'Must be at least 2 characters' },
                iso3_code: { pattern: /^[a-z]{3}$/, message: 'Must be 3 lowercase letters' },
                direction: { enum: ['LTR', 'RTL'], message: 'Must be LTR or RTL' }
            }
        };
        return rules[this.entityConfig.name] || {};
    }

    setupEventListeners() {
        document.getElementById('searchInput').addEventListener('input', (e) => this.filterEntities(e.target.value));
        setInterval(() => this.saveDrafts(), 2000);
    }

    updatePageTitle() {
        document.title = `${this.entityConfig.displayName} Management - Inline Editing`;
        const header = document.querySelector('.header h1');
        if (header) {
            header.innerHTML = `${this.entityConfig.icon} ${this.entityConfig.displayName} Management - Inline Editing`;
        }
    }

    async loadEntities() {
        try {
            const response = await fetch(`${this.baseUrl}/${this.entityConfig.endpoint}`, {
                headers: { 'X-Tenant-ID': this.tenantId }
            });
            const data = await response.json();
            this.entities = data[this.entityConfig.dataKey] || [];
            this.renderEntities();
            this.updateStatusInfo();
        } catch (error) {
            console.error(`Failed to load ${this.entityConfig.name}:`, error);
        }
    }

    renderEntities() {
        const tbody = document.getElementById('entitiesTableBody');
        if (this.entities.length === 0) {
            tbody.innerHTML = `<tr><td colspan="${this.entityConfig.columns.length + 1}" style="text-align: center;">No ${this.entityConfig.name} found</td></tr>`;
            return;
        }

        tbody.innerHTML = this.entities.map(entity => this.renderEntityRow(entity)).join('');
    }

    renderEntityRow(entity, isNew = false) {
        const rowClass = isNew ? 'editable-row adding' : 'editable-row';
        const primaryKey = entity[this.entityConfig.primaryKey];
        const deleteDisabled = entity.is_active ? 'disabled title="Deactivate first"' : '';
        
        const cells = this.entityConfig.columns.map(col => {
            let value = entity[col.field] || '';
            let displayValue = value;
            
            if (col.type === 'boolean') {
                displayValue = entity[col.field] ? '✅ Active' : '❌ Inactive';
            }
            
            return `
                <td class="editable-cell" data-field="${col.field}">
                    <span class="display-value ${col.type === 'boolean' ? (entity[col.field] ? 'status-active' : 'status-inactive') : ''}">${displayValue}</span>
                    ${this.renderEditInput(col, value)}
                    <div class="validation-msg" style="display: none;"></div>
                </td>
            `;
        }).join('');

        return `
            <tr class="${rowClass}" data-key="${primaryKey}">
                ${cells}
                <td class="row-actions">
                    <button class="btn btn-edit edit-btn" onclick="editRow('${primaryKey}')" title="Edit">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-save save-btn" onclick="saveRow('${primaryKey}')" style="display: none;" title="Save">
                        <i class="fas fa-save"></i>
                    </button>
                    <button class="btn btn-cancel cancel-btn" onclick="cancelEdit('${primaryKey}')" style="display: none;" title="Cancel">
                        <i class="fas fa-times"></i>
                    </button>
                    <button class="btn btn-delete" onclick="deleteEntity('${primaryKey}')" ${deleteDisabled} title="Delete">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `;
    }

    renderEditInput(column, value) {
        if (column.type === 'select' || column.options) {
            const options = column.options.map(opt => 
                `<option value="${opt.value}" ${value === opt.value ? 'selected' : ''}>${opt.label}</option>`
            ).join('');
            return `<select class="edit-input" style="display: none;">${options}</select>`;
        } else if (column.type === 'boolean') {
            return `
                <select class="edit-input" style="display: none;">
                    <option value="true" ${value ? 'selected' : ''}>Active</option>
                    <option value="false" ${!value ? 'selected' : ''}>Inactive</option>
                </select>
            `;
        } else {
            const inputType = column.type === 'number' ? 'number' : 'text';
            const maxLength = column.maxLength ? `maxlength="${column.maxLength}"` : '';
            return `<input type="${inputType}" class="edit-input" value="${value}" style="display: none;" ${maxLength}>`;
        }
    }

    // Additional methods for editing, validation, saving, etc.
    // ... (rest of the methods from the previous implementation)

    updateStatusInfo() {
        const total = this.entities.length;
        const active = this.entities.filter(e => e.is_active).length;
        const statusEl = document.getElementById('statusInfo');
        if (statusEl) {
            statusEl.textContent = `${total} ${this.entityConfig.name} (${active} active, ${total - active} inactive)`;
        }
    }
}

window.InlineEntityEditor = InlineEntityEditor;