# Test chef proxy

# Create a chef
curl -H "Accept: application/jsonpb" \
-H "Content-Type: application/json" \
-X POST \
-d '{"chef": {"first_name": "Simon", "last_name": "Woldemichael", "password": "l337p@$$w0rd", "location": {"latitude": 33.5779, "longitude": 101.8552 }, "allergens": ["peas", "spinach", "peanuts"], "food_preferences": "NONE"}}' \
http://localhost:8081/v1/chef


# Get a user
curl localhost:8081/v1/chef/5c68f7c03ac52b1b9ce57315
