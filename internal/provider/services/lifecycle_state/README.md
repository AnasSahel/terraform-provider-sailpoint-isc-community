# Lifecycle State Service - Best Practices Implementation

## 📋 Overview

This service implements comprehensive lifecycle state management for SailPoint ISC using modern Go and Terraform best practices.

## 🏗️ Architecture

### Core Components

```
lifecycle_state/
├── constants.go          # Centralized constants and messages
├── errors.go            # Structured error handling
├── validation.go        # Input validation logic
├── resource.go          # Main resource CRUD operations
├── helpers.go           # Conversion utilities
├── model.go            # Terraform model definitions
└── schema files...     # Schema definitions
```

### Design Principles

1. **Separation of Concerns**: Each file has a single responsibility
2. **Consistent Error Handling**: Standardized error messages and HTTP status handling
3. **Comprehensive Validation**: Input validation before API calls
4. **Code Reusability**: Shared utilities and patterns
5. **Maintainability**: Easy to extend and modify

## 🛠️ Key Features

### 1. Structured Error Handling

```go
// Centralized error handling with context
errorHandler := NewErrorHandler()
diag := errorHandler.HandleAPIError(err, 404, "lifecycle state")
```

**Benefits:**
- Consistent error messages across operations
- HTTP status code specific handling
- Context-aware error details
- Easy to extend for new error types

### 2. Comprehensive Validation

```go
// Multi-layer validation
validator := NewValidator()
diags := validator.ValidateResourceModel(model)
```

**Validation Layers:**
- Field-level validation (name, technical name, etc.)
- Business rule validation (identity states, priorities)
- Import ID validation
- Model-level validation

### 3. Constants Management

```go
// Centralized constants for consistency
const (
    ResourceTypeName = "sailpoint-isc-community_lifecycle_state"
    ErrorLifecycleStateNotFound = "The specified lifecycle state could not be found"
)
```

**Benefits:**
- Single source of truth for constants
- Easy to update error messages
- Consistent resource naming
- Reduced typos and inconsistencies

## 🔄 Resource Lifecycle

### Create Operation
1. **Input Validation** → Validate all required fields
2. **API Call** → Create lifecycle state via SailPoint API
3. **Error Handling** → Process API response with context
4. **State Management** → Update Terraform state

### Read Operation
1. **ID Validation** → Ensure valid resource ID
2. **API Retrieval** → Fetch current state from API
3. **State Sync** → Update Terraform state with API data
4. **Error Handling** → Handle not found and other errors

### Update Operation
1. **Change Detection** → Compare current vs desired state
2. **Validation** → Validate changes
3. **JSON Patch** → Create RFC 6902 patch operations
4. **API Update** → Apply changes via PATCH API
5. **State Refresh** → Update Terraform state

### Delete Operation
1. **Existence Check** → Verify resource exists
2. **API Deletion** → Remove via DELETE API
3. **Error Handling** → Handle deletion errors gracefully

## 🧪 Testing Strategy

### Unit Tests
- **Validation Tests**: Test all validation rules
- **Error Handling Tests**: Test error scenarios
- **Helper Function Tests**: Test conversion utilities

### Integration Tests
- **CRUD Operations**: Test complete resource lifecycle
- **API Integration**: Test with real SailPoint API (when TF_ACC=1)
- **Error Scenarios**: Test various failure modes

## 📚 Usage Examples

### Basic Resource Configuration

```hcl
resource "sailpoint-isc-community_lifecycle_state" "active" {
  identity_profile_id = "2c9180835d2e5168015d32f890ca1581"
  name               = "Active"
  technical_name     = "active"
  description        = "Active lifecycle state for employees"
  identity_state     = "ACTIVE"
  priority           = 10
  
  email_notification_option = jsonencode({
    notifyManagers      = true
    notifyAllAdmins     = false
    notifySpecificUsers = false
    emailAddressList    = []
  })
  
  access_actions = jsonencode([
    {
      action = "ENABLE"
      sourceIds = ["2c9180835d2e5168015d32f890ca1582"]
    }
  ])
}
```

### Import Existing Resource

```bash
terraform import sailpoint-isc-community_lifecycle_state.active "profile_id:state_id"
```

## 🚀 Performance Optimizations

1. **Efficient API Calls**: Minimize unnecessary API requests
2. **Smart State Management**: Only update changed fields
3. **Validation Caching**: Cache validation results where appropriate
4. **Error Context**: Provide meaningful error context without verbose logging

## 🔧 Development Guidelines

### Adding New Validation Rules

1. Add constants to `constants.go`
2. Implement validation in `validation.go`
3. Add tests in `validation_test.go`
4. Update resource operations as needed

### Adding New Error Types

1. Add error constants to `constants.go`
2. Implement handling in `errors.go`
3. Add tests in `errors_test.go`
4. Use in resource operations

### Extending Resource Operations

1. Follow existing patterns in `resource.go`
2. Use validation before API calls
3. Use structured error handling
4. Add comprehensive tests

## 📋 Checklist for New Features

- [ ] Constants added to `constants.go`
- [ ] Validation rules implemented
- [ ] Error handling implemented
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] Examples provided
- [ ] Integration tests pass

## 🏆 Benefits Achieved

1. **🧹 Clean Code**: Separated concerns and reduced duplication
2. **🛡️ Better UX**: Clear, actionable error messages
3. **✅ Reliability**: Comprehensive validation prevents issues
4. **🔄 Consistency**: Standardized patterns across operations
5. **📈 Maintainability**: Easy to extend and modify
6. **🧪 Testability**: Comprehensive test coverage

This implementation represents modern Terraform provider development practices and provides a solid foundation for the lifecycle state service.