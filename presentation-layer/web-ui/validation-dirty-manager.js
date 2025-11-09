/**
 * Validation and Dirty Record Management System
 * Handles schema-driven validation and dirty tracking with minimal overhead
 */

class ValidationDirtyManager {
    constructor(schema) {
        this.schema = schema;
        this.dirtyRecords = new Map(); // recordId -> {original, current, dirtyFields}
        this.validationCache = new Map(); // recordId -> {fieldName -> isValid}
        this.autoSaveEnabled = false;
        this.autoSaveDelay = 2000;
        this.autoSaveTimer = null;
    }

    // Mark field as dirty only if value actually changed
    markFieldDirty(recordId, fieldName, newValue, originalRecord) {
        const originalValue = originalRecord[fieldName];
        
        // Initialize dirty tracking if needed
        if (!this.dirtyRecords.has(recordId)) {
            this.dirtyRecords.set(recordId, {
                original: { ...originalRecord },
                current: { ...originalRecord },
                dirtyFields: new Set()
            });
        }

        const dirtyData = this.dirtyRecords.get(recordId);
        dirtyData.current[fieldName] = newValue;

        // Only mark as dirty if value changed from original
        if (newValue !== originalValue) {
            dirtyData.dirtyFields.add(fieldName);
        } else {
            dirtyData.dirtyFields.delete(fieldName);
            this.clearFieldValidation(recordId, fieldName);
        }

        // Remove record if no dirty fields
        if (dirtyData.dirtyFields.size === 0) {
            this.dirtyRecords.delete(recordId);
            this.validationCache.delete(recordId);
        }

        return dirtyData.dirtyFields.has(fieldName);
    }

    // Validate only dirty fields
    validateField(recordId, fieldName, value) {
        const fieldSchema = this.schema[fieldName];
        if (!fieldSchema) return { isValid: true };

        // Only validate if field is dirty
        const dirtyData = this.dirtyRecords.get(recordId);
        if (!dirtyData || !dirtyData.dirtyFields.has(fieldName)) {
            return { isValid: true };
        }

        const result = this.performValidation(fieldSchema, value);
        
        // Cache validation result
        if (!this.validationCache.has(recordId)) {
            this.validationCache.set(recordId, {});
        }
        this.validationCache.get(recordId)[fieldName] = result.isValid;

        return result;
    }

    performValidation(fieldSchema, value) {
        const errors = [];
        
        // Required check
        if (fieldSchema.required && this.isEmpty(value)) {
            return { isValid: false, error: `${fieldSchema.label} is required` };
        }

        // Skip other validations if empty and not required
        if (!fieldSchema.required && this.isEmpty(value)) {
            return { isValid: true };
        }

        // Type-specific validation
        switch (fieldSchema.type) {
            case 'CHAR(2)':
            case 'CHAR(3)':
                if (fieldSchema.pattern && !fieldSchema.pattern.test(value)) {
                    errors.push(`${fieldSchema.label} format invalid`);
                }
                break;
                
            case 'VARCHAR(100)':
            case 'VARCHAR(200)':
                if (fieldSchema.maxLength && value.length > fieldSchema.maxLength) {
                    errors.push(`Max ${fieldSchema.maxLength} characters`);
                }
                break;
                
            case 'SMALLINT':
            case 'INTEGER':
                const num = parseInt(value);
                if (isNaN(num)) {
                    errors.push('Must be a number');
                } else {
                    if (fieldSchema.min !== undefined && num < fieldSchema.min) {
                        errors.push(`Min value: ${fieldSchema.min}`);
                    }
                    if (fieldSchema.max !== undefined && num > fieldSchema.max) {
                        errors.push(`Max value: ${fieldSchema.max}`);
                    }
                }
                break;
                
            case 'ENUM':
                const validValues = fieldSchema.enum.map(e => e.value);
                if (!validValues.includes(value)) {
                    errors.push('Invalid selection');
                }
                break;
        }

        return {
            isValid: errors.length === 0,
            error: errors[0] || null
        };
    }

    isEmpty(value) {
        return value === null || value === undefined || value === '';
    }

    clearFieldValidation(recordId, fieldName) {
        const cache = this.validationCache.get(recordId);
        if (cache) {
            delete cache[fieldName];
        }
    }

    isRecordValid(recordId) {
        const cache = this.validationCache.get(recordId);
        if (!cache) return true;
        return Object.values(cache).every(isValid => isValid !== false);
    }

    getDirtyRecords() {
        return Array.from(this.dirtyRecords.keys());
    }

    getDirtyFields(recordId) {
        const dirtyData = this.dirtyRecords.get(recordId);
        return dirtyData ? Array.from(dirtyData.dirtyFields) : [];
    }

    getValidDirtyRecords() {
        return this.getDirtyRecords().filter(id => this.isRecordValid(id));
    }

    revertRecord(recordId) {
        this.dirtyRecords.delete(recordId);
        this.validationCache.delete(recordId);
    }

    // Auto-save functionality
    enableAutoSave(saveCallback, delay = 2000) {
        this.autoSaveEnabled = true;
        this.autoSaveDelay = delay;
        this.saveCallback = saveCallback;
    }

    triggerAutoSave() {
        if (!this.autoSaveEnabled) return;
        
        clearTimeout(this.autoSaveTimer);
        this.autoSaveTimer = setTimeout(() => {
            const validRecords = this.getValidDirtyRecords();
            if (validRecords.length > 0 && this.saveCallback) {
                this.saveCallback(validRecords);
            }
        }, this.autoSaveDelay);
    }

    // Get changes for a record
    getRecordChanges(recordId) {
        const dirtyData = this.dirtyRecords.get(recordId);
        if (!dirtyData) return null;

        const changes = {};
        dirtyData.dirtyFields.forEach(fieldName => {
            changes[fieldName] = {
                old: dirtyData.original[fieldName],
                new: dirtyData.current[fieldName]
            };
        });
        return changes;
    }

    // Get summary of all changes
    getChangesSummary() {
        const summary = {
            totalRecords: this.dirtyRecords.size,
            totalFields: 0,
            validRecords: 0,
            invalidRecords: 0
        };

        this.dirtyRecords.forEach((data, recordId) => {
            summary.totalFields += data.dirtyFields.size;
            if (this.isRecordValid(recordId)) {
                summary.validRecords++;
            } else {
                summary.invalidRecords++;
            }
        });

        return summary;
    }
}

// Export for use
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ValidationDirtyManager;
} else {
    window.ValidationDirtyManager = ValidationDirtyManager;
}