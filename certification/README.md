# Certification C3 — Luis Alejandro Meneses

## Cover

| Field | Value |
|---|---|
| Full name | Luis Alejandro Meneses |
| GitHub user | MENESESLUIS-bot |
| Group | G2 |
| Team | LMS-Library |
| Team -docs repository | https://github.com/code-corhuila/library-docs |
| Course | HCCDP-Solution Architectures Course |
| Course link | https://e.huawei.com/en/talent/ict-academy/#/ictcourses-detail?courseId=6WtKiPd1S1YelBB0AUbE1vXlPgk |
| Start date | <YYYY-MM-DD> |
| Completion date | <YYYY-MM-DD> |

## Study log

| Date | Hours | Chapters or labs | Evidence file |
|---|---|---|---|
| <YYYY-MM-DD> | <h> | Pre-class Guidance, Chapter 1 Architecture Design Methodology on Huawei Cloud | evidence/00-enrollment.png, evidence/01-chapter-01.png |
| <YYYY-MM-DD> | <h> | Chapter 2 Typical Service Implementation Solutions on the Cloud | evidence/02-chapter-02.png |
| <YYYY-MM-DD> | <h> | Chapter 3 Cloud Compute Solution Design | evidence/03-chapter-03.png |
| <YYYY-MM-DD> | <h> | Chapter 4 Cloud Storage Solution Design | evidence/04-chapter-04.png |
| <YYYY-MM-DD> | <h> | Chapter 5 Cloud Network Solution Design, Chapter 6 Cloud Database Solution Design | evidence/05-chapter-05.png, evidence/06-chapter-06.png |
| <YYYY-MM-DD> | <h> | Chapter 7 Cloud Security Solution Design | evidence/07-chapter-07.png |
| <YYYY-MM-DD> | <h> | Chapter 8 Cloud O&M Solution Design | evidence/08-chapter-08.png |
| <YYYY-MM-DD> | <h> | Chapter 9 Distributed Architecture Design | evidence/09-chapter-09.png |
| <YYYY-MM-DD> | <h> | Chapter 10 Highly Scalable System Design, Chapter 11 Disaster Recovery System Design | evidence/10-chapter-10.png, evidence/11-chapter-11.png |
| <YYYY-MM-DD> | <h> | Chapter 12 Automated Deployment, Chapter 13 Containers and Agility | evidence/12-chapter-12.png, evidence/13-chapter-13.png |
| <YYYY-MM-DD> | <h> | Chapter 14 Huawei Cloud Advanced Services, Lab 1, Lab 2 | evidence/14-chapter-14.png, evidence/15-lab-1.png, evidence/16-lab-2.png |
| <YYYY-MM-DD> | <h> | Lab 3, Lab 4, Final Exam | evidence/17-lab-3.png, evidence/18-lab-4.png, evidence/19-final-exam.png |
| Total | <sum> | | |

## Labs

| Evidence file | Lab name | What I built | What failed first and how I fixed it |
|---|---|---|---|
| evidence/15-lab-1.png | Design Cases and Practices of an O&M Solution | I configured Log Tank Service (LTS) to collect the logs of a cloud server and Cloud Eye alarm rules to monitor it. | <One sentence> |
| evidence/16-lab-2.png | High-Performance Network Solution | I built a cloud network with a VPC, an EIP, a NAT Gateway, VPC Peering and private DNS so that servers in different subnets could reach each other and the Internet. | <One sentence> |
| evidence/17-lab-3.png | Highly Elastic Service System | I deployed the Discuz! forum as a distributed system behind ELB with an Auto Scaling group that adds and removes servers based on load. | <One sentence> |
| evidence/18-lab-4.png | Highly Reliable Database Solution | I created a primary/standby RDS instance, connected the application to it and tested failover and backup restore. | <One sentence> |

## Final exam

- Score: 82/100 (passing score 70/100)
- Date: <YYYY-MM-DD>
- Evidence: evidence/19-final-exam.png

## Skills summary

I learned to evaluate a cloud design against the five principles of the course: reliability, security, performance efficiency, cost and operations.
I can spread stateless services over two availability zones behind a load balancer and scale them with auto scaling policies based on CPU or request rate.
I can choose a disaster recovery pattern from the RTO and RPO, and I now see why a daily backup cannot meet a 5-minute RPO while a cross-AZ replica set with log backups can.
I can protect a system in layers with a WAF at the edge, private subnets and security groups for the databases, and a secret manager instead of `.env` files.
I can design a CI/CD pipeline that builds an image once, promotes it from qa to production by digest and rolls back automatically when health checks fail.
I applied these techniques to redesign my team's LMS-Library system for Huawei Cloud in the design note of this folder.
