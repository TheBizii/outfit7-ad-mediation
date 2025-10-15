# Outfit7 Ad Mediation

This is my solution to the **“Ad Mediation” expertise test** provided by Outfit7 as part of my recruitment process. The solution is written in **Go**, using the **Gin web framework**.

The goal of this document is to describe, in detail, my assumptions, design decisions, and other relevant information that may help you assess my solution.

---

# Setup instructions

## Prerequisites
- Basic knowledge of Unix terminals  
- Go compiler  
- golang-migrate (only if you intend to run database migrations; migrations can also be run manually in the Postgres CLI or via a web-based dashboard)  
- Docker Engine (optional, but highly recommended)  
- PostgreSQL (hosted either locally or remotely)

---

## Step 1: Clone the repository

You can clone this repository by running:

```bash
git clone git@github.com:TheBizii/outfit7-ad-mediation.git
```

After cloning, copy the .env.example file to .env and update it as needed.

---

## Step 2: Build the Docker container

To build the container from the provided Dockerfile, navigate to the project’s root directory and run:

```bash
docker build . -t outfit7-ad-mediation
```

---

## Step 3: Build the Postgres Docker container

Create a Docker network to connect the Go container with the Postgres container:
```bash
docker network create go-outfit7-network
```

Download and start the latest Postgres image:
```bash
docker run --name outfit7-postgres --network go-outfit7-network -e POSTGRES_USER=outfit7 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=outfit7 -p 5432:5432 -d postgres
```

Verify that the container is running:
```bash
docker ps
```

Run the database migrations. This step is especially important if you’ve just created a new Postgres container:
```bash
migrate -path=./migrations -database="postgres://outfit7:secret@localhost:5432/outfit7?sslmode=disable" up
```

Finally, update the project’s .env file by setting the PSQL_HOST value to outfit7-postgres, which matches the name of the Postgres container.

---

## Step 4: Start the Golang Docker container

To start the application, run:
```bash
docker run --rm --name go-outfit7-app --publish 8080:8080 --env-file .env --network go-outfit7-network -d outfit7-ad-mediation
```

That’s it, the app should now be running. You can proceed to test the solution or continue reading this document for more details.

---

# API Endpoints

This solution exposes three main API endpoints as requested in the expertise test instructions. The full documentation is available in the `docs/swagger.yaml` file. This section focuses on my design decisions.

---

## Retrieve List of Ad Networks (Mobile App)

**GET** `/api/v1/ad_networks/{countryCode}/{adType}`

This endpoint is intended to be called from mobile apps to retrieve an ordered list of ad networks for the specified ad type. I decided to expose a **GET** endpoint at `/api/v1/ad_networks/{countryCode}/{adType}`.

It accepts two **required** parameters:

- `countryCode`: the ISO country code  
- `adType`: the ad type (e.g., banner, interstitial, rewarded_video)

I thought for quite some time whether `adType` was really necessary and came to the conclusion that mobile apps almost always know in advance what type of ad they'd like to display, depending on the placement or opportunity in the app. Therefore, it made sense to me to include it as a required parameter.

A mobile app can also optionally send additional query parameters that provide contextual information useful for filtering suitable ad networks:

- `platform`: which operating system the app is running on  
- `osVersion`: the OS version (ignored if `platform` is not provided)  
- `appName`: name of the mobile app  
- `appVersion`: version of the mobile app

---

## Dashboard Summary of All Priority Lists

**GET** `/api/v1/ad_networks/dashboard`

This endpoint is designed for internal dashboard use to visualize key information across countries and ad types. For example, it could be used to calculate the most performant ad networks across regions, or, if the backend tracked update history, to analyze correlations between ad network performance and factors like user age group or holidays.

My assumption is that the dashboard will typically request **all** priority lists at once, so this endpoint currently does not take any input parameters.
If the number of countries or ad types grows significantly, we could later add filtering (for example, by country or ad type) or pagination support.

---

## Update Ad Networks

**PUT** `/api/v1/ad_networks/{countryCode}/{adType}`

This endpoint updates or creates the ad network priority list for a given country and ad type. This design decision took the most consideration. I wasn’t sure whether Outfit7’s internal batch processes would prefer updating one country/ad type at a time, or performing a bulk update for multiple pairs at once (to reduce the number of network calls).

Both approaches have tradeoffs. I ultimately chose to update one list per request because it more closely follows REST principles and keeps the logic clean and predictable. In practice, a performance test would determine whether a batch update endpoint would be more efficient at this scale. Since I wasn’t provided with information about the expected number of countries or ad types, this single-update approach felt most appropriate.

**Possible improvement:**  
This endpoint currently lacks built-in authentication, meaning anyone could theoretically update the priority lists if the API were publicly exposed. In production, this should be protected using an authentication layer, for example, by requiring a Bearer token in the request headers and validating it through middleware before allowing access.

---

# Possible Improvements

## Support for More Advanced Filters

When displaying ads, especially those visible to children, we must take precautions to ensure safety. One improvement would be to extend the database schema to store recommended **age groups** or **content categories** for each ad network. This would not only improve safety but also allow for deeper insights into which networks perform best for specific demographics.

---

## Add an Authenitcation Layer

As mentioned above, my top priority for improvement would be implementing authentication and authorization. This is typically done via a **middleware** that validates tokens (for example, JWT or Bearer tokens).

---

## Deployment to GCP Using Cloud Run

Although this solution does not currently deploy to Google Cloud Platform, I kept GCP deployment in mind during development. Rather than purchasing a GCP plan for this test, I chose to **Dockerize** the application. This makes deployment to **Cloud Run**, **Kubernetes**, or any other container orchestration service much easier.

Additionally, I designed the Go service to be containerized separately from the database. In production, there would likely be multiple instances of the Go backend running behind a load balancer, all connected to the same database (distributed or replicated as needed).

---

