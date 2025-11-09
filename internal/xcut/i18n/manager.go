package i18n

import (
	"context"
	"fmt"
	"strings"
)

type I18nManager struct {
	translations map[string]map[string]string
	fallback     string
}

func NewI18nManager(fallback string) *I18nManager {
	im := &I18nManager{
		translations: make(map[string]map[string]string),
		fallback:     fallback,
	}
	im.loadDefaultTranslations()
	return im
}

func (im *I18nManager) GetText(ctx context.Context, key, locale string) string {
	if translations, exists := im.translations[locale]; exists {
		if text, found := translations[key]; found {
			return text
		}
	}

	// Fallback to default locale
	if translations, exists := im.translations[im.fallback]; exists {
		if text, found := translations[key]; found {
			return text
		}
	}

	return key // Return key if no translation found
}

func (im *I18nManager) FormatNumber(value float64, locale string) string {
	switch locale {
	case "en-US":
		return fmt.Sprintf("%.2f", value)
	case "de-DE":
		str := fmt.Sprintf("%.2f", value)
		return strings.Replace(str, ".", ",", 1)
	case "fr-FR":
		str := fmt.Sprintf("%.2f", value)
		return strings.Replace(str, ".", ",", 1)
	default:
		return fmt.Sprintf("%.2f", value)
	}
}

func (im *I18nManager) AddTranslation(locale, key, text string) {
	if im.translations[locale] == nil {
		im.translations[locale] = make(map[string]string)
	}
	im.translations[locale][key] = text
}

func (im *I18nManager) loadDefaultTranslations() {
	// English
	im.AddTranslation("en-US", "country.create.success", "Country created successfully")
	im.AddTranslation("en-US", "country.update.success", "Country updated successfully")
	im.AddTranslation("en-US", "error.validation.required", "This field is required")

	// Spanish
	im.AddTranslation("es-ES", "country.create.success", "País creado exitosamente")
	im.AddTranslation("es-ES", "country.update.success", "País actualizado exitosamente")
	im.AddTranslation("es-ES", "error.validation.required", "Este campo es obligatorio")

	// French
	im.AddTranslation("fr-FR", "country.create.success", "Pays créé avec succès")
	im.AddTranslation("fr-FR", "country.update.success", "Pays mis à jour avec succès")
	im.AddTranslation("fr-FR", "error.validation.required", "Ce champ est obligatoire")
}