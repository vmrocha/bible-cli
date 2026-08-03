# Security Policy

## Supported versions

Bible CLI is under active development. Security fixes are provided for
the latest published release and the current `main` branch. Older releases are
not supported; users should upgrade to the latest release before reporting a
problem that may already be fixed.

## Reporting a vulnerability

Please do not open a public issue, discussion, or pull request for a suspected
security vulnerability.

After private vulnerability reporting is enabled for this public repository,
use GitHub's **Report a vulnerability** form:

https://github.com/vmrocha/bible-cli/security/advisories/new

Until that form is available, or if it does not work, email
[vmrocha@gmail.com](mailto:vmrocha@gmail.com) with the subject
`[bible-cli security]`.

Include as much of the following as possible:

- the affected version or commit;
- the operating system and architecture;
- a description of the vulnerability and its potential impact;
- minimal reproduction steps or a proof of concept;
- any known mitigations or suggested fixes; and
- whether the issue is already public or has been disclosed elsewhere.

Do not include real credentials, personal data, or other sensitive third-party
information in a report. Use clearly marked test values and redact logs where
possible.

You should receive an acknowledgment within five business days. The maintainer
will aim to provide a status update at least every ten business days until the
report is resolved or closed. These are response targets, not guarantees.

Please allow time to investigate and release a fix before public disclosure.
The maintainer will coordinate disclosure timing and credit with the reporter.

## Scope

Security reports may concern the CLI, embedded data handling, translation
import pipeline, release artifacts, installation instructions, or build and
release automation. Bible translation wording issues and ordinary functional
bugs without a security impact should be reported through the public issue
tracker after the repository becomes public.
