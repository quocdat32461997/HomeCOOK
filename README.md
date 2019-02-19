## Homecook @ HackDFW 2019

Project Structure
- data/
  - Datasets for training and testing price generation model (Grocery data, Cook Rate, Delivery Cost, Our fee)
- frontend/
  - Front-end design mockups and code
- go/
  - api/
    - Protocol Buffers + gRPC definitions
  - cmd/
    - Docker entrypoints/microservices
  - internal/
    - Service implementations
- kubernetes/
  - GKE configuration
- node/
  - Node.js API for pulling recipies
- scripts/
  - Some testing scripts

Ignored files
  1. .env - service environments
  2. homecook.json - Google Cloud Platform service account key
