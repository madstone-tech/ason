# UOW-013: Add Comprehensive Tests for XDG Package (0% Coverage)

## Overview
**Phase**: 5 - Testing
**Priority**: Medium
**Estimated Effort**: 4-6 hours
**Dependencies**: Phase 2 completion

## Problem Description
The XDG package currently has 0% test coverage, which is problematic because:
- XDG directory functions are critical for cross-platform compatibility
- No validation of directory creation and permission handling
- Missing tests for edge cases and error conditions
- No verification of XDG Base Directory Specification compliance
- Potential issues on different operating systems go undetected

## Acceptance Criteria
- [ ] Achieve >95% test coverage for XDG package
- [ ] Test all XDG Base Directory Specification functions
- [ ] Cross-platform compatibility tests (Windows, macOS, Linux)
- [ ] Edge case testing (permissions, missing directories, invalid paths)
- [ ] Environment variable override testing
- [ ] Error condition handling verification
- [ ] Mock filesystem testing for isolation

## Technical Approach

### Implementation Strategy
1. Analyze current XDG package implementation
2. Create comprehensive unit tests for all functions
3. Add cross-platform compatibility tests
4. Implement filesystem mocking for test isolation
5. Add edge case and error condition tests

### Current XDG Package Analysis

The XDG package contains three main functions that need comprehensive testing:

1. **ConfigHome()** - Returns XDG_CONFIG_HOME/ason directory
2. **DataHome()** - Returns XDG_DATA_HOME/ason directory
3. **CacheHome()** - Returns XDG_CACHE_HOME/ason directory

Each function follows the same pattern:
- Check for environment variable (XDG_CONFIG_HOME, XDG_DATA_HOME, XDG_CACHE_HOME)
- Fall back to default locations if env var not set (~/.config, ~/.local/share, ~/.cache)
- Append "ason" subdirectory to the path
- Return error if os.UserHomeDir() fails

<function_calls>
<invoke name="Read">
<parameter name="file_path">/Users/andhi/code/mdstn/ason/internal/xdg/xdg.go