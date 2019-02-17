# Test user proxy

# Create a user
curl -H "Accept: application/jsonpb" \
-H "Content-Type: application/json" \
-X POST \
-d '{"user": {"first_name": "Simon", "last_name": "Woldemichael", "password": "l337p@$$w0rd", "location": {"latitude": 33.5779, "longitude": 101.8552 }, "allergens": ["peas", "spinach", "peanuts"], "food_preferences": "NONE"}}' \
http://localhost:8080/v1/user


# Get a user
# http://localhost:8080/v1/user/
curl localhost:8080/v1/user/5c68cf8d680afd7cda4146b1
