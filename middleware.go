package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ctxUserKey = "user"

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		uidFloat, ok := claims["userId"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			return
		}

		uid := uint(uidFloat)
		u, err := GetUserByID(uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set(ctxUserKey, u)
		c.Next()
	}
}

func GetCtxUser(c *gin.Context) (*User, bool) {
	v, ok := c.Get(ctxUserKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*User)
	return u, ok
}

func RBAC(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := GetCtxUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}

		// Check if user has the required permission through their roles
		var count int64
		DB.Raw(
			`SELECT COUNT(1) FROM role_permissions rp
			 JOIN user_roles ur ON ur.role_id = rp.role_id
			 WHERE ur.user_id = ? AND rp.permission_id = (
			   SELECT id FROM permissions WHERE name = ?
			 )`,
			u.ID, permissionName,
		).Scan(&count)

		if count > 0 {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

func ReBACResource(relationType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := GetCtxUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}

		id := c.Param("id")
		var rec PatientRecord
		if err := DB.First(&rec, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		// Allow owner to access their own records
		if rec.OwnerID == u.ID {
			c.Next()
			return
		}

		// Check for direct relationship: user -> record
		var directRel Relationship
		if err := DB.Where("subject_id = ? AND object_id = ? AND type = ?", u.ID, rec.ID, relationType).First(&directRel).Error; err == nil {
			c.Next()
			return
		}

		// Check for indirect relationship: user -> owner (e.g., doctor assigned to patient)
		var indirectRel Relationship
		if err := DB.Where("subject_id = ? AND object_id = ? AND type = ?", u.ID, rec.OwnerID, relationType).First(&indirectRel).Error; err == nil {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no relationship found"})
	}
}
