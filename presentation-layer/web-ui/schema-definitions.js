/**
 * Schema Definition System for LASANI Platform
 * Provides database schema metadata for frontend validation and UI generation
 */

class SchemaDefinitionSystem {
    constructor() {
        this.schemas = new Map();
        this.initializeSchemas();
    }

    initializeSchemas() {
        // Countries Schema
        this.registerSchema('countries', {
            tableName: 'domain_reference_master_geopolitical.countries',
            primaryKey: 'country_code',
            displayName: 'Countries',
            fields: {
                country_id: {
                    type: 'UUID',
                    primaryKey: true,
                    generated: true,
                    label: 'ID',
                    editable: false,
                    visible: false
                },
                country_code: {
                    type: 'CHAR(2)',
                    required: true,
                    unique: true,
                    pattern: /^[A-Z]{2}$/,
                    maxLength: 2,
                    label: 'Code',
                    editable: true,
                    businessKey: true,
                    collation: 'C',
                    description: 'ISO 3166-1 alpha-2 country code'
                },
                country_name: {
                    type: 'VARCHAR(100)',
                    required: true,
                    minLength: 2,
                    maxLength: 100,
                    label: 'Country Name',
                    editable: true,
                    searchable: true,
                    description: 'Official country name'
                },
                iso3_code: {
                    type: 'CHAR(3)',
                    unique: true,
                    pattern: /^[A-Z]{3}$/,
                    maxLength: 3,
                    label: 'ISO3 Code',
                    editable: true,
                    collation: 'C',
                    description: 'ISO 3166-1 alpha-3 country code'
                },
                numeric_code: {
                    type: 'SMALLINT',
                    unique: true,
                    min: 1,
                    max: 999,
                    label: 'Numeric Code',
                    editable: true,
                    description: 'ISO 3166-1 numeric country code'
                },
                official_name: {
                    type: 'VARCHAR(200)',
                    maxLength: 200,
                    label: 'Official Name',
                    editable: true,
                    description: 'Full official country name'
                },
                capital_city: {
                    type: 'VARCHAR(100)',
                    maxLength: 100,
                    label: 'Capital City',
                    editable: true,
                    description: 'Capital city name'
                },
                continent_code: {
                    type: 'ENUM',
                    enumName: 'continent_enum',
                    enum: [
                        { value: 'AF', label: 'Africa', description: 'African continent' },
                        { value: 'AS', label: 'Asia', description: 'Asian continent' },
                        { value: 'EU', label: 'Europe', description: 'European continent' },
                        { value: 'NA', label: 'North America', description: 'North American continent' },
                        { value: 'SA', label: 'South America', description: 'South American continent' },
                        { value: 'OC', label: 'Oceania', description: 'Oceania continent' },
                        { value: 'AN', label: 'Antarctica', description: 'Antarctic continent' }
                    ],
                    label: 'Continent',
                    editable: true,
                    description: 'Continental classification'
                },
                region_id: {
                    type: 'UUID',
                    foreignKey: {
                        table: 'domain_reference_master_geopolitical.regions',
                        field: 'region_id',
                        displayField: 'region_name'
                    },
                    label: 'Region',
                    editable: true,
                    description: 'Geographic region reference'
                },
                primary_language_id: {
                    type: 'UUID',
                    foreignKey: {
                        table: 'domain_reference_master_geopolitical.languages',
                        field: 'language_id',
                        displayField: 'language_name'
                    },
                    label: 'Primary Language',
                    editable: true,
                    description: 'Primary language reference'
                },
                currency_id: {
                    type: 'UUID',
                    foreignKey: {
                        table: 'ref.currencies',
                        field: 'currency_id',
                        displayField: 'currency_code'
                    },
                    label: 'Currency',
                    editable: true,
                    description: 'Primary currency reference'
                },
                phone_prefix: {
                    type: 'VARCHAR(10)',
                    pattern: /^\\+\\d{1,4}$/,
                    maxLength: 10,
                    label: 'Phone Prefix',
                    editable: true,
                    collation: 'C',
                    description: 'International dialing prefix'
                },
                is_active: {
                    type: 'BOOLEAN',
                    default: true,
                    required: true,
                    enum: [
                        { value: true, label: 'Active', description: 'Country is active' },
                        { value: false, label: 'Inactive', description: 'Country is inactive' }
                    ],
                    label: 'Active Status',
                    editable: true,
                    description: 'Whether the country is active'
                },
                // LASANI Audit Fields (read-only)
                created_at: {
                    type: 'TIMESTAMPTZ',
                    label: 'Created At',
                    editable: false,
                    audit: true,
                    description: 'Record creation timestamp'
                },
                created_by: {
                    type: 'UUID',
                    label: 'Created By',
                    editable: false,
                    audit: true,
                    description: 'User who created the record'
                },
                updated_at: {
                    type: 'TIMESTAMPTZ',
                    label: 'Updated At',
                    editable: false,
                    audit: true,
                    description: 'Last update timestamp'
                },
                updated_by: {
                    type: 'UUID',
                    label: 'Updated By',
                    editable: false,
                    audit: true,
                    description: 'User who last updated the record'
                },
                version: {
                    type: 'INTEGER',
                    default: 1,
                    label: 'Version',
                    editable: false,
                    audit: true,
                    description: 'Optimistic locking version'
                }
            },
            indexes: [
                { name: 'idx_countries_active', fields: ['is_active'], where: 'is_deleted = FALSE' },
                { name: 'idx_countries_continent', fields: ['continent_code'] },
                { name: 'idx_countries_name', fields: ['country_name'] }
            ],
            constraints: [
                { name: 'chk_countries_soft_delete', check: '(is_deleted = TRUE AND deleted_at IS NOT NULL) OR (is_deleted = FALSE AND deleted_at IS NULL)' }
            ]
        });

        // Currencies Schema (for reference)
        this.registerSchema('currencies', {
            tableName: 'ref.currencies',
            primaryKey: 'currency_id',
            displayName: 'Currencies',
            fields: {
                currency_id: {
                    type: 'UUID',
                    primaryKey: true,
                    generated: true,
                    label: 'ID',
                    editable: false
                },
                currency_code: {
                    type: 'CHAR(3)',
                    required: true,
                    unique: true,
                    pattern: /^[A-Z]{3}$/,
                    maxLength: 3,
                    label: 'Currency Code',
                    editable: true,
                    businessKey: true,
                    description: 'ISO 4217 currency code'
                },
                currency_name: {
                    type: 'VARCHAR(100)',
                    required: true,
                    maxLength: 100,
                    label: 'Currency Name',
                    editable: true,
                    description: 'Full currency name'
                },
                numeric_code: {
                    type: 'SMALLINT',
                    unique: true,
                    min: 1,
                    max: 999,
                    label: 'Numeric Code',
                    editable: true,
                    description: 'ISO 4217 numeric code'
                },
                symbol: {
                    type: 'VARCHAR(10)',
                    maxLength: 10,
                    label: 'Symbol',
                    editable: true,
                    description: 'Currency symbol'
                },
                minor_unit: {
                    type: 'SMALLINT',
                    default: 2,
                    min: 0,
                    max: 4,
                    label: 'Minor Unit',
                    editable: true,
                    description: 'Number of decimal places'
                },
                is_active: {
                    type: 'BOOLEAN',
                    default: true,
                    required: true,
                    enum: [
                        { value: true, label: 'Active' },
                        { value: false, label: 'Inactive' }
                    ],
                    label: 'Active',
                    editable: true
                }
            }
        });
    }

    registerSchema(entityName, schema) {
        this.schemas.set(entityName, schema);
    }

    getSchema(entityName) {
        return this.schemas.get(entityName);
    }

    getFieldSchema(entityName, fieldName) {
        const schema = this.getSchema(entityName);
        return schema?.fields?.[fieldName];
    }

    getEditableFields(entityName) {
        const schema = this.getSchema(entityName);
        if (!schema) return [];
        
        return Object.entries(schema.fields)
            .filter(([_, field]) => field.editable)
            .map(([name, field]) => ({ name, ...field }));
    }

    getRequiredFields(entityName) {
        const schema = this.getSchema(entityName);
        if (!schema) return [];
        
        return Object.entries(schema.fields)
            .filter(([_, field]) => field.required)
            .map(([name, field]) => ({ name, ...field }));
    }

    validateField(entityName, fieldName, value) {
        const field = this.getFieldSchema(entityName, fieldName);
        if (!field) return { isValid: true };

        const errors = [];
        
        // Required validation
        if (field.required && (value === null || value === undefined || value === '')) {
            errors.push(`${field.label} is required`);
        }
        
        // Skip other validations if empty and not required
        if (!field.required && (value === null || value === undefined || value === '')) {
            return { isValid: true };
        }

        const stringValue = value?.toString() || '';
        
        // Type-specific validations
        switch (field.type) {
            case 'CHAR(2)':
            case 'CHAR(3)':
                if (field.pattern && !field.pattern.test(stringValue)) {
                    errors.push(`${field.label} must match pattern: ${field.pattern.source}`);
                }
                if (field.maxLength && stringValue.length !== field.maxLength) {
                    errors.push(`${field.label} must be exactly ${field.maxLength} characters`);
                }
                break;
                
            case 'VARCHAR(100)':
            case 'VARCHAR(200)':
            case 'VARCHAR(10)':
                if (field.minLength && stringValue.length < field.minLength) {
                    errors.push(`${field.label} must be at least ${field.minLength} characters`);
                }
                if (field.maxLength && stringValue.length > field.maxLength) {
                    errors.push(`${field.label} must be at most ${field.maxLength} characters`);
                }
                if (field.pattern && !field.pattern.test(stringValue)) {
                    errors.push(`${field.label} format is invalid`);
                }
                break;
                
            case 'SMALLINT':
            case 'INTEGER':
                const numValue = parseInt(value);
                if (isNaN(numValue)) {
                    errors.push(`${field.label} must be a valid number`);
                } else {
                    if (field.min !== undefined && numValue < field.min) {
                        errors.push(`${field.label} must be at least ${field.min}`);
                    }
                    if (field.max !== undefined && numValue > field.max) {
                        errors.push(`${field.label} must be at most ${field.max}`);
                    }
                }
                break;
                
            case 'BOOLEAN':
                if (typeof value !== 'boolean' && value !== 'true' && value !== 'false') {
                    errors.push(`${field.label} must be true or false`);
                }
                break;
                
            case 'ENUM':
                const validValues = field.enum.map(e => e.value);
                if (!validValues.includes(value)) {
                    errors.push(`${field.label} must be one of: ${validValues.join(', ')}`);
                }
                break;
                
            case 'UUID':
                const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
                if (stringValue && !uuidPattern.test(stringValue)) {
                    errors.push(`${field.label} must be a valid UUID`);
                }
                break;
        }

        return {
            isValid: errors.length === 0,
            errors: errors,
            field: field
        };
    }

    generateFormField(entityName, fieldName, value = '', options = {}) {
        const field = this.getFieldSchema(entityName, fieldName);
        if (!field || !field.editable) return null;

        const fieldId = options.fieldId || `${entityName}_${fieldName}`;
        const onChange = options.onChange || `handleFieldChange('${fieldId}', this.value)`;
        const onBlur = options.onBlur || `validateField('${fieldId}')`;
        
        let inputHtml = '';
        
        if (field.type === 'ENUM') {
            const optionsHtml = field.enum.map(opt => 
                `<option value="${opt.value}" ${value == opt.value ? 'selected' : ''} title="${opt.description || ''}">${opt.label}</option>`
            ).join('');
            
            inputHtml = `
                <select id="${fieldId}" name="${fieldName}" onchange="${onChange}" onblur="${onBlur}" ${field.required ? 'required' : ''}>
                    <option value="">Select ${field.label}...</option>
                    ${optionsHtml}
                </select>
            `;
        } else if (field.type === 'BOOLEAN') {
            const optionsHtml = field.enum.map(opt => 
                `<option value="${opt.value}" ${value == opt.value ? 'selected' : ''}>${opt.label}</option>`
            ).join('');
            
            inputHtml = `
                <select id="${fieldId}" name="${fieldName}" onchange="${onChange}" onblur="${onBlur}" ${field.required ? 'required' : ''}>
                    ${optionsHtml}
                </select>
            `;
        } else {
            const inputType = field.type.includes('INT') ? 'number' : 'text';
            const maxLength = field.maxLength ? `maxlength="${field.maxLength}"` : '';
            const min = field.min !== undefined ? `min="${field.min}"` : '';
            const max = field.max !== undefined ? `max="${field.max}"` : '';
            const pattern = field.pattern ? `pattern="${field.pattern.source}"` : '';
            
            inputHtml = `
                <input type="${inputType}" 
                       id="${fieldId}" 
                       name="${fieldName}" 
                       value="${value}" 
                       onchange="${onChange}" 
                       onblur="${onBlur}"
                       ${field.required ? 'required' : ''}
                       ${maxLength}
                       ${min}
                       ${max}
                       ${pattern}
                       placeholder="${field.description || field.label}">
            `;
        }

        return `
            <div class="form-group">
                <label for="${fieldId}">${field.label}${field.required ? ' *' : ''}</label>
                ${inputHtml}
                <div class="validation-message" id="${fieldId}_validation"></div>
                ${field.description ? `<small class="field-description">${field.description}</small>` : ''}
            </div>
        `;
    }

    getTableColumns(entityName, includeAudit = false) {
        const schema = this.getSchema(entityName);
        if (!schema) return [];
        
        return Object.entries(schema.fields)
            .filter(([_, field]) => {
                if (!includeAudit && field.audit) return false;
                return field.visible !== false;
            })
            .map(([name, field]) => ({
                name,
                label: field.label,
                type: field.type,
                editable: field.editable,
                sortable: true,
                ...field
            }));
    }
}

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = SchemaDefinitionSystem;
} else {
    window.SchemaDefinitionSystem = SchemaDefinitionSystem;
}