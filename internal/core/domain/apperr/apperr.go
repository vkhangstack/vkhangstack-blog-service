package apperr

import "github.com/vkhangstack/hexagonal-architecture/internal/core/domain"

func New(code domain.ErrorCode) string {
	switch code {
	case domain.ErrorCodePayloadBadRequest:
		return "Invalid request payload"
	case domain.ErrorCodeBlogCategoryNotFound:
		return "Blog category not found"
	case domain.ErrorCodeBlogPostNotFound:
		return "Blog post not found"
	case domain.ErrorCodeTokenNotFound:
		return "Token not found"
	case domain.ErrorCodeUserNotFound:
		return "User not found"
	case domain.ErrorCodeUserInactive:
		return "User is inactive"
	case domain.ErrorCodeInsufficientPermission:
		return "Insufficient permission"
	case domain.ErrorCodeEmailExists:
		return "Email already exists"
	case domain.ErrorCodeEmailNotExists:
		return "Email does not exist"
	case domain.ErrorCodeFullName:
		return "Full name is invalid"
	case domain.ErrorCodeRole:
		return "Role is invalid"
	case domain.ErrorCodeTenantID:
		return "Tenant ID is invalid"
	case domain.ErrorCodePassword:
		return "Password is invalid"
	case domain.ErrorCodeUserStatusNotFound:
		return "User status not found"
	case domain.ErrorCodeProductNotFound:
		return "Product not found"
	case domain.ErrorCodeInventoryNotFound:
		return "Inventory not found"
	case domain.ErrorCodeOrderNotFound:
		return "Order not found"
	case domain.ErrorCodeUnAuthorization:
		return "Unauthorized"
	case domain.ErrorCodeForbidden:
		return "Forbidden"
	case domain.ErrorCodeInternalServerError:
		return "Internal server error"
	case domain.ErrorCodeTooManyRequests:
		return "Too many requests"
	default:
		return "Unknown error"
	}
}
