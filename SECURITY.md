# Security Policy

## Supported Versions

We actively support the following versions of Ason with security updates:

| Version | Supported |
| ------- | --------- |
| 1.x.x   | ✅ Yes    |
| < 1.0   | ❌ No     |

## Reporting a Vulnerability

### Critical Security Vulnerabilities

**For critical security vulnerabilities that could compromise user data or systems, please DO NOT create a public GitHub issue.**

Instead, please report them privately by emailing: **<security@madstone.io>**

#### What constitutes a critical vulnerability?

- Remote code execution
- Arbitrary file write/read outside intended directories
- Authentication/authorization bypasses
- Injection vulnerabilities (command injection, path traversal, etc.)
- Exposure of sensitive data
- Denial of service with minimal effort

### Non-Critical Security Issues

For security improvements, hardening suggestions, or minor security concerns that don't pose immediate risk, you can:

- Create a public issue using the "Security Report" template
- Email <security@madstone.io> if you're unsure

## Response Process

### For Critical Vulnerabilities

1. **Acknowledgment**: We'll acknowledge receipt within 48 hours
2. **Initial Assessment**: We'll provide an initial assessment within 5 business days
3. **Investigation**: We'll investigate and develop a fix
4. **Coordination**: We'll coordinate disclosure timeline with you
5. **Release**: We'll release a security update and advisory
6. **Credit**: We'll provide appropriate credit if desired

### Timeline Expectations

- **Acknowledgment**: Within 48 hours
- **Initial response**: Within 5 business days
- **Fix development**: Depends on severity and complexity
- **Security release**: As soon as possible after fix is ready

## Security Measures

### Current Security Practices

- **Input validation**: All user inputs are validated
- **Path sanitization**: File paths are sanitized to prevent traversal attacks
- **Dependency scanning**: Dependencies are regularly scanned for vulnerabilities
- **Code review**: All code changes undergo security-conscious review
- **Automated testing**: Security-focused tests are included in our test suite

### Planned Security Improvements

See our [security roadmap](./roadmap/phase1-security/) for upcoming security enhancements:

- Enhanced path traversal protection
- Improved input validation
- Security-focused refactoring

## Security Best Practices for Users

### Template Security

When creating or using templates:

- **Validate template sources**: Only use templates from trusted sources
- **Review template content**: Check templates for suspicious patterns
- **Limit template variables**: Don't pass sensitive data as template variables
- **Sandbox execution**: Run Ason in isolated environments when processing untrusted templates

### Installation Security

- **Verify checksums**: Always verify binary checksums when downloading releases
- **Use package managers**: Prefer official package managers when available
- **Keep updated**: Regularly update to the latest version
- **Monitor advisories**: Subscribe to security advisories

### Configuration Security

- **File permissions**: Set appropriate permissions on configuration files
- **Secrets management**: Don't store secrets in configuration files
- **Access control**: Limit access to Ason data directories
- **Network isolation**: Use network restrictions in sensitive environments

## Known Security Considerations

### Template Processing Risks

- **Untrusted templates**: Processing untrusted templates may pose security risks
- **File permissions**: Generated files inherit template file permissions
- **Path handling**: Template paths are validated but defense in depth is recommended
- **Resource usage**: Large templates may consume significant system resources

### Mitigation Strategies

- Use Ason in containerized or sandboxed environments for untrusted content
- Regularly audit generated output for unexpected files or permissions
- Implement resource limits in production environments
- Monitor file system access when processing templates

## Security Contact

- **Email**: <security@madstone.io>
- **PGP Key**: Available on request
- **Response Time**: We aim to respond to security reports within 48 hours

## Security Updates

Security updates are distributed through:

- **GitHub Releases**: All security updates are clearly marked
- **Security Advisories**: Published on GitHub Security Advisories
- **Changelog**: Security fixes are documented in CHANGELOG.md
- **Email notifications**: For critical vulnerabilities (if contact provided)

## Disclosure Policy

We follow responsible disclosure principles:

- **Coordination**: We coordinate with reporters on disclosure timelines
- **Credit**: We provide appropriate credit to security researchers
- **Transparency**: We publish security advisories after fixes are available
- **Timeline**: We aim for disclosure within 90 days of initial report

## Bug Bounty

Currently, we don't have a formal bug bounty program, but we:

- Recognize and credit security researchers
- Prioritize security fixes
- Welcome responsible disclosure
- Consider small tokens of appreciation for significant findings

## Legal Protection

We support security research conducted in good faith and will not pursue legal action against researchers who:

- Follow responsible disclosure practices
- Don't access or modify user data without permission
- Don't perform attacks that could harm our users or systems
- Report findings privately to <security@madstone.io>

---

Thank you for helping keep Ason and our users safe!

