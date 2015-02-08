# TimeOff

TimeOff is a web application for tracking a worker's annual leave (vacation and public holidays),
sickness and time in lieu.

## HTTP/2

Timeoff supports the HTTP/2 protocol over TLS for browsers that understand it.
Other browsers will continue to use HTTP/1.1.

## Project status

Currently in alpha. Many features are not yet implemented.

### Features TODO

- Annual leave
- Public holidays/bank holidays
- Sickness
- Time off in lieu (TOIL)
- Managerial approvals
- Annotations (e.g. "holiday to Spain")
- Custom leave day types (e.g. privilege days)
- Custom leave years (per-employee or global for whole company)
- Reporting
- Notifications to employee and manager to notify of time not taken, nearing end of year, etc.
- Currently assumes full time, need to allow for part-time

## Supported user types TODO

- Organisation admin (super user)
- Line manager
- Read-only user (auditor)
- Employee
