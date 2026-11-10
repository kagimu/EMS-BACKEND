# EMS Backend Structure

EMS uses one package per business feature. Each feature follows the same shape:

```
domain/ems/<feature>/
├── api/
│   ├── requests/     request DTOs and validation
│   └── responses/    response DTOs
├── controllers/      Gin handlers
├── models/           GORM models
├── repositories/     database queries (when needed)
├── services/         business rules and transactions
└── middleware/       feature-specific guards (when needed)
```

Routes are registered centrally in `api/routes/Routes.go`; the controller imports
the feature service rather than putting business logic in the handler.

## EMS feature map

| Domain | Responsibility |
|---|---|
| `organisation` | Organisation profile, units, custom fields |
| `stakeholders` | Owners, people, employee intake |
| `assignments` | User assignments |
| `authority` | Decision authority |
| `reporting` | Reporting lines |
| `workflow` | Approval workflow engine |
| `licensing` | Licences, scopes, entitlements, units, activation, enforcement |
| `addons` | Addon/module catalogue and enablement state |
| `audit` | Organisation and registry audit |
| `notifications` | In-app notifications |
| `support` | Tickets, messages, attachments, support staff and operations |
| `commercial` | Customer profiles, commercial products, commercial audit |
| `owner` | Provider/owner console read models and aggregates |

## Conventions

- Use UUID `Number` fields and `gorm.Model` timestamps/soft delete unless a
  model has a specific reason not to.
- Keep request/response DTOs separate from GORM models.
- Put transactional or multi-record logic in `services/`, not controllers.
- Register migrations additively in `migrations/Migrations.go`.
- Record immutable history in the relevant event/audit table; do not overwrite
  historical facts.
- Owner-only/provider APIs use the explicit software-owner guard. Do not grant
  them through ordinary client Admin permissions.
- Integration tests live in `backend/integration` and use the disposable test
  database (`db_wholesome_e2_test`), never production.

## Related modules outside `domain/ems`

- `domain/commercial` — onboarding/setup/subscription order and package records.
- `domain/auth` — users, JWT, passwords.
- `domain/permissions` — roles and permissions.
- `domain/eos` — entities, memberships and active organisation.
- `domain/ems/commercial` — commercial customer/product catalogue and audit.

The Education addon uses the same backend `domain/<feature>` pattern and the same
central route registration, so a developer familiar with Education can navigate
EMS by finding the matching feature package.
