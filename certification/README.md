# Certification C3 — Luis Alejandro Meneses Hernandez

## Cover

| Field | Value |
|---|---|
| Full name | Luis Alejandro Meneses Hernandez |
| GitHub user | MENESESLUIS-bot |
| Group | G2 |
| Team | LMS-Library |
| Team -docs repository | https://github.com/code-corhuila/library-docs |
| Course | HCCDP-Solution Architectures Course |
| Course link | https://e.huawei.com/en/talent/ict-academy/#/ictcourses-detail?courseId=6WtKiPd1S1YelBB0AUbE1vXlPgk |
| Start date | 2026-11-03 |
| Completion date | 2026-11-03 |

Certificate of Completion: `certificate.pdf` (code ICT20261103000003, issue date 2026-11-03).

## Study log

| Date | Hours | Chapters or labs | Evidence file |
|---|---|---|---|
| 2026-11-03 | 0.67 | Enrollment and Pre-class Guidance; Chapters 1 to 14; Lab 1 — Design Cases and Practices of an O&M Solution; Lab 2 — High-Performance Network Solution; Lab 3 — Highly Elastic Service System; Lab 4 — Highly Reliable Database Solution; Final exam | evidence/00-enrollment.png to evidence/14-chapter-14.png, evidence/15-lab-1.png, evidence/16-lab-2.png, evidence/17-lab-3.png, evidence/18-lab-4.png, evidence/19-final-exam.png |
| Total | 0.67 | | |

## Labs

| Evidence file | Lab name | What I built | What failed first and how I fixed it |
|---|---|---|---|
| evidence/15-lab-1.png | Design Cases and Practices of an O&M Solution | I configured Log Tank Service (LTS) to collect the logs of a cloud server and Cloud Eye alarm rules to monitor it. | <One sentence> |
| evidence/16-lab-2.png | High-Performance Network Solution | I built a cloud network with a VPC, an EIP, a NAT Gateway, VPC Peering and private DNS so that servers in different subnets could reach each other and the Internet. | <One sentence> |
| evidence/17-lab-3.png | Highly Elastic Service System | I deployed the Discuz! forum as a distributed system behind ELB with an Auto Scaling group that adds and removes servers based on load. | <One sentence> |
| evidence/18-lab-4.png | Highly Reliable Database Solution | I created a primary/standby RDS instance, connected the application to it and tested failover and backup restore. | <One sentence> |

## Final exam

- Score: 82 points (Passed, passing score 70/100)
- Date: 2026-11-03
- Evidence: evidence/19-final-exam.png

## Skills summary

I learned to evaluate a cloud design against the five principles of the course: reliability, security, performance efficiency, cost and operations.
I can spread stateless services over two availability zones behind a load balancer and scale them with auto scaling policies based on CPU or request rate.
I can choose a disaster recovery pattern from the RTO and RPO, and I now see why a daily backup cannot meet a 5-minute RPO while a cross-AZ replica set with log backups can.
I can protect a system in layers with a WAF at the edge, private subnets and security groups for the databases, and a secret manager instead of `.env` files.
I can design a CI/CD pipeline that builds an image once, promotes it from qa to production by digest and rolls back automatically when health checks fail.
I applied these techniques to redesign my team's LMS-Library system for Huawei Cloud in the design note of this folder.
