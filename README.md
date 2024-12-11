# PointaFam

## Introduction

**PointaFam** is an innovative e-commerce platform bridging the gap between local farmers and urban retailers. By connecting farmers directly to retailers, PointaFarm ensures the delivery of fresh, high-quality, and fairly priced produce. The platform focuses on reducing supply chain inefficiencies, empowering farmers, and making sourcing easier for retailers. Built with cutting-edge technologies, it emphasizes simplicity, reliability, and seamless user experiences.

- **Backend**: Golang, Gin-Gonic, SQLite  
- **Frontend**: HTMX, TailwindCSS  

---

### Author's LinkedIn  
Connect with me on LinkedIn: [Ferdynand Odhiambo](https://www.linkedin.com/in/ferdynand-odhiambo)

---

## Installation

To set up and run PointaFam locally, follow these steps:

### Prerequisites  

Ensure the following are installed on your machine:
- [Go](https://golang.org/doc/install) (version 1.19+ recommended)  
- [SQLite](https://www.sqlite.org/download.html)  
- Git  

---

### Clone the Repository  

```bash
git clone https://github.com/MeFerdi/pointa-fam.git

cd pointa-fam
```
## Backend Setup

- Install Go dependencies:
```bash
go mod download
```

- Run database migrations:
```bash

go run migrations/migrate.go
```
- Start the Backend Server:
```bash

go run main.go
```
- The backend server will be accessible at http://localhost:8080.

### Frontend Setup
PointaFarm's frontend is built with HTMX and TailwindCSS. Serve the frontend files directly from the static folder using the backend, or integrate them into your workflow.

### Usage
- For Retailers:

**Browse Farms**: Discover farms and available produce based on location and preferences.

**Filter Products**: Use filters for organic produce, price range, or availability.

**Add to Cart**: Add selected products to your cart and prepare to check out.

- For Farmers:

**Register Farms**: Farmers can create accounts, register farms, and list available produce.

**Update Products**: Maintain up-to-date product lists, pricing, and seasonal offerings.

**Manage Orders**: View and manage incoming orders, ensuring timely delivery to retailers.

### Developer Usage
- Testing Locally:
Ensure both the backend server and frontend are properly set up and running.
Access the backend at http://localhost:8080.

### Contributing
Contributions are welcome! Follow the steps below to contribute:

- Fork the Repository:
Click the "Fork" button at the top of this repository to create a copy under your GitHub account.
Clone Your Fork:
```bash

git clone https://github.com/your-username/pointa-fam.git
cd pointa-fam
```
- Create a New Branch:
```bash

git checkout -b feature/your-feature-name
```
- Make Changes and Commit:
- Make code changes.
- Add tests where applicable.
- Run tests to ensure everything works.
- Commit your changes:
```bash

git add .
git commit -m "Add feature: your feature description"
```
- Push Your Changes:
```bash

git push origin feature/your-feature-name
```
- Open a Pull Request:
- Go to the original repository on GitHub.
- Click "New Pull Request."
- Select your branch and submit the pull request.