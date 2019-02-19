# Test user proxy

# Create a user
curl -H "Accept: application/jsonpb" \
-H "Content-Type: application/json" \
-X POST \
-d '{"user": {"first_name": "Simon", "last_name": "Woldemichael", "password": "l337p@$$w0rd", "location": {"latitude": 33.5779, "longitude": 101.8552 }, "allergens": ["peas", "spinach", "peanuts"], "food_preferences": "NONE"}}' \
http://35.193.17.77:8080/v1/user


# Get a user
curl 35.193.17.77:8080/v1/user/5c68f7c13ac52b1b9ce57316
