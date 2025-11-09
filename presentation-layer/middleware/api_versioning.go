package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIVersionHandler struct {
	supportedVersions map[string]bool
	defaultVersion    string
}

func NewAPIVersionHandler(defaultVersion string, supportedVersions []string) *APIVersionHandler {
	versions := make(map[string]bool)
	for _, v := range supportedVersions {
		versions[v] = true
	}
	return &APIVersionHandler{
		supportedVersions: versions,
		defaultVersion:    defaultVersion,
	}
}

func (avh *APIVersionHandler) VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := avh.extractVersion(c)
		
		if !avh.supportedVersions[version] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "API_VERSION_001",
					"message": "Unsupported API version",
					"supported_versions": avh.getSupportedVersions(),
				},
			})
			c.Abort()
			return
		}

		c.Set("api_version", version)
		c.Header("API-Version", version)
		c.Next()
	}
}

func (avh *APIVersionHandler) extractVersion(c *gin.Context) string {
	// Check Accept header first (e.g., application/vnd.api+json;version=1)
	accept := c.GetHeader("Accept")
	if strings.Contains(accept, "version=") {
		parts := strings.Split(accept, "version=")
		if len(parts) > 1 {
			version := strings.Split(parts[1], ";")[0]
			version = strings.Split(version, ",")[0]
			return strings.TrimSpace(version)
		}
	}

	// Check API-Version header
	if version := c.GetHeader("API-Version"); version != "" {
		return version
	}

	// Check query parameter
	if version := c.Query("version"); version != "" {
		return version
	}

	// Extract from URL path (e.g., /api/v1/countries)
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/api/v") {
		parts := strings.Split(path, "/")
		if len(parts) > 2 {
			return strings.TrimPrefix(parts[2], "v")
		}
	}

	return avh.defaultVersion
}

func (avh *APIVersionHandler) getSupportedVersions() []string {
	var versions []string
	for v := range avh.supportedVersions {
		versions = append(versions, v)
	}
	return versions
}

func (avh *APIVersionHandler) HandleVersion(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiVersion := c.GetString("api_version")
		
		switch apiVersion {
		case "1":
			avh.handleV1(c)
		case "2":
			avh.handleV2(c)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "API_VERSION_002",
					"message": "Version handler not implemented",
				},
			})
		}
	}
}

func (avh *APIVersionHandler) handleV1(c *gin.Context) {
	// V1 specific handling
	c.Set("response_format", "v1")
	c.Next()
}

func (avh *APIVersionHandler) handleV2(c *gin.Context) {
	// V2 specific handling with enhanced features
	c.Set("response_format", "v2")
	c.Set("include_metadata", true)
	c.Next()
}