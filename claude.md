# Project Specification Notes

Project is a ip / domain / whois / email reputation 



api routes : 

/ => 185.147.177.200
/185.147.177.200 -> 


https://echoip.ir/185.147.177.200 i want something like this 

## Frontend

### Stack

- - only Api
  
  

### Requirements

- All assets must be hosted on the server.
  
  this project just has api endpoints
  
  ratelimit should set as config 

## Backend

### Core Stack

- Golang 
- gin web framework
- Hexagonal Architecture
- Modular Structure
- redis (for ratelimit)

### Database

- ORM / PDO 
  Migrations
- `.env`
  
  ### Infrastructure

- Routes
- Cron Jobs
- Git
- Deployment

---

## AI / Agent Requirements

### Claude

- Agent Skills support
- Agent Command Support
  subagent for backend senior 
  subagent for forntend
  subagent for devops

### Rules

1. Responses must be in English.
2. Assets should remain inside the project.
3. Before starting any phase, ask questions if requirements are unclear.
4. Before each phase, write a specification (spec).
5. Maintain a skills folder/system.
6. Orders should be tracked in an `orders` file.
7. After each phase, provide deployment instructions.
8. UI should be implemented inside the UI module.
9. Use AJAX where appropriate.
10. Use Git semantic commit messages.
11. After each push, create a Git tag.
12. Before each phase, ask up to 5 clarification questions when necessary.

---

## Documentation

### Required Files

- `CHANGELOG.md`
- `ORDERS.md`

### Project Documentation

Include:

- Project description
- Architecture explanation
- Commands
- Server IP address
- Domain
- Deployment notes

---

## Recommendations

After each phase provide recommendations for:

- Database choice
- Deployment strategy
- Cron configuration
- Additional required skills
- Missing specifications

---

## Development Workflow

1. Define the phase.
2. Create/update specification.
3. Implement.
4. Verify and test.
5. Commit using semantic Git messages.
6. Push.
7. Tag release.
8. Update documentation.
9. Provide recommendations for the next phase.

---

## Preferred Frontend Stack

- Alpine.js
- Tailwind CSS
- DaisyUI

### Deployment

- use docker
