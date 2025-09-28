# Pull Request

## Summary

<!-- Provide a brief summary of the changes in this PR -->

## Related Issues

<!-- Link to any related issues -->
Closes #<!-- issue number -->
Relates to #<!-- issue number -->

## Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] 🐛 Bug fix (non-breaking change that fixes an issue)
- [ ] ✨ New feature (non-breaking change that adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to change)
- [ ] 📚 Documentation update
- [ ] 🔧 Code refactoring (no functional changes)
- [ ] ⚡ Performance improvement
- [ ] 🔒 Security improvement
- [ ] 🧪 Test coverage improvement
- [ ] 🔨 Build/CI improvements
- [ ] 🏗️ Infrastructure changes

## Changes Made

<!-- Describe the changes made in this PR -->

### Modified Files
<!-- List the main files that were changed -->
- `path/to/file.go` - Description of changes
- `path/to/test.go` - Description of changes

### New Files
<!-- List any new files added -->
- `path/to/new/file.go` - Purpose of the new file

## Testing

<!-- Describe how you tested these changes -->

### Test Coverage
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] All existing tests pass

### Test Results
```bash
# Show test command output
go test ./...
```

### Manual Testing Steps
<!-- If manual testing was required, list the steps -->
1. Step one
2. Step two
3. Step three

## Performance Impact

<!-- If applicable, describe any performance implications -->
- [ ] No performance impact
- [ ] Performance improvement (describe)
- [ ] Performance regression (justify why acceptable)
- [ ] Performance testing completed

### Benchmarks
<!-- Include benchmark results if applicable -->
```bash
# Before changes
BenchmarkFunction-8    1000000    1234 ns/op

# After changes
BenchmarkFunction-8    1000000    1000 ns/op
```

## Security Considerations

<!-- Address any security implications -->
- [ ] No security implications
- [ ] Security improvement (describe)
- [ ] Security review completed
- [ ] No sensitive data exposed

## Breaking Changes

<!-- If this includes breaking changes, describe them -->
- [ ] No breaking changes
- [ ] Breaking changes documented below

### Migration Guide
<!-- If breaking changes exist, provide migration guidance -->

## Documentation

<!-- Documentation changes -->
- [ ] Code comments updated
- [ ] README updated
- [ ] API documentation updated
- [ ] CHANGELOG updated
- [ ] No documentation changes needed

## Checklist

<!-- Ensure all items are completed before requesting review -->

### Code Quality
- [ ] Code follows project conventions and style guidelines
- [ ] Self-review completed
- [ ] Code is well-commented, particularly complex areas
- [ ] No TODO comments left in production code
- [ ] Error handling is appropriate
- [ ] Logging is appropriate (not too verbose/quiet)

### Testing
- [ ] All tests pass locally
- [ ] New tests added for new functionality
- [ ] Edge cases considered and tested
- [ ] No test files were accidentally committed with `.only` or similar

### Security
- [ ] No secrets or sensitive data in code
- [ ] Input validation added where needed
- [ ] Security best practices followed
- [ ] Dependencies are up to date and secure

### Performance
- [ ] No obvious performance regressions
- [ ] Memory usage considered
- [ ] Database queries optimized (if applicable)
- [ ] Caching strategy appropriate (if applicable)

### Documentation
- [ ] Public APIs are documented
- [ ] Complex logic is explained
- [ ] README updated if needed
- [ ] Examples provided where helpful

### Git History
- [ ] Commits are atomic and well-described
- [ ] No merge commits (rebased if needed)
- [ ] Branch is up to date with target branch
- [ ] Commit messages follow conventional format

## Screenshots/Recordings

<!-- If UI changes, include screenshots or recordings -->

## Additional Notes

<!-- Any additional information for reviewers -->

## Review Guidance

<!-- Help reviewers focus on important areas -->
- **Focus Areas**: What should reviewers pay special attention to?
- **Risk Areas**: What parts of the change are most likely to have issues?
- **Testing Notes**: Any special testing considerations?

---

**For Reviewers:**
- [ ] Code review completed
- [ ] Tests reviewed and adequate
- [ ] Documentation reviewed
- [ ] Security considerations reviewed
- [ ] Performance impact assessed
- [ ] Breaking changes documented and justified