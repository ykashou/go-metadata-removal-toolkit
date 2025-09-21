# Security Operations

This directory contains security-related documentation, policies, and scanning configurations for the go-metadata-removal-toolkit project.

## Directory Structure

```
security/
├── compliance/     # Compliance documentation and checklists
├── incidents/      # Security incident reports and responses
├── policies/       # Security policies and procedures
└── scans/         # Security scan results and configurations
```

## Security Scanning

### Dependency Scanning
```bash
# Check for vulnerable dependencies
go list -json -m all | nancy sleuth

# Alternative using govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Static Analysis
```bash
# Run gosec for security issues
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -fmt json -out ops/security/scans/gosec-report.json ./...

# Run staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

### Container Scanning
```bash
# Scan container with Trivy
podman run --rm -v /var/run/podman/podman.sock:/var/run/docker.sock \
    aquasec/trivy image metadata-remover:latest

# Scan with Grype
grype metadata-remover:latest
```

## Security Policies

1. **Code Review**: All code must be reviewed before merging
2. **Dependencies**: Regular updates and vulnerability scanning
3. **Containers**: Distroless base images in production
4. **Secrets**: No hardcoded credentials or sensitive data
5. **Least Privilege**: Run with minimal permissions

## Incident Response

In case of a security incident:
1. Document the incident in `incidents/` directory
2. Assess impact and severity
3. Implement immediate mitigation
4. Create permanent fix
5. Update security measures

## Compliance Checklist

- [ ] OWASP Top 10 reviewed
- [ ] CWE/SANS Top 25 reviewed
- [ ] Input validation implemented
- [ ] Error handling secure
- [ ] Logging appropriate
- [ ] Authentication/Authorization proper
- [ ] Data encryption at rest/transit
- [ ] Regular security updates

## Security Contacts

Report security vulnerabilities to:
- Email: security@example.com
- GPG Key: [Link to public key]

## Resources

- [Go Security Guidelines](https://golang.org/doc/security)
- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Language_Security_Cheat_Sheet.html)
- [CIS Benchmarks](https://www.cisecurity.org/)
