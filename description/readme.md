I`m splitting the online marketplace platform into four microservices. Every microservice concentrates on one business function.

Services:
User service - authorization users, handling account operations
Product service - managing product catalog and getting product list.
Order service - storing information about order
Delivery service - storing the history of delivery

Databases:
	SQL databases for user service, order service, and product service because these databases have a relational structure. Every table has defined schemas and every row has the same structure. This way delivers data consistency on account operations and orders.
NoSQL databases for product service and delivery service. First, most users spend time looking for items in the marketplace. NoSQL server is easy to horizontal scaling (add new server) and works with multiple instances, without replication. Delivery service tracks every detail about the package and the status of delivery.

Scaling:
Product service and delivery service can be scaled horizontally, the most of usage is searching products and changing delivery status.

CI/CD
CI - Linting and running tests before integration code with the main branch.
CD - Continuous Deployment - automatically building applications and deployment.

